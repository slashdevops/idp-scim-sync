package core

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/slashdevops/idp-scim-sync/internal/model"
)

var (
	ErrIdentiyProviderServiceNil = errors.New("sync: identity provider service cannot be nil")
	ErrSCIMServiceNil            = errors.New("sync: SCIM service cannot be nil")
	ErrStateRepositoryNil        = errors.New("sync: state repository cannot be nil")
)

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

func (ss *SyncService) SyncGroupsAndTheirMembers() error {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	// get data from the identity provider
	idpUsersResult, idpGroupsResult, idpGroupsUsersResult, err := getIdentityProviderData(ss.ctx, ss.prov, ss.provGroupsFilter)
	if err != nil {
		return fmt.Errorf("sync: error getting data from the identity provider: %w", err)
	}

	if idpUsersResult.Items == 0 {
		log.Warn("sync: there are no users in the identity provider")
	}

	if idpGroupsResult.Items == 0 {
		log.Warnf("sync: there are no groups in the identity provider that match with this filter: %s", ss.provGroupsFilter)
	}

	if idpGroupsUsersResult.Items == 0 {
		log.Warn("sync: there are no group users in the identity provider")
	}

	// get the state data from the repository
	state, err := ss.repo.GetState(ss.ctx)
	if err != nil {
		return fmt.Errorf("sync: error getting state data from the repository: %w", err)
	}

	// these variables are used to store the data that will be used to create or delete users and groups in SCIM
	// the differents between the data in the identity provider and these is that these have already set the SCIMID
	// after the creation of the element in SCIM
	var createdGroupsResult model.GroupsResult
	var createdUsersResult model.UsersResult
	var createdGroupsUsersResult model.GroupsUsersResult

	// first time syncing
	if state.LastSync == "" {
		// Check SCIM side to see if there are elelemnts to be
		// reconciled. Bassically check if SCIM is not clean before the first sync
		// and we need to reconcile the SCIM side with the identity provider side.
		// In case of migration from a different tool and we want to keep the state
		// of the users and groups in the SCIM side, just no recreate, keep the existing ones.
		// NOTE: This only work if the user id is the same in both sides, in our case is rthe email address.

		log.Info("state without lastsync time, first time syncing")
		log.Warn("reconciling the SCIM data with the Identity Provider data, the first syncing")

		scimUsersResult, scimGroupsResult, scimGroupsUsersResult, err := getSCIMData(ss.ctx, ss.scim)
		if err != nil {
			return fmt.Errorf("sync: error getting data from the SCIM service: %w", err)
		}

		log.WithFields(log.Fields{
			"idp_quantity":  idpGroupsResult.Items,
			"scim_quantity": scimGroupsResult.Items,
		}).Info("starting reconciling groups")
		groupsCreate, groupsUpdate, _, groupsDelete := groupsOperations(idpGroupsResult, scimGroupsResult)

		rgrc, rgru, err := reconcilingSCIMGroups(ss.ctx, ss.scim, groupsCreate, groupsUpdate, groupsDelete)
		if err != nil {
			return fmt.Errorf("sync: error reconciling groups: %w", err)
		}

		// merge in only one data structure the groups created and updated who has the SCIMID
		createdGroupsResult = mergeGroupsResult(rgrc, rgru)

		log.WithFields(log.Fields{
			"idp_quantity":  idpUsersResult.Items,
			"scim_quantity": scimUsersResult.Items,
		}).Info("starting reconciling users")
		usersCreate, usersUpdate, _, usersDelete := usersOperations(idpUsersResult, scimUsersResult)

		rurc, ruru, err := reconcilingSCIMUsers(ss.ctx, ss.scim, usersCreate, usersUpdate, usersDelete)
		if err != nil {
			return fmt.Errorf("sync: error reconciling users: %w", err)
		}

		// merge in only one data structure the users created and updated who has the SCIMID
		createdUsersResult = mergeUsersResult(rurc, ruru)

		log.WithFields(log.Fields{
			"idp_quantity":  idpGroupsUsersResult.Items,
			"scim_quantity": scimGroupsUsersResult.Items,
		}).Info("starting reconciling groups members")
		ugCreate, _, ugDelete := groupsUsersOperations(idpGroupsUsersResult, scimGroupsUsersResult)

		createdGroupsUsersResult = *ugCreate

		if err := reconcilingSCIMGroupsUsers(ss.ctx, ss.scim, ugCreate, ugDelete); err != nil {
			return fmt.Errorf("sync: error reconciling groups users: %w", err)
		}

	} else { // This is not the first time syncing
		log.WithField("lastsync", state.LastSync).Info("state with lastsync time, it is not first time syncing")

		if idpGroupsResult.HashCode == state.Resources.Groups.HashCode {
			log.Info("provider groups and state groups are the same, nothing to do with groups")
		} else {
			log.Info("provider groups and state groups are diferent")
			// now here we have the google fresh data and the last sync data state
			// we need to compare the data and decide what to do
			// see differences between the two data sets

			log.WithFields(log.Fields{
				"idp_quantity":   idpGroupsResult.Items,
				"state_quantity": &state.Resources.Groups.Items,
			}).Info("starting reconciling groups")
			groupsCreate, groupsUpdate, _, groupsDelete := groupsOperations(idpGroupsResult, &state.Resources.Groups)

			rgrc, rgru, err := reconcilingSCIMGroups(ss.ctx, ss.scim, groupsCreate, groupsUpdate, groupsDelete)
			if err != nil {
				return fmt.Errorf("sync: error reconciling groups: %w", err)
			}

			// merge in only one data structure the groups created and updated who has the SCIMID
			createdGroupsResult = mergeGroupsResult(rgrc, rgru)
		}

		if idpUsersResult.HashCode == state.Resources.Users.HashCode {
			log.Info("provider users and state users are the same, nothing to do with users")
		} else {
			log.Info("provider users and state users are diferent")

			log.WithFields(log.Fields{
				"idp_quantity":   idpUsersResult.Items,
				"state_quantity": &state.Resources.Users.Items,
			}).Info("starting reconciling users")
			usersCreate, usersUpdate, _, usersDelete := usersOperations(idpUsersResult, &state.Resources.Users)

			rurc, ruru, err := reconcilingSCIMUsers(ss.ctx, ss.scim, usersCreate, usersUpdate, usersDelete)
			if err != nil {
				return fmt.Errorf("sync: error reconciling users: %w", err)
			}

			// merge in only one data structure the users created and updated who has the SCIMID
			createdUsersResult = mergeUsersResult(rurc, ruru)
		}

		if idpGroupsUsersResult.HashCode == state.Resources.GroupsUsers.HashCode {
			log.Info("provider groups-members and state groups-members are the same, nothing to do with groups-members")
		} else {
			log.Info("provider groups-members and state groups-members are diferent")

			log.WithFields(log.Fields{
				"idp_quantity":   idpGroupsUsersResult.Items,
				"state_quantity": &state.Resources.GroupsUsers.Items,
			}).Info("starting reconciling groups members")
			ugCreate, _, ugDelete := groupsUsersOperations(idpGroupsUsersResult, &state.Resources.GroupsUsers)

			createdGroupsUsersResult = *ugCreate

			if err := reconcilingSCIMGroupsUsers(ss.ctx, ss.scim, ugCreate, ugDelete); err != nil {
				return fmt.Errorf("sync: error reconciling groups users: %w", err)
			}
		}

	}

	// after be sure all the SCIM side is aligned with the Identity Provider side
	// we can update the state with the identity provider data
	newState := &model.State{
		SchemaVersion: "1.0.0",
		CodeVersion:   "0.0.1",
		LastSync:      time.Now().Format(time.RFC3339),
		Resources: model.StateResources{
			Groups:      createdGroupsResult,
			Users:       createdUsersResult,
			GroupsUsers: createdGroupsUsersResult,
		},
	}

	log.WithFields(log.Fields{
		"groups_quantity": createdGroupsResult.Items,
		"users_quantity":  createdUsersResult.Items,
	}).Info("setting the new state")
	if err := ss.repo.SetState(ss.ctx, newState); err != nil {
		return fmt.Errorf("sync: error saving state: %w", err)
	}

	return nil
}

