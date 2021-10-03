package disk

import (
	"encoding/json"
	"io"
	"sync"

	"github.com/pkg/errors"

	"github.com/slashdevops/idp-scim-sync/internal/model"
	"github.com/slashdevops/idp-scim-sync/internal/repository"
)

// implement core.SateRepository

// consume io.ReadWriter

var (
	ErrDataFilesNil               = errors.New("data files are nil")
	ErrStateFileNil               = errors.New("state file is nil")
	ErrGroupsFileNil              = errors.New("groups file is nil")
	ErrUsersFileNil               = errors.New("users file is nil")
	ErrGroupsUsersFileNil         = errors.New("groups users file is nil")
	ErrInconsistentHashCode       = errors.New("inconsistent hash code")
	ErrInconsistentNumberElements = errors.New("inconsistent number of elements")
)

type DBFiles struct {
	state           io.ReadWriter
	groups          io.ReadWriter
	groupsMeta      io.ReadWriter
	users           io.ReadWriter
	usersMeta       io.ReadWriter
	groupsUsers     io.ReadWriter
	groupsUsersMeta io.ReadWriter
}

type DiskRepository struct {
	mu *sync.RWMutex
	db *DBFiles
}

func NewDiskRepository(db *DBFiles) (*DiskRepository, error) {
	if db == nil {
		return nil, errors.Wrapf(ErrDataFilesNil, "NewDiskRepository")
	}

	return &DiskRepository{
		mu: &sync.RWMutex{},
		db: db,
	}, nil
}

func (dr *DiskRepository) GetGroups() (*model.GroupsResult, error) {
	dr.mu.RLock()
	defer dr.mu.RUnlock()

	grps := make([]model.Group, 0)
	d := json.NewDecoder(dr.db.groups)
	for {
		var grp model.Group
		if err := d.Decode(&grp); err == io.EOF {
			break
		} else if err != nil {
			return nil, errors.Wrapf(err, "disk.GetGroups")
		}
		grps = append(grps, grp)
	}

	grpsMeta, err := dr.GetGroupsMeta()
	if err != nil {
		return nil, errors.Wrapf(err, "disk.GetGroups.GetGroupsMeta")
	}

	// check consistency
	// if len(grps) != grpsMeta.Items {
	// 	return nil, errors.Wrapf(ErrInconsistentNumberElements, "disk.GetGroups")
	// }

	// if hash.Get(grps) != grpsMeta.HashCode {
	// 	return nil, errors.Wrapf(ErrInconsistentHashCode, "disk.GetGroups")
	// }

	gr := &model.GroupsResult{
		Items:     grpsMeta.Items,
		HashCode:  grpsMeta.HashCode,
		Resources: grps,
	}

	return gr, nil
}

func (dr *DiskRepository) GetGroupsMeta() (*repository.GroupsMetaIndex, error) {
	dr.mu.RLock()
	defer dr.mu.RUnlock()

	var grpsMeta repository.GroupsMetaIndex
	if err := json.NewDecoder(dr.db.groupsMeta).Decode(&grpsMeta); err == io.EOF {
	} else if err != nil {
		return nil, errors.Wrapf(err, "disk.GetGroupsMeta")
	}

	return &grpsMeta, nil
}

func (dr *DiskRepository) GetUsers() (*model.UsersResult, error) {
	dr.mu.RLock()
	defer dr.mu.RUnlock()

	usrs := make([]model.User, 0)
	d := json.NewDecoder(dr.db.users)
	for {
		var usr model.User
		if err := d.Decode(&usr); err == io.EOF {
			break
		} else if err != nil {
			return nil, errors.Wrapf(err, "disk.GetUsers")
		}
		usrs = append(usrs, usr)
	}

	usrsMeta, err := dr.GetUsersMeta()
	if err != nil {
		return nil, errors.Wrapf(err, "disk.GetUsers.GetUsersMeta")
	}

	// check consistency
	// if len(usrs) != usrsMeta.Items {
	// 	return nil, errors.Wrapf(ErrInconsistentNumberElements, "disk.GetUsers")
	// }

	// if hash.Get(usrs) != usrsMeta.HashCode {
	// 	return nil, errors.Wrapf(ErrInconsistentHashCode, "disk.GetUsers")
	// }

	us := &model.UsersResult{
		Items:     usrsMeta.Items,
		HashCode:  usrsMeta.HashCode,
		Resources: usrs,
	}

	return us, nil
}

