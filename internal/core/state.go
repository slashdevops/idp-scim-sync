package core

import "github.com/slashdevops/idp-scim-sync/internal/state"

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -package=mocks -destination=../mocks/core/state_mocks.go -source=state.go

type SyncState interface {
	GetName() string
	Empty() bool
	Build(groups *state.StoreGroupsResult, groupsUsers *state.StoreGroupsUsersResult, users *state.StoreUsersResult) (*state.State, error)
}
