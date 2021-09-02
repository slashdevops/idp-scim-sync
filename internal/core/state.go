package core

import "github.com/slashdevops/idp-scim-sync/internal/state"

//go:generate mockgen -package=mocks -destination=../mocks/core/state_mocks.go -source=state.go

type SyncState interface {
	GetName() string
	Empty() bool
	Build(groups *state.StoreGroupsResult, groupsUsers *state.StoreGroupsUsersResult, users *state.StoreUsersResult) (*state.State, error)
}