func (dr *DiskRepository) GetUsersMeta() (*repository.UsersMetaIndex, error) {
	dr.mu.RLock()
	defer dr.mu.RUnlock()

	var usrsMeta repository.UsersMetaIndex
	if err := json.NewDecoder(dr.db.usersMeta).Decode(&usrsMeta); err == io.EOF {
	} else if err != nil {
		return nil, errors.Wrapf(err, "disk.GetUsersMeta")
	}

	return &usrsMeta, nil
}

func (dr *DiskRepository) GetGroupsUsers() (*model.GroupsUsersResult, error) {
	dr.mu.RLock()
	defer dr.mu.RUnlock()

	grpsUsrs := make([]model.GroupUsers, 0)
	d := json.NewDecoder(dr.db.groupsUsers)
	for {
		var grpUsrs model.GroupUsers
		if err := d.Decode(&grpUsrs); err == io.EOF {
			break
		} else if err != nil {
			return nil, errors.Wrapf(err, "disk.GetGroupsUsers")
		}
		grpsUsrs = append(grpsUsrs, grpUsrs)
	}

	grpsUsrsMeta, err := dr.GetGroupsUsersMeta()
	if err != nil {
		return nil, errors.Wrapf(err, "disk.GetGroupsUsers.GetGroupsUsersMeta")
	}

	// check consistency
	// if len(grpsUsrs) != grpsUsrsMeta.Items {
	// 	return nil, errors.Wrapf(ErrInconsistentNumberElements, "disk.GetUsers")
	// }

	// if hash.Get(grpsUsrs) != grpsUsrsMeta.HashCode {
	// 	return nil, errors.Wrapf(ErrInconsistentHashCode, "disk.GetUsers")
	// }

	us := &model.GroupsUsersResult{
		Items:     grpsUsrsMeta.Items,
		HashCode:  grpsUsrsMeta.HashCode,
		Resources: grpsUsrs,
	}

	return us, nil
}

func (dr *DiskRepository) GetGroupsUsersMeta() (*repository.GroupsUsersMetaIndex, error) {
	dr.mu.RLock()
	defer dr.mu.RUnlock()

	var grpsUsrsMeta repository.GroupsUsersMetaIndex
	if err := json.NewDecoder(dr.db.groupsUsersMeta).Decode(&grpsUsrsMeta); err == io.EOF {
	} else if err != nil {
		return nil, errors.Wrapf(err, "disk.GetGroupsUsersMeta")
	}

	return &grpsUsrsMeta, nil
}

func (dr *DiskRepository) GetState() (*model.State, error) {
	dr.mu.RLock()
	defer dr.mu.RUnlock()
	var err error

	grps, err := dr.GetGroups()
	if err != nil {
		return nil, errors.Wrapf(err, "disk.GetState.GetGroups")
	}

	usrs, err := dr.GetUsers()
	if err != nil {
		return nil, errors.Wrapf(err, "disk.GetState.GetUsers")
	}

	grpsUsrs, err := dr.GetGroupsUsers()
	if err != nil {
		return nil, errors.Wrapf(err, "disk.GetState.GetGroupsUsers")
	}

	// read only the metadata not the Resoruces
	var loadedState repository.StateMetaIndex
	if err = json.NewDecoder(dr.db.state).Decode(&loadedState); err == io.EOF {
	} else if err != nil {
		return nil, errors.Wrapf(err, "GetState")
	}

	r := model.StateResources{
		Groups:      *grps,
		Users:       *usrs,
		GroupsUsers: *grpsUsrs,
	}
	state := &model.State{
		LastSync:  loadedState.LastSync,
		HashCode:  loadedState.HashCode,
		Resources: r,
	}

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
