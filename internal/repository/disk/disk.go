package disk

import (
	"encoding/json"
	"io"
	"sync"

	"github.com/pkg/errors"

	"github.com/slashdevops/idp-scim-sync/internal/model"
)

// implement core.SateRepository

// consume io.ReadWriter

var ErrStateFileNil = errors.New("state file is nil")

type DiskRepository struct {
	mu        *sync.RWMutex
	stateFile io.ReadWriter
}

func NewDiskRepository(stateFile io.ReadWriter) (*DiskRepository, error) {
	if stateFile == nil {
		return nil, errors.Wrapf(ErrStateFileNil, "NewDiskRepository")
	}

	return &DiskRepository{
		mu:        &sync.RWMutex{},
		stateFile: stateFile,
	}, nil
}

func (dr *DiskRepository) GetState() (*model.State, error) {
	dr.mu.RLock()
	defer dr.mu.RUnlock()
	var err error

	var state model.State
	if err = json.NewDecoder(dr.stateFile).Decode(&state); err == io.EOF {
	} else if err != nil {
		return nil, errors.Wrapf(err, "GetState")
	}

	return &state, nil
}

func (dr *DiskRepository) UpdateState(state *model.State) error {
	dr.mu.Lock()
	defer dr.mu.Unlock()

	enc := json.NewEncoder(dr.stateFile)
	enc.SetIndent("", "  ")
	err := enc.Encode(state)
	if err != nil {
		return errors.Wrapf(err, "UpdateState")
	}

	return nil
}
