package disk

import (
	"sync"

	"github.com/slashdevops/idp-scim-sync/internal/core"
)

type diskRepository struct {
	mu   *sync.RWMutex
	path string
}

func NewDiskRepository(path string) core.SyncRepository {
	return &diskRepository{}
}

func (dr *diskRepository) StoreGroups(gr *core.GroupsResult) (core.StoreGroupsResult, error) {
	return core.StoreGroupsResult{}, nil
}

func (dr *diskRepository) StoreUsers(ur *core.UsersResult) (core.StoreUsersResult, error) {
	return core.StoreUsersResult{}, nil
}

func (dr *diskRepository) StoreGroupsMembers(gr *core.GroupsMembersResult) (core.StoreGroupsMembersResult, error) {
	return core.StoreGroupsMembersResult{}, nil
}

func (dr *diskRepository) StoreState(state *core.SyncState) (core.StoreStateResult, error) {
	return core.StoreStateResult{}, nil
}

func (dr *diskRepository) GetState() (core.SyncState, error) {
	return core.SyncState{}, nil
}

func (dr *diskRepository) GetGroups(place string) (*core.GroupsResult, error) {
	return &core.GroupsResult{}, nil
}

func (dr *diskRepository) GetUsers(place string) (*core.UsersResult, error) {
	return &core.UsersResult{}, nil
}

func (dr *diskRepository) GetGroupsMembers(place string) (*core.GroupsMembersResult, error) {
	return &core.GroupsMembersResult{}, nil
}
