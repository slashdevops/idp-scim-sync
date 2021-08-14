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

func (dr *diskRepository) StoreGroups(gr *s.GroupsResult) error {

	return nil
}

func (dr *diskRepository) StoreUsers(ur *s.UsersResult) error {
	return nil
}

func (dr *diskRepository) StoreGroupsMembers(gr *s.GroupsMembersResult) error {
	return nil
}
