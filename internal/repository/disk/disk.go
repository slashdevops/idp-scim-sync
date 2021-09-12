package disk

import (
	"errors"
	"sync"

	"github.com/slashdevops/idp-scim-sync/internal/model"
	"github.com/spf13/afero"
)

// implement core.Repository

// consume afero.Fs

type DiskRepository struct {
	mu *sync.RWMutex
	fs afero.Fs
}

func NewDiskRepository(fs afero.Fs) *DiskRepository {
	return &DiskRepository{
		mu: &sync.RWMutex{},
		fs: fs,
	}
}

func (dr *DiskRepository) GetState() (*model.State, error) {
	return nil, errors.New("not implemented")
}

func (dr *DiskRepository) GetGroups() (*model.GroupsResult, error) {
	return nil, errors.New("not implemented")
}

func (dr *DiskRepository) GetUsers() (*model.UsersResult, error) {
	return nil, errors.New("not implemented")
}

func (dr *DiskRepository) GetGroupsUsers() (*model.GroupsUsersResult, error) {
	return nil, errors.New("not implemented")
}

func (dr *DiskRepository) UpdateState(state *model.State) error {
	return errors.New("not implemented")
}
