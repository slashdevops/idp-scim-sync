package core

import "github.com/slashdevops/idp-scim-sync/internal/state"

type SyncState interface {
	GetName() string
	Empty() bool
	Build(groups *state.StoreGroupsResult, groupsUsers *state.StoreGroupsUsersResult, users *state.StoreUsersResult) (*state.State, error)
}
