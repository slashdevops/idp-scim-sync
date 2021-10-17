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
	ErrProviderServiceNil = errors.New("identity provider service cannot be nil")
	ErrSCIMServiceNil     = errors.New("SCIM service cannot be nil")
	ErrRepositoryNil      = errors.New("repository cannot be nil")
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
		return nil, ErrProviderServiceNil
	}
	if scim == nil {
		return nil, ErrSCIMServiceNil
	}
	if repo == nil {
		return nil, ErrRepositoryNil
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

	log := log.WithFields(log.Fields{
		"filter": ss.provGroupsFilter,
	})

	// get data from the identity provider
	pUsersResult, pGroupsResult, pGroupsUsersResult, err := getIdentityProviderData(ss.ctx, ss.prov, ss.provGroupsFilter)
	if err != nil {
		return fmt.Errorf("error getting data from the identity provider: %w", err)
	}

	// get the state metadata from the repository
	state, err := ss.repo.GetState(ss.ctx)
	if err != nil {
		return fmt.Errorf("error getting state from the repository: %w", err)
	}

	// first time syncing
	if state.LastSync == "" {

		log.Info("state without lastsync time, first time syncing")
		log.Warn("reconciling the SCIM data with the Identity Provider data, the first syncing")

		// Check SCIM side to see if there are elelemnts to be
		// reconciled. I mean SCIM is not clean before the first sync
		// and we need to reconcile the SCIM side with the identity provider side.
		// In case of migration from a different tool and we want to keep the state
		// of the users and groups in the SCIM side, just no recreate, keep the existing ones.
		//
		// During this check the groups could be empty, but we need to check if users exist
		// even if there are not groups. Becuase the users, groups their members in the scim needs to be
		// controlled by the sync process.
		sUsersResult, sGroupsResult, sGroupsUsersResult, err := getSCIMData(ss.ctx, ss.scim)
		if err != nil {
			return fmt.Errorf("error getting data from the SCIM service: %w", err)
		}

		log.Info("starting reconciling groups")
		gCreate, gUpdate, _, gDelete := groupsOperations(pGroupsResult, sGroupsResult)

		rgrc, rgru, rgrd, err := reconcilingSCIMGroups(ss.ctx, ss.scim, gCreate, gUpdate, gDelete)
		if err != nil {
			return fmt.Errorf("error reconciling groups: %w", err)
		}

		log.Info("starting reconciling users")
		uCreate, uUpdate, _, uDelete := usersOperations(pUsersResult, sUsersResult)

		rurc, ruru, rurd, err := reconcilingSCIMUsers(ss.ctx, ss.scim, uCreate, uUpdate, uDelete)
		if err != nil {
			return fmt.Errorf("error reconciling users: %w", err)
		}

		log.Info("starting reconciling groups members")
		ugCreate, _, ugDelete := groupsUsersOperations(pGroupsUsersResult, sGroupsUsersResult)

		rgurc, rgurd, err := reconcilingSCIMGroupsUsers(ss.ctx, ss.scim, ugCreate, ugDelete)
		if err != nil {
			return fmt.Errorf("error reconciling groups users: %w", err)
		}

		// WIP
		_ = rgrc
		_ = rgrd
		_ = rurc
		_ = ruru
		_ = rurd
		_ = rgru
		_ = rgurd
		_ = rgurc
		_ = rgurd

	} else { // This is not the first time syncing
		log.WithField("lastsync", state.LastSync).Info("state with lastsync time, it is not first time syncing")

		if pGroupsResult.HashCode == state.Resources.Groups.HashCode {
			log.Info("provider groups and state groups are the same, nothing to do with groups")
		}

		log.Info("provider groups and state groups are diferent")
		// now here we have the google fresh data and the last sync data state
		// we need to compare the data and decide what to do
		// see differences between the two data sets

		// syncing groups
		gCreate, gUpdate, _, gDelete := groupsOperations(pGroupsResult, &state.Resources.Groups)

		rgrc, rgru, rgrd, err := reconcilingSCIMGroups(ss.ctx, ss.scim, gCreate, gUpdate, gDelete)
		if err != nil {
			return fmt.Errorf("error reconciling groups: %w", err)
		}

		if pUsersResult.HashCode == state.Resources.Users.HashCode {
			log.Info("provider users and state users are the same, nothing to do with users")
		}

		log.Info("provider users and state users are diferent")

		// syncing users
		uCreate, uUpdate, _, uDelete := usersOperations(pUsersResult, &state.Resources.Users)

		rurc, ruru, rurd, err := reconcilingSCIMUsers(ss.ctx, ss.scim, uCreate, uUpdate, uDelete)
		if err != nil {
			return fmt.Errorf("error reconciling users: %w", err)
		}

		if pGroupsUsersResult.HashCode == state.Resources.GroupsUsers.HashCode {
			log.Info("provider groups-members and state groups-members are the same, nothing to do with groups-members")
		}
		log.Info("provider groups-members and state groups-members are diferent")

		// syncing groups-users --> groups members
		ugCreate, _, ugDelete := groupsUsersOperations(pGroupsUsersResult, &state.Resources.GroupsUsers)

		rgurc, rgurd, err := reconcilingSCIMGroupsUsers(ss.ctx, ss.scim, ugCreate, ugDelete)
		if err != nil {
			return fmt.Errorf("error reconciling groups users: %w", err)
		}

		// WIP
		_ = rgrc
		_ = rgrd
		_ = rurc
		_ = ruru
		_ = rurd
		_ = rgru
		_ = rgurd
		_ = rgurc
		_ = rgurd

	}

	// after be sure all the SCIM side is aligned with the Identity Provider side
	// we can update the state with the identity provider data
	newState := &model.State{
		SchemaVersion: "1.0.0",
		CodeVersion:   "0.0.1",
		LastSync:      time.Now().Format(time.RFC3339),
		Resources: model.StateResources{
			Groups:      *pGroupsResult,
			Users:       *pUsersResult,
			GroupsUsers: *pGroupsUsersResult,
		},
	}

	if err := ss.repo.SaveState(ss.ctx, newState); err != nil {
		return fmt.Errorf("error saving state: %w", err)
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
	// retrive data from the identity provider
	// always theses steps are necessary
	groupsResult, err := ip.GetGroups(ctx, groupFilter)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error getting groups from the identity provider: %w", err)
	}

	if groupsResult.Items == 0 {
		log.Warnf("there are no groups in the identity provider that match with this filter: %s", groupFilter)
	}

	usersResult, groupsUsersResult, err := ip.GetUsersAndGroupsUsers(ctx, groupsResult)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error getting users and groups and their users: %w", err)
	}

	if groupsUsersResult.Items == 0 {
		log.Warn("there are no group users in the identity provider")
	}

	if usersResult.Items == 0 {
		log.Warn("there are no users in the identity provider")
	}

	return usersResult, groupsResult, groupsUsersResult, nil
}

