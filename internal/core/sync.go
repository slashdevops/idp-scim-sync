package core

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sync"
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
	ctx              context.Context
	mu               *sync.RWMutex
	provGroupsFilter []string
	provUsersFilter  []string
	prov             IdentityProviderService
	scim             SCIMService
	repo             StateRepository
}

// NewSyncService creates a new sync service.
func NewSyncService(ctx context.Context, prov IdentityProviderService, scim SCIMService, repo StateRepository, opts ...SyncServiceOption) (*SyncService, error) {
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
		ctx:              ctx,
		mu:               &sync.RWMutex{},
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
func (ss *SyncService) SyncGroupsAndTheirMembers() error {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	log.WithFields(log.Fields{
		"group_filter": ss.provGroupsFilter,
	}).Info("getting Identity Provider data")

	idpGroupsResult, err := ss.prov.GetGroups(ss.ctx, ss.provGroupsFilter)
	if err != nil {
		return fmt.Errorf("error getting groups from the identity provider: %w", err)
	}

	log.Tracef("idpGroupsResult: %s\n", utils.ToJSON(idpGroupsResult))

	idpUsersResult, err := ss.prov.GetUsers(ss.ctx, []string{""})
	if err != nil {
		return fmt.Errorf("error getting users from the identity provider: %w", err)
	}

	log.Tracef("idpUsersResult: %s\n", utils.ToJSON(idpUsersResult))

	idpGroupsMembersResult, err := ss.prov.GetGroupsMembers(ss.ctx, idpGroupsResult)
	if err != nil {
		return fmt.Errorf("error getting groups members: %w", err)
	}

	log.Tracef("idpGroupsMembersResult: %s\n", utils.ToJSON(idpGroupsMembersResult))

	if idpUsersResult.Items == 0 {
		log.Warn("there are no users in the identity provider")
	}

	if idpGroupsResult.Items == 0 {
		log.Warnf("there are no groups in the identity provider that match with this filter: %s", ss.provGroupsFilter)
	}

	if idpGroupsMembersResult.Items == 0 {
		log.Warn("there are no group with members in the identity provider")
	}

	log.Info("getting state data")
	state, err := ss.repo.GetState(ss.ctx)
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
	var totalGroupsResult model.GroupsResult
	var totalUsersResult model.UsersResult
	var totalGroupsMembersResult model.GroupsMembersResult

	// first time syncing
	if state.LastSync == "" {
		// Check SCIM side to see if there are elelemnts to be
		// reconciled. Bassically check if SCIM is not clean before the first sync
		// and we need to reconcile the SCIM side with the identity provider side.
		// In case of migration from a different tool and we want to keep the state
		// of the users and groups in the SCIM side, just no recreate, keep the existing ones when the:
		// - Groups names are equals in both sides
		// - Users emails are equals in both sides

		log.Warn("syncing from scim service, first time syncing")
		log.Warn("reconciling the SCIM data with the Identity Provider data")

		log.Info("getting SCIM Groups")
		scimGroupsResult, err := ss.scim.GetGroups(ss.ctx)
		if err != nil {
			return fmt.Errorf("error getting groups from the SCIM service: %w", err)
		}

		log.WithFields(log.Fields{
			"idp":  idpGroupsResult.Items,
			"scim": scimGroupsResult.Items,
		}).Info("reconciling groups")
		groupsCreate, groupsUpdate, groupsEqual, groupsDelete := groupsOperations(idpGroupsResult, scimGroupsResult)

		groupsCreated, groupsUpdated, err := reconcilingGroups(ss.ctx, ss.scim, groupsCreate, groupsUpdate, groupsDelete)
		if err != nil {
			return fmt.Errorf("error reconciling groups: %w", err)
		}

		// groupsCreated + groupsUpdated + groupsEqual = groups total
		totalGroupsResult = mergeGroupsResult(groupsCreated, groupsUpdated, groupsEqual)

		log.Info("getting SCIM Users")
		scimUsersResult, err := ss.scim.GetUsers(ss.ctx)
		if err != nil {
			return fmt.Errorf("error getting users from the SCIM service: %w", err)
		}

		log.WithFields(log.Fields{
			"idp":  idpUsersResult.Items,
			"scim": scimUsersResult.Items,
		}).Info("reconciling users")
		usersCreate, usersUpdate, usersEqual, usersDelete := usersOperations(idpUsersResult, scimUsersResult)

		usersCreated, usersUpdated, err := reconcilingUsers(ss.ctx, ss.scim, usersCreate, usersUpdate, usersDelete)
		if err != nil {
			return fmt.Errorf("error reconciling users: %w", err)
		}

		// usersCreated + usersUpdated + usersEqual = users total
		totalUsersResult = mergeUsersResult(usersCreated, usersUpdated, usersEqual)

		// log.Tracef("totalGroupsResult: %s", utils.ToJSON(totalGroupsResult))

		log.Info("getting SCIM Groups Members")
		// scimGroupsMembersResult, err := ss.scim.GetGroupsMembers(ss.ctx, &totalGroupsResult) // not supported yet
		scimGroupsMembersResult, err := ss.scim.GetGroupsMembersBruteForce(ss.ctx, &totalGroupsResult, &totalUsersResult)
		if err != nil {
			return fmt.Errorf("error getting groups members from the SCIM service: %w", err)
		}

		// log.Tracef("scimGroupsMembersResult: %s", utils.ToJSON(scimGroupsMembersResult))

		log.WithFields(log.Fields{
			"idp":  idpGroupsMembersResult.Items,
			"scim": scimGroupsMembersResult.Items,
		}).Info("reconciling groups members")
		membersCreate, membersEqual, membersDelete := membersOperations(idpGroupsMembersResult, scimGroupsMembersResult)

		// log.Tracef("membersCreate: %s, membersEqual: %s, membersDelete: %s", utils.ToJSON(membersCreate), utils.ToJSON(membersEqual), utils.ToJSON(membersDelete))

		membersCreated, err := reconcilingGroupsMembers(ss.ctx, ss.scim, membersCreate, membersDelete)
		if err != nil {
			return fmt.Errorf("error reconciling groups members: %w", err)
		}

		// log.Tracef("membersCreated: %s\n, membersEqual: %s\n", utils.ToJSON(membersCreated), utils.ToJSON(membersEqual))
		// membersCreate + membersEqual = members total
		totalGroupsMembersResult = mergeGroupsMembersResult(membersCreated, membersEqual)

		// log.Tracef("totalGroupsMembersResult: %s", utils.ToJSON(totalGroupsMembersResult))

	} else { // This is not the first time syncing

		lastSyncTime, err := time.Parse(time.RFC3339, state.LastSync)
		if err != nil {
			return fmt.Errorf("error parsing last sync time: %w", err)
		}
		deltaTime := time.Since(lastSyncTime)

		deltaHours := fmt.Sprintf("%.0f", math.Floor(deltaTime.Hours()))
		deltaMinutes := fmt.Sprintf("%.0f", math.Floor(deltaTime.Minutes()))
		deltaSeconds := fmt.Sprintf("%.0f", math.Floor(deltaTime.Seconds()))

		log.WithFields(log.Fields{
			"lastsync": state.LastSync,
			"since":    deltaHours + "h, " + deltaMinutes + "m, " + deltaSeconds + "s",
		}).Info("syncing from state")

		log.Tracef("state.Resources.Groups: %s", utils.ToJSON(state.Resources.Groups))

		if idpGroupsResult.HashCode == state.Resources.Groups.HashCode {
			log.Info("provider groups and state groups are the same, nothing to do with groups")

			totalGroupsResult = state.Resources.Groups
		} else {
			log.Info("provider groups and state groups are diferent")
			// now here we have the google fresh data and the last sync data state
			// we need to compare the data and decide what to do
			// see differences between the two data sets

			log.WithFields(log.Fields{
				"idp":   idpGroupsResult.Items,
				"state": state.Resources.Groups.Items,
			}).Info("reconciling groups")
			groupsCreate, groupsUpdate, groupsEqual, groupsDelete := groupsOperations(idpGroupsResult, &state.Resources.Groups)

			groupsCreated, groupsUpdated, err := reconcilingGroups(ss.ctx, ss.scim, groupsCreate, groupsUpdate, groupsDelete)
			if err != nil {
				return fmt.Errorf("error reconciling groups: %w", err)
			}

			// merge in only one data structure the groups created and updated who has the SCIMID
			totalGroupsResult = mergeGroupsResult(groupsCreated, groupsUpdated, groupsEqual)
		}

		log.Tracef("state.Resources.Users: %s", utils.ToJSON(state.Resources.Users))

		if idpUsersResult.HashCode == state.Resources.Users.HashCode {
			log.Info("provider users and state users are the same, nothing to do with users")

			totalUsersResult = state.Resources.Users
		} else {
			log.Info("provider users and state users are diferent")

			log.WithFields(log.Fields{
				"idp":   idpUsersResult.Items,
				"state": state.Resources.Users.Items,
			}).Info("reconciling users")
			usersCreate, usersUpdate, usersEqual, usersDelete := usersOperations(idpUsersResult, &state.Resources.Users)

			usersCreated, usersUpdated, err := reconcilingUsers(ss.ctx, ss.scim, usersCreate, usersUpdate, usersDelete)
			if err != nil {
				return fmt.Errorf("error reconciling users: %w", err)
			}

			// usersCreated + usersUpdated + usersEqual = users total
			totalUsersResult = mergeUsersResult(usersCreated, usersUpdated, usersEqual)
		}

		log.Tracef("state.Resources.GroupsMembers: %s", utils.ToJSON(state.Resources.GroupsMembers))

		if idpGroupsMembersResult.HashCode == state.Resources.GroupsMembers.HashCode {
			log.Info("provider groups-members and state groups-members are the same, nothing to do with groups-members")

			totalGroupsMembersResult = state.Resources.GroupsMembers
		} else {
			log.Info("provider groups-members and state groups-members are diferent")

			log.WithFields(log.Fields{
				"idp":   idpGroupsMembersResult.Items,
				"state": state.Resources.GroupsMembers.Items,
			}).Info("reconciling groups members")
			membersCreate, membersEqual, membersDelete := membersOperations(idpGroupsMembersResult, &state.Resources.GroupsMembers)

			membersCreated, err := reconcilingGroupsMembers(ss.ctx, ss.scim, membersCreate, membersDelete)
			if err != nil {
				return fmt.Errorf("error reconciling groups members: %w", err)
			}

			// membersCreate + membersEqual = members total
			totalGroupsMembersResult = mergeGroupsMembersResult(membersCreated, membersEqual)

		}
	}

	// after be sure all the SCIM side is aligned with the Identity Provider side
	// we can update the state with the identity provider data
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
		"members":  totalGroupsMembersResult.Items,
	}).Info("storing the new state")

	// TODO: avoid this step using a cmd flag, could be a nice feature
	if err := ss.repo.SetState(ss.ctx, newState); err != nil {
		return fmt.Errorf("error storing the state: %w", err)
	}

	log.Tracef("state data: %s", utils.ToJSON(newState))
	log.Info("sync completed")
	return nil
}

