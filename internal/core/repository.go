package core

import "github.com/slashdevops/idp-scim-sync/internal/model"

type SyncRepository interface {
	StoreGroups(gr *model.GroupsResult) (model.StoreGroupsResult, error)
	StoreUsers(ur *model.UsersResult) (model.StoreUsersResult, error)
	StoreGroupsMembers(gr *model.GroupsMembersResult) (model.StoreGroupsMembersResult, error)
	StoreState(state *model.SyncState) (model.StoreStateResult, error)
	GetState() (model.SyncState, error)
	GetGroups(place string) (*model.GroupsResult, error)
	GetUsers(place string) (*model.UsersResult, error)
	GetGroupsMembers(place string) (*model.GroupsMembersResult, error)
}