func (ss *SyncService) SyncGroupsAndUsers() error {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	return errors.New("not implemented")
}

// getIdentityProviderData return the users, groups and groups and their users from Identity Provider Service
func getIdentityProviderData(ctx context.Context, ip IdentityProviderService, groupFilter []string) (*model.UsersResult, *model.GroupsResult, *model.GroupsUsersResult, error) {
	groupsResult, err := ip.GetGroups(ctx, groupFilter)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("sync: error getting groups from the identity provider: %w", err)
	}

	usersResult, groupsUsersResult, err := ip.GetUsersAndGroupsUsers(ctx, groupsResult)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("sync: error getting users and groups and their users: %w", err)
	}

	return usersResult, groupsResult, groupsUsersResult, nil
}

// getSCIMData return the users, groups and groups and their users from SCIM Service
func getSCIMData(ctx context.Context, scim SCIMService) (*model.UsersResult, *model.GroupsResult, *model.GroupsUsersResult, error) {
	groupsResult, err := scim.GetGroups(ctx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("sync: error getting groups from the SCIM provider: %w", err)
	}

	usersResult, groupsUsersResult, err := scim.GetUsersAndGroupsUsers(ctx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("sync: error getting users and groups and their users from SCIM provider: %w", err)
	}

	return usersResult, groupsResult, groupsUsersResult, nil
}