// SyncGroupsAndUsers this method is used to sync the usersm groups and their members from the identity provider to the SCIM
func (ss *SyncService) SyncGroupsAndUsers() error {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	return errors.New("not implemented")
}

// reconcilingGroups receives lists of groups to create, update, equals and delete in the SCIM service
// returns the lists of groups created and updated in the SCIM service with the Ids of these groups.
func reconcilingGroups(ctx context.Context, scim SCIMService, create *model.GroupsResult, update *model.GroupsResult, delete *model.GroupsResult) (created *model.GroupsResult, updated *model.GroupsResult, e error) {
	var err error

	if create.Items == 0 {
		log.Info("no groups to be create")
		created = &model.GroupsResult{Items: 0, Resources: []model.Group{}}
	} else {
		log.WithField("quantity", create.Items).Warn("creating groups")
		created, err = scim.CreateGroups(ctx, create)
		if err != nil {
			return nil, nil, fmt.Errorf("error creating groups in SCIM Provider: %w", err)
		}

	}

	if update.Items == 0 {
		log.Info("no groups to be updated")
		updated = &model.GroupsResult{Items: 0, Resources: []model.Group{}}
	} else {
		log.WithField("quantity", update.Items).Warn("updating groups")
		updated, err = scim.UpdateGroups(ctx, update)
		if err != nil {
			return nil, nil, fmt.Errorf("error updating groups in SCIM Provider: %w", err)
		}

	}

	if delete.Items == 0 {
		log.Info("no groups to be deleted")
	} else {
		log.WithField("quantity", delete.Items).Warn("deleting groups")
		if err := scim.DeleteGroups(ctx, delete); err != nil {
			return nil, nil, fmt.Errorf("error deleting groups in SCIM Provider: %w", err)
		}

	}

	return
}

