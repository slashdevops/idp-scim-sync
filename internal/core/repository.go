package core

import (
	"github.com/slashdevops/idp-scim-sync/internal/model"
)

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -package=mocks -destination=../mocks/core/repository_mocks.go -source=repository.go

type StateRepository interface {
	GetState() (*model.State, error)
	GetStateMetadata() (*model.StateMetadata, error)
	UpdateState(state *model.State) error

	GetGroups() (*model.GroupsResult, error)
	GetGroupsMetadata() (*model.GroupsMetadata, error)

	GetUsers() (*model.UsersResult, error)
	GetUsersMetadata() (*model.UsersMetadata, error)

	GetGroupsUsers() (*model.GroupsUsersResult, error)
	GetGroupsUsersMetadata() (*model.GroupsUsersMetadata, error)
}
