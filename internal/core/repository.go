package core

import "github.com/slashdevops/idp-scim-sync/internal/model"

type SyncRepository interface {
	StoreGroups(gr *model.GroupsResult) (model.StoreGroupsResult, error)
	StoreUsers(ur *model.UsersResult) (model.StoreUsersResult, error)
	StoreGroupsUsers(gr *model.GroupsUsersResult) (model.StoreGroupsUsersResult, error)

	StoreState(state *model.State) (model.StoreStateResult, error)
	GetState(name string) (model.State, error)
	GetGroups(place string) (*model.GroupsResult, error)
	GetUsers(place string) (*model.UsersResult, error)
	GetGroupsUsers(place string) (*model.GroupsUsersResult, error)
}