// reconcilingSCIMGroups receives lists of groups to create, update and delete in the SCIM service
// returns the lists of groups created and updated in the SCIM service with the Ids of these groups in the SCIM service
func reconcilingSCIMGroups(ctx context.Context, scim SCIMService, create *model.GroupsResult, update *model.GroupsResult, delete *model.GroupsResult) (created *model.GroupsResult, updated *model.GroupsResult, e error) {
	if create.Items == 0 {
		log.Info("no groups to be create")
	} else {
		log.WithField("quantity", create.Items).Info("creating groups")
	}

	created, err := scim.CreateGroups(ctx, create)
	if err != nil {
		return nil, nil, fmt.Errorf("sync: error creating groups in SCIM Provider: %w", err)
	}

	if update.Items == 0 {
		log.Info("no groups to be updated")
	} else {
		log.WithField("quantity", update.Items).Info("updating groups")
	}

	updated, err = scim.UpdateGroups(ctx, update)
	if err != nil {
		return nil, nil, fmt.Errorf("sync: error updating groups in SCIM Provider: %w", err)
	}

	if delete.Items == 0 {
		log.Info("no groups to be deleted")
	} else {
		log.WithField("quantity", delete.Items).Info("deleting groups")
	}

	if err = scim.DeleteGroups(ctx, delete); err != nil {
		return nil, nil, fmt.Errorf("sync: error deleting groups in SCIM Provider: %w", err)
	}

	return
}

// reconcilingSCIMUsers receives lists of users to create, update and delete in the SCIM service
// returns the lists of users created and updated in the SCIM service with the Ids of these users in the SCIM service
func reconcilingSCIMUsers(ctx context.Context, scim SCIMService, create *model.UsersResult, update *model.UsersResult, delete *model.UsersResult) (created *model.UsersResult, updated *model.UsersResult, e error) {
	if create.Items == 0 {
		log.Info("no users to be created")
	} else {
		log.WithField("quantity", create.Items).Info("creating users")
	}

	created, err := scim.CreateUsers(ctx, create)
	if err != nil {
		return nil, nil, fmt.Errorf("sync: error creating users in SCIM Provider: %w", err)
	}

	if update.Items == 0 {
		log.Info("no users to be updated")
	} else {
		log.WithField("quantity", update.Items).Info("updating users")
	}

	updated, err = scim.UpdateUsers(ctx, update)
	if err != nil {
		return nil, nil, fmt.Errorf("sync: error updating users in SCIM Provider: %w", err)
	}

	if delete.Items == 0 {
		log.Info("no users to be deleted")
	} else {
		log.WithField("quantity", delete.Items).Info("deleting users")
	}

	if err := scim.DeleteUsers(ctx, delete); err != nil {
		return nil, nil, fmt.Errorf("sync: error deleting users in SCIM Provider: %w", err)
	}

	return
}

// reconcilingSCIMGroupsUsers
func reconcilingSCIMGroupsUsers(ctx context.Context, scim SCIMService, create *model.GroupsUsersResult, delete *model.GroupsUsersResult) error {
	if create.Items == 0 {
		log.Info("no users to be joined to groups")
	} else {
		log.WithField("quantity", create.Items).Info("joining users to groups")
		if err := scim.CreateGroupsMembers(ctx, create); err != nil {
			return fmt.Errorf("sync: error creating groups members in SCIM Provider: %w", err)
		}
	}

	if delete.Items == 0 {
		log.Info("no users to be removed from groups")
	} else {
		log.WithField("quantity", delete.Items).Info("removing users to groups")
		if err := scim.DeleteGroupsMembers(ctx, delete); err != nil {
			return fmt.Errorf("sync: error removing users from groups in SCIM Provider: %w", err)
		}
	}

	return nil
}
