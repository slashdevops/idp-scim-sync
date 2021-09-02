package core

import (
	"github.com/slashdevops/idp-scim-sync/internal/model"
	"github.com/slashdevops/idp-scim-sync/internal/state"
)

//go:generate mockgen -package=mocks -destination=../mocks/core/repository_mocks.go -source=repository.go

type SyncRepository interface {
	StoreGroups(gr *model.GroupsResult) (state.StoreGroupsResult, error)
	StoreUsers(ur *model.UsersResult) (state.StoreUsersResult, error)
	StoreGroupsUsers(gr *model.GroupsUsersResult) (state.StoreGroupsUsersResult, error)

	StoreState(sgr *state.StoreGroupsResult, sur *state.StoreUsersResult, sgur *state.StoreGroupsUsersResult) (state.StoreStateResult, error)
	GetState() (state.State, error)
	GetGroups() (*model.GroupsResult, error)
	GetUsers() (*model.UsersResult, error)
	GetGroupsUsers() (*model.GroupsUsersResult, error)
}
