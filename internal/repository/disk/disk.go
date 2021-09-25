package disk

import (
	"encoding/json"
	"io"
	"sync"

	"github.com/slashdevops/idp-scim-sync/internal/model"
)

// implement core.SateRepository

// consume afero.Fs

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

func NewDiskRepository(db *DBFiles) *DiskRepository {
	// paths := &repository.DB{
	// 	StateFile:       "state.json",
	// 	GroupsFile:      "groups.json",
	// 	UsersFile:       "users.json",
	// 	GroupsUsersFile: "groups_users.json",
	// }

	// stateFile, err := os.Open(paths.StateFile)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// groupsFile, err := os.Open(paths.GroupsFile)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// usersFile, err := os.Open(paths.UsersFile)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// groupsUsersFile, err := os.Open(paths.GroupsUsersFile)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	return &DiskRepository{
		mu: &sync.RWMutex{},
		db: db,
	}
}

func (dr *DiskRepository) GetState() (*model.State, error) {
	dr.mu.RLock()
	defer dr.mu.RUnlock()
	var err error

	var gr *model.GroupsResult
	groupsDecoder := json.NewDecoder(dr.db.groups)
	err = groupsDecoder.Decode(&gr)
	if err != nil {
		return nil, err
	}

	var ur *model.UsersResult
	usersDecoder := json.NewDecoder(dr.db.users)
	err = usersDecoder.Decode(&ur)
	if err != nil {
		return nil, err
	}

	var gur *model.GroupsUsersResult
	groupsUsersDecoder := json.NewDecoder(dr.db.groupsUsers)
	err = groupsUsersDecoder.Decode(&gur)
	if err != nil {
		return nil, err
	}

	// read only the metadata not the Resoruces
	var loadedState *model.State
	decoder := json.NewDecoder(dr.db.state)
	err = decoder.Decode(&loadedState)
	if err != nil {
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
