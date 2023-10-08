package core

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	log "github.com/sirupsen/logrus"
	"github.com/slashdevops/idp-scim-sync/internal/model"
	"github.com/slashdevops/idp-scim-sync/internal/repository"
	"github.com/slashdevops/idp-scim-sync/internal/version"
)

var (
	// ErrIdentityProviderServiceNil is returned when the Identity Provider Service is nil
	ErrIdentityProviderServiceNil = errors.New("identity provider service cannot be nil")

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
		return nil, ErrIdentityProviderServiceNil
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
	}).Info("getting identity provider data")

	idpGroupsResult, err := ss.prov.GetGroups(ctx, ss.provGroupsFilter)
	if err != nil {
		return fmt.Errorf("error getting groups from the identity provider: %w", err)
	}

	log.WithFields(
		log.Fields{
			"group_filter": ss.provGroupsFilter,
			"groups":       idpGroupsResult.Items,
		}).Info("groups retrieved from the identity provider for syncing that match the filter")

	idpGroupsMembersResult, err := ss.prov.GetGroupsMembers(ctx, idpGroupsResult)
	if err != nil {
		return fmt.Errorf("error getting groups members: %w", err)
	}

	log.WithFields(
		log.Fields{
			"group_filter":   ss.provGroupsFilter,
			"groups_members": idpGroupsMembersResult.Items,
		}).Info("groups members retrieved from the identity provider for syncing that match the filter")

	log.WithFields(log.Fields{
		"group_filter": ss.provGroupsFilter,
	}).Info("getting users (using groups members) from the identity provider")

	idpUsersResult, err := ss.prov.GetUsersByGroupsMembers(ctx, idpGroupsMembersResult)
	if err != nil {
		return fmt.Errorf("error getting users from the identity provider: %w", err)
	}

	log.WithFields(
		log.Fields{
			"group_filter": ss.provGroupsFilter,
			"users":        idpUsersResult.Items,
		}).Info("users retrieved from the identity provider for syncing that match the filter")

	log.Info("getting state data")
	state, err := ss.repo.GetState(ctx)
	if err != nil {
		var nsk *types.NoSuchKey
		var StateFileEmpty *repository.ErrStateFileEmpty

		if errors.As(err, &nsk) || errors.As(err, &StateFileEmpty) {
			log.Warn("no state file found in the state repository, creating a new one")
			state = model.StateBuilder().Build()
		} else {
			return fmt.Errorf("error getting state data from the repository: %w", err)
		}
	}

	var (
		totalGroupsResult        *model.GroupsResult
		totalUsersResult         *model.UsersResult
		totalGroupsMembersResult *model.GroupsMembersResult
	)

	// first time syncing
	if state.LastSync == "" {
		// Check SCIM side to see if there are elements to be reconciled.
		// Basically, checks if SCIM is not clean before the first sync
		// and we need to reconcile the SCIM side with the identity provider side.
		// In case of migration from a different tool and we want to keep the state
		// of the users and groups in the SCIM side, just no recreation, keep the existing ones when:
		// - Groups names are equals on both sides, update only the external id (coming from the identity provider)
		// - Users emails are equals on both sides, update only the external id (coming from the identity provider)
		log.Warn("syncing from scim service, first time syncing")
		totalGroupsResult, totalUsersResult, totalGroupsMembersResult, err = scimSync(
			ctx,
			ss.scim,
			idpGroupsResult,
			idpUsersResult,
			idpGroupsMembersResult,
		)
		if err != nil {
			return fmt.Errorf("error doing the first sync: %w", err)
		}
	} else {
		log.Warn("syncing from state, it's not the first time syncing")
		totalGroupsResult, totalUsersResult, totalGroupsMembersResult, err = stateSync(
			ctx,
			state,
			ss.scim,
			idpGroupsResult,
			idpUsersResult,
			idpGroupsMembersResult,
		)
		if err != nil {
			return fmt.Errorf("error syncing state: %w", err)
		}
	}

	// after be sure all the SCIM side is aligned with the identity provider side
	// we can update the state with the last data coming from the reconciliation
	newState := model.StateBuilder().
		WithCodeVersion(version.Version).
		WithLastSync(time.Now().Format(time.RFC3339)).
		WithGroups(totalGroupsResult).
		WithUsers(totalUsersResult).
		WithGroupsMembers(totalGroupsMembersResult).
		Build()

	log.WithFields(log.Fields{
		"lastSync": newState.LastSync,
		"groups":   totalGroupsResult.Items,
		"users":    totalUsersResult.Items,
	}).Info("storing the new state")

	if err := ss.repo.SetState(ctx, newState); err != nil {
		return fmt.Errorf("error storing the state: %w", err)
	}

	log.WithFields(log.Fields{
		"date": time.Now().Format(time.RFC3339),
	}).Info("sync completed")
	return nil
}
