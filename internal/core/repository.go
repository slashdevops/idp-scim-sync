package core

import (
	"github.com/slashdevops/idp-scim-sync/internal/model"
)

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -package=mocks -destination=../mocks/core/repository_mocks.go -source=repository.go

type Repository interface {
	GetState() (*model.State, error)
	GetGroups() (*model.GroupsResult, error)
	GetUsers() (*model.UsersResult, error)
	GetGroupsUsers() (*model.GroupsUsersResult, error)
	UpdateState(state *model.State) error
}
