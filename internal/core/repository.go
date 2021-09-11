package core

import (
	"github.com/slashdevops/idp-scim-sync/internal/model"
)

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -package=mocks -destination=../mocks/core/repository_mocks.go -source=repository.go

type StateRepository interface {
	StoreGroups(gr *model.GroupsResult) (model.StoreGroupsResult, error)
	StoreUsers(ur *model.UsersResult) (model.StoreUsersResult, error)
	StoreGroupsUsers(gr *model.GroupsUsersResult) (model.StoreGroupsUsersResult, error)

	StoreState(sgr *model.StoreGroupsResult, sur *model.StoreUsersResult, sgur *model.StoreGroupsUsersResult) (model.StoreStateResult, error)

	GetState() (model.State, error)
	GetGroups() (*model.GroupsResult, error)
	GetUsers() (*model.UsersResult, error)
	GetGroupsUsers() (*model.GroupsUsersResult, error)
}
