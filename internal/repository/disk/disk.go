package disk

import (
	"sync"

	"github.com/slashdevops/idp-scim-sync/internal/model"
)

type DiskRepository struct {
	mu   *sync.RWMutex
	path string
}

func NewDiskRepository(path string) *DiskRepository {
	return &DiskRepository{}
}

func (dr *DiskRepository) StoreGroups(gr *model.GroupsResult) (model.StoreGroupsResult, error) {
	return model.StoreGroupsResult{}, nil
}

func (dr *DiskRepository) StoreUsers(ur *model.UsersResult) (model.StoreUsersResult, error) {
	return model.StoreUsersResult{}, nil
}

func (dr *DiskRepository) StoreGroupsMembers(gr *model.GroupsMembersResult) (model.StoreGroupsMembersResult, error) {
	return model.StoreGroupsMembersResult{}, nil
}

func (dr *DiskRepository) StoreState(state *model.SyncState) (model.StoreStateResult, error) {
	return model.StoreStateResult{}, nil
}

func (dr *DiskRepository) GetState() (model.SyncState, error) {
	return model.SyncState{}, nil
}

func (dr *DiskRepository) GetGroups(place string) (*model.GroupsResult, error) {
	return &model.GroupsResult{}, nil
}

func (dr *DiskRepository) GetUsers(place string) (*model.UsersResult, error) {
	return &model.UsersResult{}, nil
}

func (dr *DiskRepository) GetGroupsMembers(place string) (*model.GroupsMembersResult, error) {
	return &model.GroupsMembersResult{}, nil
}
