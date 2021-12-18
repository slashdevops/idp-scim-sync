package core

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	log "github.com/sirupsen/logrus"
	"github.com/slashdevops/idp-scim-sync/internal/model"
	"github.com/slashdevops/idp-scim-sync/internal/utils"
	"github.com/slashdevops/idp-scim-sync/internal/version"
)

var (
	// ErrIdentiyProviderServiceNil is returned when the Identity Provider Service is nil
	ErrIdentiyProviderServiceNil = errors.New("identity provider service cannot be nil")

	// ErrSCIMServiceNil is returned when the SCIM Service is nil
	ErrSCIMServiceNil = errors.New("SCIM service cannot be nil")

	// ErrStateRepositoryNil is returned when the State Repository is nil
	ErrStateRepositoryNil = errors.New("state repository cannot be nil")
)

// SyncService represent the sync service and the core of the sync process
type SyncService struct {
	provGroupsFilter []string
	provUsersFilter  []string
	prov             IdentityProviderService
	scim             SCIMService
	repo             StateRepository
}

// NewSyncService creates a new sync service.
func NewSyncService(prov IdentityProviderService, scim SCIMService, repo StateRepository, opts ...SyncServiceOption) (*SyncService, error) {
	if prov == nil {
		return nil, ErrIdentiyProviderServiceNil
	}
	if scim == nil {
		return nil, ErrSCIMServiceNil
	}
	if repo == nil {
		return nil, ErrStateRepositoryNil
	}

	ss := &SyncService{
		prov:             prov,
		provGroupsFilter: []string{}, // fill in with the opts
		provUsersFilter:  []string{}, // fill in with the opts
		scim:             scim,
		repo:             repo,
	}

	for _, opt := range opts {
		opt(ss)
	}

	return ss, nil
}