// reconcilingUsers receives lists of users to create, update and delete in the SCIM service
// returns the lists of users created and updated in the SCIM service with the Ids of these users in the SCIM service
func reconcilingUsers(ctx context.Context, scim SCIMService, create *model.UsersResult, update *model.UsersResult, delete *model.UsersResult) (created *model.UsersResult, updated *model.UsersResult, e error) {
	var err error

	if create.Items == 0 {
		log.Info("no users to be created")
		created = &model.UsersResult{Items: 0, Resources: []model.User{}}
	} else {
		log.WithField("quantity", create.Items).Warn("creating users")
		created, err = scim.CreateUsers(ctx, create)
		if err != nil {
			return nil, nil, fmt.Errorf("error creating users in SCIM Provider: %w", err)
		}
	}

	if update.Items == 0 {
		log.Info("no users to be updated")
		updated = &model.UsersResult{Items: 0, Resources: []model.User{}}
	} else {
		log.WithField("quantity", update.Items).Warn("updating users")
		updated, err = scim.UpdateUsers(ctx, update)
		if err != nil {
			return nil, nil, fmt.Errorf("error updating users in SCIM Provider: %w", err)
		}
	}

	if delete.Items == 0 {
		log.Info("no users to be deleted")
	} else {
		log.WithField("quantity", delete.Items).Warn("deleting users")
		if err := scim.DeleteUsers(ctx, delete); err != nil {
			return nil, nil, fmt.Errorf("error deleting users in SCIM Provider: %w", err)
		}
	}

	return
}

// reconcilingGroupsMembers
func reconcilingGroupsMembers(ctx context.Context, scim SCIMService, create *model.GroupsMembersResult, delete *model.GroupsMembersResult) (created *model.GroupsMembersResult, e error) {
	var err error

	if create.Items == 0 {
		log.Info("no users to be joined to groups")
		created = &model.GroupsMembersResult{Items: 0, Resources: []model.GroupMembers{}}
	} else {
		log.WithField("quantity", create.Items).Warn("joining users to groups")
		created, err = scim.CreateGroupsMembers(ctx, create)
		if err != nil {
			return nil, fmt.Errorf("error creating groups members in SCIM Provider: %w", err)
		}
	}

	if delete.Items == 0 {
		log.Info("no users to be removed from groups")
	} else {
		log.WithField("quantity", delete.Items).Warn("removing users to groups")
		if err := scim.DeleteGroupsMembers(ctx, delete); err != nil {
			return nil, fmt.Errorf("error removing users from groups in SCIM Provider: %w", err)
		}
	}

	return
}
