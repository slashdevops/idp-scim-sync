package core

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/slashdevops/idp-scim-sync/internal/model"
)

var (
	ErrCreateGroupsResultNil        = fmt.Errorf("create Groups Result is nil")
	ErrUpdateGroupsResultNil        = fmt.Errorf("update Groups Result is nil")
	ErrDeleteGroupsResultNil        = fmt.Errorf("delete Groups Result is nil")
	ErrCreateUsersResultNil         = fmt.Errorf("create Users Result is nil")
	ErrUpdateUsersResultNil         = fmt.Errorf("update Users Result is nil")
	ErrDeleteUsersResultNil         = fmt.Errorf("delete Users Result is nil")
	ErrCreateGroupsMembersResultNil = fmt.Errorf("create Groups Members Result is nil")
	ErrDeleteGroupsMembersResultNil = fmt.Errorf("delete Groups Members Result is nil")
)

// reconcilingGroups receives lists of groups to create, update, equals and delete in the SCIM service
// returns the lists of groups created and updated in the SCIM service with the Ids of these groups.
func reconcilingGroups(ctx context.Context, scim SCIMService, create *model.GroupsResult, update *model.GroupsResult, delete *model.GroupsResult) (created *model.GroupsResult, updated *model.GroupsResult, e error) {
	if scim == nil {
		return nil, nil, ErrSCIMServiceNil
	}
	if create == nil {
		return nil, nil, ErrCreateGroupsResultNil
	}
	if update == nil {
		return nil, nil, ErrUpdateGroupsResultNil
	}
	if delete == nil {
		return nil, nil, ErrDeleteGroupsResultNil
	}

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

// reconcilingUsers creates, updates and deletes users in the SCIM service
// returns the lists of users created and updated in the SCIM service with the Ids of these users in the SCIM service
func reconcilingUsers(ctx context.Context, scim SCIMService, create *model.UsersResult, update *model.UsersResult, delete *model.UsersResult) (created *model.UsersResult, updated *model.UsersResult, e error) {
	if scim == nil {
		return nil, nil, ErrSCIMServiceNil
	}
	if create == nil {
		return nil, nil, ErrCreateUsersResultNil
	}
	if update == nil {
		return nil, nil, ErrUpdateUsersResultNil
	}
	if delete == nil {
		return nil, nil, ErrDeleteUsersResultNil
	}

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

// reconcilingGroupsMembers creates and deletes the members of the groups in the SCIM service
// returns the lists of groups members created in the SCIM service with the Ids of these groups members in the SCIM service
func reconcilingGroupsMembers(ctx context.Context, scim SCIMService, create *model.GroupsMembersResult, delete *model.GroupsMembersResult) (created *model.GroupsMembersResult, e error) {
	if scim == nil {
		return nil, ErrSCIMServiceNil
	}
	if create == nil {
		return nil, ErrCreateGroupsMembersResultNil
	}
	if delete == nil {
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
