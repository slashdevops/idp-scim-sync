package disk

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"sync"

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
		return nil, ErrDataFilesNil
	}

	if db.state == nil {
		return nil, ErrStateFileNil
	}

	if db.groups == nil {
		return nil, ErrGroupsFileNil
	}

	if db.users == nil {
		return nil, ErrUsersFileNil
	}

	if db.groupsUsers == nil {
		return nil, ErrGroupsUsersFileNil
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

	var gr *model.GroupsResult
	groupsDecoder := json.NewDecoder(dr.db.groups)
	if err = groupsDecoder.Decode(&gr); err == io.EOF {
		log.Println("no data to decode")
		gr = &model.GroupsResult{
			Items:     0,
			HashCode:  "",
			Resources: make([]model.Group, 0),
		}
	} else if err != nil {
		return nil, err
	}

	var ur *model.UsersResult
	usersDecoder := json.NewDecoder(dr.db.users)
	if err = usersDecoder.Decode(&ur); err == io.EOF {
		log.Println("no data to decode")
		ur = &model.UsersResult{
			Items:     0,
			HashCode:  "",
			Resources: make([]model.User, 0),
		}
	} else if err != nil {
		return nil, err
	}

	var gur *model.GroupsUsersResult
	groupsUsersDecoder := json.NewDecoder(dr.db.groupsUsers)
	if err = groupsUsersDecoder.Decode(&gur); err == io.EOF {
		log.Println("no data to decode")
		gur = &model.GroupsUsersResult{
			Items:     0,
			HashCode:  "",
			Resources: make([]model.GroupUsers, 0),
		}
	} else if err != nil {
		return nil, err
	}

	// read only the metadata not the Resoruces
	var loadedState *model.State
	decoder := json.NewDecoder(dr.db.state)
	if err = decoder.Decode(loadedState); err == io.EOF {
		log.Println("no state data to decode")
		loadedState = &model.State{}
	} else if err != nil {
		return nil, err
	}

	state := &model.State{
		LastSync: loadedState.LastSync,
		Resources: model.StateResources{
			Groups:      gr,
			Users:       ur,
			GroupsUsers: gur,
		},
	}

	log.Printf("state: %v", state)

	return state, nil
}

func (dr *DiskRepository) UpdateState(state *model.State) error {
	dr.mu.Lock()
	defer dr.mu.Unlock()
	var err error

	groupsEncoder := json.NewEncoder(dr.db.groups)
	err = groupsEncoder.Encode(&state.Resources.Groups)
	if err != nil {
		return err
	}

	usersEncoder := json.NewEncoder(dr.db.users)
	err = usersEncoder.Encode(&state.Resources.Users)
	if err != nil {
		return err
	}

	groupsUsersEncoder := json.NewEncoder(dr.db.groupsUsers)
	err = groupsUsersEncoder.Encode(&state.Resources.GroupsUsers)
	if err != nil {
		return err
	}

	newState := &model.State{
		LastSync:  state.LastSync,
		Resources: model.StateResources{},
	}

	decoder := json.NewEncoder(dr.db.state)
	err = decoder.Encode(&newState)
	if err != nil {
		return err
	}

	return nil
}
