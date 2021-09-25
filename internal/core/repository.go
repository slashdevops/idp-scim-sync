package core

import (
	"github.com/slashdevops/idp-scim-sync/internal/model"
)

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -package=mocks -destination=../mocks/core/repository_mocks.go -source=repository.go

type StateRepository interface {
	GetState() (*model.State, error)
	UpdateState(state *model.State) error
}