// SyncGroupsAndTheirMembers the default sync method tha syncs groups and their members
func (ss *SyncService) SyncGroupsAndTheirMembers(ctx context.Context) error {
	log.WithFields(log.Fields{
		"group_filter": ss.provGroupsFilter,
	}).Info("getting Identity Provider data")

	idpGroupsResult, err := ss.prov.GetGroups(ctx, ss.provGroupsFilter)
	if err != nil {
		return fmt.Errorf("error getting groups from the identity provider: %w", err)
	}

	idpUsersResult, err := ss.prov.GetUsers(ctx, []string{""})
	if err != nil {
		return fmt.Errorf("error getting users from the identity provider: %w", err)
	}

	idpGroupsMembersResult, err := ss.prov.GetGroupsMembers(ctx, idpGroupsResult)
	if err != nil {
		return fmt.Errorf("error getting groups members: %w", err)
	}

	if idpUsersResult.Items == 0 {
		log.Warn("there are no users in the identity provider")
	}

	if idpGroupsResult.Items == 0 {
		log.Warnf("there are no groups in the identity provider that match with this filter: %s", ss.provGroupsFilter)
	}

	if idpGroupsMembersResult.Items == 0 {
		log.Warn("there are no groups with members in the identity provider")
	}

	log.Info("getting state data")
	state, err := ss.repo.GetState(ctx)
	if err != nil {
		var nsk *types.NoSuchKey
		if errors.As(err, &nsk) {
			log.Warn("no state file found in the state repository, creating this")
			state = &model.State{}
		} else {
			return fmt.Errorf("error getting state data from the repository: %w", err)
		}
	}

	// these variables are used to store the data that will be used to create, delete and are equal for users and groups in SCIM
	// the differents between the data in the identity provider and these is that these have already have the SCIMID
	// after the creation of the element in SCIM
	var (
		totalGroupsResult        model.GroupsResult
		totalUsersResult         model.UsersResult
		totalGroupsMembersResult model.GroupsMembersResult
	)

	// first time syncing
	if state.LastSync == "" {
		// check SCIM side to see if there are elements to be
		// reconciled. Basically, checks if SCIM is not clean before the first sync
		// and we need to reconcile the SCIM side with the identity provider side.
		// In case of migration from a different tool and we want to keep the state
		// of the users and groups in the SCIM side, just no recreation, keep the existing ones when the:n
		// - Groups names are equals on both sides, update only the external id (coming from the identity provider)
		// - Users emails are equals on both sides, update only the external id (coming from the identity provider)

		log.Warn("syncing from scim service, first time syncing")
		log.Warn("reconciling the SCIM data with the Identity Provider data")

		log.Info("getting SCIM Groups")
		scimGroupsResult, err := ss.scim.GetGroups(ctx)
		if err != nil {
			return fmt.Errorf("error getting groups from the SCIM service: %w", err)
		}

		log.WithFields(log.Fields{
			"idp":  idpGroupsResult.Items,
			"scim": scimGroupsResult.Items,
		}).Info("reconciling groups")
		groupsCreate, groupsUpdate, groupsEqual, groupsDelete, err := model.GroupsOperations(idpGroupsResult, scimGroupsResult)
		if err != nil {
			return fmt.Errorf("error reconciling groups: %w", err)
		}

		groupsCreated, groupsUpdated, err := reconcilingGroups(ctx, ss.scim, groupsCreate, groupsUpdate, groupsDelete)
		if err != nil {
			return fmt.Errorf("error reconciling groups: %w", err)
		}

		// groupsCreated + groupsUpdated + groupsEqual = groups total
		totalGroupsResult = model.MergeGroupsResult(groupsCreated, groupsUpdated, groupsEqual)

		log.Info("getting SCIM Users")
		scimUsersResult, err := ss.scim.GetUsers(ctx)
		if err != nil {
			return fmt.Errorf("error getting users from the SCIM service: %w", err)
		}

		log.WithFields(log.Fields{
			"idp":  idpUsersResult.Items,
			"scim": scimUsersResult.Items,
		}).Info("reconciling users")
		usersCreate, usersUpdate, usersEqual, usersDelete, err := model.UsersOperations(idpUsersResult, scimUsersResult)
		if err != nil {
			return fmt.Errorf("error operating with users: %w", err)
		}

		usersCreated, usersUpdated, err := reconcilingUsers(ctx, ss.scim, usersCreate, usersUpdate, usersDelete)
		if err != nil {
			return fmt.Errorf("error reconciling users: %w", err)
		}

		// usersCreated + usersUpdated + usersEqual = users total
		totalUsersResult = model.MergeUsersResult(usersCreated, usersUpdated, usersEqual)

		log.Info("getting SCIM Groups Members")
		// unfortunately, the SCIM service does not support the getGroupsMembers method in and efficient way
		// see: "Nor Supported" section in: https://docs.aws.amazon.com/singlesignon/latest/developerguide/listgroups.html
		// scimGroupsMembersResult, err := ss.scim.GetGroupsMembers(ctx, &totalGroupsResult) // not supported yet
		scimGroupsMembersResult, err := ss.scim.GetGroupsMembersBruteForce(ctx, &totalGroupsResult, &totalUsersResult)
		if err != nil {
			return fmt.Errorf("error getting groups members from the SCIM service: %w", err)
		}

		log.Tracef("idpGroupsMembersResult: %s\n, scimGroupsMembersResult: %s\n", utils.ToJSON(idpGroupsMembersResult), utils.ToJSON(scimGroupsMembersResult))

		log.WithFields(log.Fields{
			"idp":  idpGroupsMembersResult.Items,
			"scim": scimGroupsMembersResult.Items,
		}).Info("reconciling groups members")
		membersCreate, membersEqual, membersDelete, err := model.MembersOperations(idpGroupsMembersResult, scimGroupsMembersResult)
		if err != nil {
			return fmt.Errorf("error reconciling groups members: %w", err)
		}

		log.Tracef("membersCreate: %s\n, membersEqual: %s\n, membersDelete: %s\n", utils.ToJSON(membersCreate), utils.ToJSON(membersEqual), utils.ToJSON(membersDelete))

		membersCreated, err := reconcilingGroupsMembers(ctx, ss.scim, membersCreate, membersDelete)
		if err != nil {
			return fmt.Errorf("error reconciling groups members: %w", err)
		}

		// membersCreate + membersEqual = members total
		totalGroupsMembersResult = model.MergeGroupsMembersResult(membersCreated, membersEqual)
	} else { // This is not the first time syncing
		lastSyncTime, err := time.Parse(time.RFC3339, state.LastSync)
		if err != nil {
			return fmt.Errorf("error parsing last sync time: %w", err)
		}

		log.WithFields(log.Fields{
			"lastsync": state.LastSync,
			"since":    time.Since(lastSyncTime).String(),
		}).Info("syncing from state")

		if idpGroupsResult.HashCode == state.Resources.Groups.HashCode {
			log.Info("provider groups and state groups are the same, nothing to do with groups")

			totalGroupsResult = state.Resources.Groups
		} else {
			log.Info("provider groups and state groups are different")
			// now here we have the google fresh data and the last sync data state
			// we need to compare the data and decide what to do
			// see differences between the two datasets

			log.WithFields(log.Fields{
				"idp":   idpGroupsResult.Items,
				"state": state.Resources.Groups.Items,
			}).Info("reconciling groups")
			groupsCreate, groupsUpdate, groupsEqual, groupsDelete, err := model.GroupsOperations(idpGroupsResult, &state.Resources.Groups)
			if err != nil {
				return fmt.Errorf("error reconciling groups: %w", err)
			}

			groupsCreated, groupsUpdated, err := reconcilingGroups(ctx, ss.scim, groupsCreate, groupsUpdate, groupsDelete)
			if err != nil {
				return fmt.Errorf("error reconciling groups: %w", err)
			}

			// merge in only one data structure the groups created, updated amd equals who has the SCIMID
			totalGroupsResult = model.MergeGroupsResult(groupsCreated, groupsUpdated, groupsEqual)
		}

		if idpUsersResult.HashCode == state.Resources.Users.HashCode {
			log.Info("provider users and state users are the same, nothing to do with users")

			totalUsersResult = state.Resources.Users
		} else {
			log.Info("provider users and state users are different")

			log.WithFields(log.Fields{
				"idp":   idpUsersResult.Items,
				"state": state.Resources.Users.Items,
			}).Info("reconciling users")
			usersCreate, usersUpdate, usersEqual, usersDelete, err := model.UsersOperations(idpUsersResult, &state.Resources.Users)
			if err != nil {
				return fmt.Errorf("error operating with users: %w", err)
			}

			usersCreated, usersUpdated, err := reconcilingUsers(ctx, ss.scim, usersCreate, usersUpdate, usersDelete)
			if err != nil {
				return fmt.Errorf("error reconciling users: %w", err)
			}

			// usersCreated + usersUpdated + usersEqual = users total
			totalUsersResult = model.MergeUsersResult(usersCreated, usersUpdated, usersEqual)
		}

		if idpGroupsMembersResult.HashCode == state.Resources.GroupsMembers.HashCode {
			log.Info("provider groups-members and state groups-members are the same, nothing to do with groups-members")

			totalGroupsMembersResult = state.Resources.GroupsMembers
		} else {
			log.Info("provider groups-members and state groups-members are different")

			// if we create a group or user during the sync, we need the scimid of these new groups/users
			// because to add members to a group the scim api needs that.
			// so this function will fill the scimid of the new groups/users
			groupsMembers := model.UpdateGroupsMembersSCIMID(idpGroupsMembersResult, &totalGroupsResult, &totalUsersResult)

			log.Tracef("groupsMembers: %s\n", utils.ToJSON(groupsMembers))

			log.WithFields(log.Fields{
				"idp":   idpGroupsMembersResult.Items,
				"state": state.Resources.GroupsMembers.Items,
			}).Info("reconciling groups members")

			membersCreate, membersEqual, membersDelete, err := model.MembersOperations(groupsMembers, &state.Resources.GroupsMembers)
			if err != nil {
				return fmt.Errorf("error reconciling groups members: %w", err)
			}

			log.Tracef("membersCreate: %s\n, membersEqual: %s\n, membersDelete: %s\n", utils.ToJSON(membersCreate), utils.ToJSON(membersEqual), utils.ToJSON(membersDelete))

			_, err = reconcilingGroupsMembers(ctx, ss.scim, membersCreate, membersDelete)
			if err != nil {
				return fmt.Errorf("error reconciling groups members: %w", err)
			}

			totalGroupsMembersResult = model.MergeGroupsMembersResult(groupsMembers)
		}
	}

	// after be sure all the SCIM side is aligned with the Identity Provider side
	// we can update the state with the last data coming from the reconciliation
	newState := &model.State{
		Resources: model.StateResources{
			Groups:        totalGroupsResult,
			Users:         totalUsersResult,
			GroupsMembers: totalGroupsMembersResult,
		},
	}
	// calculate the hash with the data payload
	newState.SetHashCode()
	newState.SchemaVersion = model.StateSchemaVersion
	newState.CodeVersion = version.Version
	newState.LastSync = time.Now().Format(time.RFC3339)

	log.WithFields(log.Fields{
		"lastSycn": newState.LastSync,
		"groups":   totalGroupsResult.Items,
		"users":    totalUsersResult.Items,
	}).Info("storing the new state")

	// TODO: avoid this step using a cmd flag, could be a nice feature
	if err := ss.repo.SetState(ctx, newState); err != nil {
		return fmt.Errorf("error storing the state: %w", err)
	}

	log.Tracef("state data: %s", utils.ToJSON(newState))
	log.Info("sync completed")
	return nil
}

// SyncGroupsAndUsers this method is used to sync the usersm groups and their members from the identity provider to the SCIM
func (ss *SyncService) SyncGroupsAndUsers() error {
	return errors.New("not implemented yet")
}
