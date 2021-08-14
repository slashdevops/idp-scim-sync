package disk

import (
	"sync"

	s "github.com/slashdevops/idp-scim-sync/internal/sync"
)

type diskRepository struct {
	mu   *sync.RWMutex
	path string
}

func NewDiskRepository(path string) s.SyncRepository {
	return &diskRepository{}
}

func (dr *diskRepository) StoreGroups(gr *s.GroupsResult) (s.StoreGroupsResult, error) {

	return s.StoreGroupsResult{}, nil
}

func (dr *diskRepository) StoreUsers(ur *s.UsersResult) (s.StoreUsersResult, error) {
	return s.StoreUsersResult{}, nil
}

func (dr *diskRepository) StoreGroupsMembers(gr *s.GroupsMembersResult) (s.StoreGroupsMembersResult, error) {
	return s.StoreGroupsMembersResult{}, nil
}

func (dr *diskRepository) StoreState(state *s.SyncState) (s.StoreStateResult, error) {
	return s.StoreStateResult{}, nil
}