// getSCIMData return the users, groups and groups and their users from SCIM Service
func getSCIMData(ctx context.Context, scim SCIMService) (*model.UsersResult, *model.GroupsResult, *model.GroupsUsersResult, error) {
	groupsResult, err := scim.GetGroups(ctx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error getting groups from the SCIM service: %w", err)
	}

	if groupsResult.Items == 0 {
		log.Warn("there are no groups in the SCIM Service")
	}

	usersResult, groupsUsersResult, err := scim.GetUsersAndGroupsUsers(ctx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error getting users and groups and their users from SCIM provider: %w", err)
	}

	if usersResult.Items == 0 {
		log.Warn("no SCIM users to be reconciling")
	}

	if groupsUsersResult.Items == 0 {
		log.Warn("no SCIM groups and members to be reconciling")
	}

	return usersResult, groupsResult, groupsUsersResult, nil
}

// reconcilingSCIMGroups
func reconcilingSCIMGroups(ctx context.Context, scim SCIMService, create *model.GroupsResult, update *model.GroupsResult, delete *model.GroupsResult) (c *model.GroupsResult, u *model.GroupsResult, d *model.GroupsResult, e error) {
	if create.Items == 0 {
		log.Info("no groups to be create")
	}

	log.WithField("quantity", create.Items).Info("creating groups")
	c, err := scim.CreateGroups(ctx, create)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error creating groups in SCIM Provider: %w", err)
	}

	if update.Items == 0 {
		log.Info("no groups to be updated")
	}

	log.WithField("quantity", update.Items).Info("updating groups")
	u, err = scim.UpdateGroups(ctx, update)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error updating groups in SCIM Provider: %w", err)
	}

	if delete.Items == 0 {
		log.Info("no groups to be deleted")
	}

	log.WithField("quantity", delete.Items).Info("deleting groups")
	d, err = scim.DeleteGroups(ctx, delete)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error deleting groups in SCIM Provider: %w", err)
	}

	return
}

// reconcilingSCIMUsers
func reconcilingSCIMUsers(ctx context.Context, scim SCIMService, create *model.UsersResult, update *model.UsersResult, delete *model.UsersResult) (c *model.UsersResult, u *model.UsersResult, d *model.UsersResult, e error) {
	if create.Items == 0 {
		log.Info("no users to be created")
	}

	log.WithField("quantity", create.Items).Info("creating users")
	c, err := scim.CreateUsers(ctx, create)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error creating users in SCIM Provider: %w", err)
	}

	if update.Items == 0 {
		log.Info("no users to be updated")
	}

	log.WithField("quantity", update.Items).Info("updating users")
	u, err = scim.UpdateUsers(ctx, update)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error updating users in SCIM Provider: %w", err)
	}

	if delete.Items == 0 {
		log.Info("no users to be deleted")
	}

	log.WithField("quantity", delete.Items).Info("deleting users")
	d, err = scim.DeleteUsers(ctx, delete)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error deleting users in SCIM Provider: %w", err)
	}

	return
}

// reconcilingSCIMGroupsUsers
func reconcilingSCIMGroupsUsers(ctx context.Context, scim SCIMService, create *model.GroupsUsersResult, delete *model.GroupsUsersResult) (c *model.GroupsUsersResult, d *model.GroupsUsersResult, e error) {
	if create.Items == 0 {
		log.Info("no users to be joined to groups")
	}

	log.WithField("quantity", create.Items).Info("joining users to groups")
	c, err := scim.CreateMembers(ctx, create)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating groups members in SCIM Provider: %w", err)
	}

	if delete.Items == 0 {
		log.Info("no usrs to be removed from groups")
	}

	log.WithField("quantity", delete.Items).Info("removing users to groups")
	d, err = scim.DeleteMembers(ctx, delete)
	if err != nil {
		return nil, nil, fmt.Errorf("error removing users from groups in SCIM Provider: %w", err)
	}

	return
}
