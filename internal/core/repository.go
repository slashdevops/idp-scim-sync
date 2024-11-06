package core

import (
	"context"

	"github.com/slashdevops/idp-scim-sync/internal/model"
)

//go:generate go run go.uber.org/mock/mockgen@v0.5.0 -package=mocks -destination=../../mocks/core/repository_mocks.go -source=repository.go

// StateRepository is an interface for a repository that stores the state of the
// synchronization process.
// This interface needs to be implemented by the repository service.
type StateRepository interface {
	// GetState returns the state of the synchronization process.
	GetState(ctx context.Context) (*model.State, error)

	// SetState sets the state of the synchronization process.
	SetState(ctx context.Context, state *model.State) error
}
