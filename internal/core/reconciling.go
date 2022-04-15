package core

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/slashdevops/idp-scim-sync/internal/model"
)

var (
	// ErrCreateGroupsResultNil is returned when the create *model.GroupsResult argument is nil
	ErrCreateGroupsResultNil = fmt.Errorf("create Groups Result is nil")

	// ErrUpdateGroupsResultNil is returned when the update *model.GroupsResult argument is nil
	ErrUpdateGroupsResultNil = fmt.Errorf("update Groups Result is nil")

	// ErrDeleteGroupsResultNil is returned when the delete *model.GroupsResult argument is nil
	ErrDeleteGroupsResultNil = fmt.Errorf("remove Groups Result is nil")

	// ErrCreateUsersResultNil is returned when the create *model.UsersResult argument is nil
	ErrCreateUsersResultNil = fmt.Errorf("create Users Result is nil")

	// ErrUpdateUsersResultNil is returned when the update *model.UsersResult argument is nil
	ErrUpdateUsersResultNil = fmt.Errorf("update Users Result is nil")

	// ErrDeleteUsersResultNil is returned when the delete *model.UsersResult argument is nil
	ErrDeleteUsersResultNil = fmt.Errorf("remove Users Result is nil")

	// ErrCreateGroupsMembersResultNil is returned when the SCIM *model.GroupsMembersResult argument is nil
	ErrCreateGroupsMembersResultNil = fmt.Errorf("create Groups Members Result is nil")

	// ErrDeleteGroupsMembersResultNil is returned when the SCIM *model.GroupsMembersResult argument is nil
	ErrDeleteGroupsMembersResultNil = fmt.Errorf("remove Groups Members Result is nil")
)

// reconcilingGroups receives lists of groups to create, update, equals and remove in the SCIM service
// returns the lists of groups created and updated in the SCIM service with the Ids of these groups.
func reconcilingGroups(ctx context.Context, scim SCIMService, create, update, remove *model.GroupsResult) (created, updated *model.GroupsResult, e error) {
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

	var err error

	if create.Items == 0 {
		log.Info("no groups to be create")
		created = &model.GroupsResult{Items: 0, Resources: []*model.Group{}}
	} else {
		log.WithField("quantity", create.Items).Warn("creating groups")
		created, err = scim.CreateGroups(ctx, create)
		if err != nil {
			return nil, nil, fmt.Errorf("error creating groups in SCIM Provider: %w", err)
		}
	}

	if update.Items == 0 {
		log.Info("no groups to be updated")
		updated = &model.GroupsResult{Items: 0, Resources: []*model.Group{}}
	} else {
		log.WithField("quantity", update.Items).Warn("updating groups")
		updated, err = scim.UpdateGroups(ctx, update)
		if err != nil {
			return nil, nil, fmt.Errorf("error updating groups in SCIM Provider: %w", err)
		}
	}

	if remove.Items == 0 {
		log.Info("no groups to be removed")
	} else {
		log.WithField("quantity", remove.Items).Warn("deleting groups")
		if err := scim.DeleteGroups(ctx, remove); err != nil {
			return nil, nil, fmt.Errorf("error deleting groups in SCIM Provider: %w", err)
		}
	}

	return
}

// reconcilingUsers creates, updates and removes users in the SCIM service
// returns the lists of users created and updated in the SCIM service with the Ids of these users in the SCIM service
func reconcilingUsers(ctx context.Context, scim SCIMService, create, update, remove *model.UsersResult) (created, updated *model.UsersResult, e error) {
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

	var err error

	if create.Items == 0 {
		log.Info("no users to be created")
		created = &model.UsersResult{Items: 0, Resources: []*model.User{}}
	} else {
		log.WithField("quantity", create.Items).Warn("creating users")
		created, err = scim.CreateUsers(ctx, create)
		if err != nil {
			return nil, nil, fmt.Errorf("error creating users in SCIM Provider: %w", err)
		}
	}

	if update.Items == 0 {
		log.Info("no users to be updated")
		updated = &model.UsersResult{Items: 0, Resources: []*model.User{}}
	} else {
		log.WithField("quantity", update.Items).Warn("updating users")
		updated, err = scim.UpdateUsers(ctx, update)
		if err != nil {
			return nil, nil, fmt.Errorf("error updating users in SCIM Provider: %w", err)
		}
	}

	if remove.Items == 0 {
		log.Info("no users to be removed")
	} else {
		log.WithField("quantity", remove.Items).Warn("deleting users")
		if err := scim.DeleteUsers(ctx, remove); err != nil {
			return nil, nil, fmt.Errorf("error deleting users in SCIM Provider: %w", err)
		}
	}

	return
}

// reconcilingGroupsMembers creates and removes the members of the groups in the SCIM service
// returns the lists of groups members created in the SCIM service with the Ids of these groups members in the SCIM service
func reconcilingGroupsMembers(ctx context.Context, scim SCIMService, create, remove *model.GroupsMembersResult) (created *model.GroupsMembersResult, e error) {
	if scim == nil {
		return nil, ErrSCIMServiceNil
	}
	if create == nil {
		return nil, ErrCreateGroupsMembersResultNil
	}
	if remove == nil {
		return nil, ErrDeleteGroupsMembersResultNil
	}

	var err error

	if create.Items == 0 {
		log.Info("no users to be joined to groups")
		created = &model.GroupsMembersResult{Items: 0, Resources: []*model.GroupMembers{}}
	} else {
		log.WithField("quantity", create.Items).Warn("joining users to groups")
		created, err = scim.CreateGroupsMembers(ctx, create)
		if err != nil {
			return nil, fmt.Errorf("error creating groups members in SCIM Provider: %w", err)
		}
	}

	if remove.Items == 0 {
		log.Info("no users to be removed from groups")
	} else {
		log.WithField("quantity", remove.Items).Warn("removing users from groups")
		if err := scim.DeleteGroupsMembers(ctx, remove); err != nil {
			return nil, fmt.Errorf("error removing users from groups in SCIM Provider: %w", err)
		}
	}

	return
}
