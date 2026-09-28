package core

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/slashdevops/idp-scim-sync/internal/model"
)

var (
	// ErrCreateGroupsResultNil is returned when the create *model.GroupsResult argument is nil
	ErrCreateGroupsResultNil = errors.New("create Groups Result is nil")

	// ErrUpdateGroupsResultNil is returned when the update *model.GroupsResult argument is nil
	ErrUpdateGroupsResultNil = errors.New("update Groups Result is nil")

	// ErrDeleteGroupsResultNil is returned when the delete *model.GroupsResult argument is nil
	ErrDeleteGroupsResultNil = errors.New("delete Groups Result is nil")

	// ErrCreateUsersResultNil is returned when the create *model.UsersResult argument is nil
	ErrCreateUsersResultNil = errors.New("create Users Result is nil")

	// ErrUpdateUsersResultNil is returned when the update *model.UsersResult argument is nil
	ErrUpdateUsersResultNil = errors.New("update Users Result is nil")

	// ErrDeleteUsersResultNil is returned when the delete *model.UsersResult argument is nil
	ErrDeleteUsersResultNil = errors.New("remove Users Result is nil")

	// ErrCreateGroupsMembersResultNil is returned when the SCIM *model.GroupsMembersResult argument is nil
	ErrCreateGroupsMembersResultNil = errors.New("create Groups Members Result is nil")

	// ErrDeleteGroupsMembersResultNil is returned when the SCIM *model.GroupsMembersResult argument is nil
	ErrDeleteGroupsMembersResultNil = errors.New("remove Groups Members Result is nil")
)

// reconcilingGroups creates, update and removes from groups in SCIM service
// returns the lists of groups created and updated in the SCIM provider
// with the ids of these groups.
func reconcilingGroups(
	ctx context.Context,
	scim SCIMService,
	create, update, remove *model.GroupsResult,
) (*model.GroupsResult, *model.GroupsResult, error) {
	if scim == nil {
		return nil, nil, ErrSCIMServiceNil
	}
	if create == nil {
		return nil, nil, ErrCreateGroupsResultNil
	}
	if update == nil {
		return nil, nil, ErrUpdateGroupsResultNil
	}
	if remove == nil {
		return nil, nil, ErrDeleteGroupsResultNil
	}

	var (
		created, updated *model.GroupsResult
		err              error
	)

	if create.Items == 0 {
		slog.Info("no groups to be create")
		created = model.GroupsResultBuilder().Build()
	} else {
		slog.Warn("creating groups", "groups", create.Items)
		created, err = scim.CreateGroups(ctx, create)
		if err != nil {
			return nil, nil, fmt.Errorf("error creating groups from SCIM provider: %w", err)
		}
	}

	if update.Items == 0 {
		slog.Info("no groups to be updated")
		updated = model.GroupsResultBuilder().Build()
	} else {
		slog.Warn("updating groups", "groups", update.Items)
		updated, err = scim.UpdateGroups(ctx, update)
		if err != nil {
			return nil, nil, fmt.Errorf("error updating groups from SCIM provider: %w", err)
		}
	}

	if remove.Items == 0 {
		slog.Info("no groups to be deleted")
	} else {
		slog.Warn("deleting groups", "groups", remove.Items)
		if err := scim.DeleteGroups(ctx, remove); err != nil {
			return nil, nil, fmt.Errorf("error deleting groups from SCIM provider: %w", err)
		}
	}

	return created, updated, nil
}

// reconcilingUsers creates, updates and removes users in SCIM provider
// returns the lists of users created and updated in the SCIM provider
// with the ids of these users.
func reconcilingUsers(
	ctx context.Context,
	scim SCIMService,
	create, update,
	remove *model.UsersResult,
) (*model.UsersResult, *model.UsersResult, error) {
	if scim == nil {
		return nil, nil, ErrSCIMServiceNil
	}
	if create == nil {
		return nil, nil, ErrCreateUsersResultNil
	}
	if update == nil {
		return nil, nil, ErrUpdateUsersResultNil
	}
	if remove == nil {
		return nil, nil, ErrDeleteUsersResultNil
	}

	var (
		created, updated *model.UsersResult
		err              error
	)

	if create.Items == 0 {
		slog.Info("no users to be created")
		created = model.UsersResultBuilder().Build()
	} else {
		slog.Warn("creating users", "users", create.Items)
		created, err = scim.CreateUsers(ctx, create)
		if err != nil {
			return nil, nil, fmt.Errorf("error creating users from SCIM provider: %w", err)
		}
	}

	if update.Items == 0 {
		slog.Info("no users to be updated")
		updated = model.UsersResultBuilder().Build()
	} else {
		slog.Warn("updating users", "users", update.Items)
		updated, err = scim.UpdateUsers(ctx, update)
		if err != nil {
			return nil, nil, fmt.Errorf("error updating users from SCIM provider: %w", err)
		}
	}

	if remove.Items == 0 {
		slog.Info("no users to be removed")
	} else {
		slog.Warn("deleting users", "users", remove.Items)
		if err := scim.DeleteUsers(ctx, remove); err != nil {
			return nil, nil, fmt.Errorf("error deleting users from SCIM provider: %w", err)
		}
	}

	return created, updated, nil
}

// reconcilingGroupsMembers creates and removes the members of the groups in SCIM provider
// returns the lists of groups members created in the SCIM provider
// with the ids of these groups members.
func reconcilingGroupsMembers(
	ctx context.Context,
	scim SCIMService, create,
	remove *model.GroupsMembersResult,
) (*model.GroupsMembersResult, error) {
	if scim == nil {
		return nil, ErrSCIMServiceNil
	}
	if create == nil {
		return nil, ErrCreateGroupsMembersResultNil
	}
	if remove == nil {
		return nil, ErrDeleteGroupsMembersResultNil
	}

	var (
		created *model.GroupsMembersResult
		err     error
	)

	if create.Items == 0 {
		slog.Info("no users to be joined to groups")
		created = model.GroupsMembersResultBuilder().Build()
	} else {
		slog.Warn("joining users to groups", "groups", create.Items)
		created, err = scim.CreateGroupsMembers(ctx, create)
		if err != nil {
			return nil, fmt.Errorf("error creating groups members in SCIM provider: %w", err)
		}
	}

	if remove.Items == 0 {
		slog.Info("no users to be removed from groups")
	} else {
		slog.Warn("removing users from groups", "groups", remove.Items)
		if err := scim.DeleteGroupsMembers(ctx, remove); err != nil {
			return nil, fmt.Errorf("error removing groups members from SCIM provider: %w", err)
		}
	}

	return created, nil
}
