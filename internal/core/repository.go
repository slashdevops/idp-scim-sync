package core

import (
	"context"

	"github.com/slashdevops/idp-scim-sync/internal/model"
)

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -package=mocks -destination=../mocks/core/repository_mocks.go -source=repository.go

type StateRepository interface {
	GetState(ctx context.Context) (*model.State, error)
	UpdateState(ctx context.Context, state *model.State) error
}
