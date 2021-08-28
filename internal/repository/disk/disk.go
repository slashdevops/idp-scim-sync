package disk

import (
	"sync"

	"github.com/slashdevops/idp-scim-sync/internal/model"
	"github.com/slashdevops/idp-scim-sync/internal/repository"
)

type diskRepository struct {
	mu   *sync.RWMutex
	path string
}

func NewDiskRepository(path string) repository.SyncRepository {
	return &diskRepository{}
}

func (dr *diskRepository) StoreGroups(gr *model.GroupsResult) (model.StoreGroupsResult, error) {
	return model.StoreGroupsResult{}, nil
}

func (dr *diskRepository) StoreUsers(ur *model.UsersResult) (model.StoreUsersResult, error) {
	return model.StoreUsersResult{}, nil
}

func (dr *diskRepository) StoreGroupsMembers(gr *model.GroupsMembersResult) (model.StoreGroupsMembersResult, error) {
	return model.StoreGroupsMembersResult{}, nil
}

func (dr *diskRepository) StoreState(state *model.SyncState) (model.StoreStateResult, error) {
	return model.StoreStateResult{}, nil
}

func (dr *diskRepository) GetState() (model.SyncState, error) {
	return model.SyncState{}, nil
}

func (dr *diskRepository) GetGroups(place string) (*model.GroupsResult, error) {
	return &model.GroupsResult{}, nil
}

func (dr *diskRepository) GetUsers(place string) (*model.UsersResult, error) {
	return &model.UsersResult{}, nil
}

func (dr *diskRepository) GetGroupsMembers(place string) (*model.GroupsMembersResult, error) {
	return &model.GroupsMembersResult{}, nil
}
