package disk

import (
	"encoding/json"
	"io"
	"sync"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/slashdevops/idp-scim-sync/internal/model"
)

// implement core.SateRepository

// consume io.ReadWriter

var (
	ErrDataFilesNil       = errors.New("data files are nil")
	ErrStateFileNil       = errors.New("state file is nil")
	ErrGroupsFileNil      = errors.New("groups file is nil")
	ErrUsersFileNil       = errors.New("users file is nil")
	ErrGroupsUsersFileNil = errors.New("groups users file is nil")
)

type DBFiles struct {
	state       io.ReadWriter
	groups      io.ReadWriter
	users       io.ReadWriter
	groupsUsers io.ReadWriter
}

type DiskRepository struct {
	mu *sync.RWMutex
	db *DBFiles
}

func NewDiskRepository(db *DBFiles) (*DiskRepository, error) {
	if db == nil {
		return nil, errors.Wrapf(ErrDataFilesNil, "NewDiskRepository")
	}

	if db.state == nil {
		return nil, errors.Wrapf(ErrStateFileNil, "NewDiskRepository")
	}

	if db.groups == nil {
		return nil, errors.Wrapf(ErrGroupsFileNil, "NewDiskRepository")
	}

	if db.users == nil {
		return nil, errors.Wrapf(ErrUsersFileNil, "NewDiskRepository")
	}

	if db.groupsUsers == nil {
		return nil, errors.Wrapf(ErrGroupsUsersFileNil, "NewDiskRepository")
	}

	return &DiskRepository{
		mu: &sync.RWMutex{},
		db: db,
	}, nil
}

func (dr *DiskRepository) GetState() (*model.State, error) {
	dr.mu.RLock()
	defer dr.mu.RUnlock()
	var err error

	var gr model.GroupsResult
	if err = json.NewDecoder(dr.db.groups).Decode(&gr); err == io.EOF {
	} else if err != nil {
		return nil, errors.Wrapf(err, "GetState")
	}

	var ur model.UsersResult
	if err = json.NewDecoder(dr.db.users).Decode(&ur); err == io.EOF {
	} else if err != nil {
		return nil, errors.Wrapf(err, "GetState")
	}

	var grUr model.GroupsUsersResult
	if err = json.NewDecoder(dr.db.groupsUsers).Decode(&grUr); err == io.EOF {
	} else if err != nil {
		return nil, errors.Wrapf(err, "GetState")
	}

	// read only the metadata not the Resoruces
	var loadedState model.State
	if err = json.NewDecoder(dr.db.state).Decode(&loadedState); err == io.EOF {
	} else if err != nil {
		return nil, errors.Wrapf(err, "GetState")
	}

	r := model.StateResources{
		Groups:      gr,
		Users:       ur,
		GroupsUsers: grUr,
	}
	state := &model.State{
		LastSync:  loadedState.LastSync,
		HashCode:  loadedState.HashCode,
		Resources: r,
	}

	log.Printf("loadedState: %v", loadedState)
	log.Printf("state: %v", state)

	return state, nil
}

func (dr *DiskRepository) UpdateState(state *model.State) error {
	dr.mu.Lock()
	defer dr.mu.Unlock()
	var err error

	err = json.NewEncoder(dr.db.groups).Encode(&state.Resources.Groups)
	if err != nil {
		return errors.Wrapf(err, "UpdateState")
	}

	err = json.NewEncoder(dr.db.users).Encode(&state.Resources.Users)
	if err != nil {
		return errors.Wrapf(err, "UpdateState")
	}

	err = json.NewEncoder(dr.db.groupsUsers).Encode(&state.Resources.GroupsUsers)
	if err != nil {
		return errors.Wrapf(err, "UpdateState")
	}

	newState := &model.State{
		LastSync:  state.LastSync,
		Resources: model.StateResources{},
	}

	err = json.NewEncoder(dr.db.state).Encode(&newState)
	if err != nil {
		return errors.Wrapf(err, "UpdateState")
	}

	return nil
}
