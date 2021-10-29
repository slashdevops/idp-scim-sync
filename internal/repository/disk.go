package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/pkg/errors"

	"github.com/slashdevops/idp-scim-sync/internal/model"
)

// implement core.SateRepository

// consume io.ReadWriter

const (
	stateFileName = "state.json"
)

var ErrStateFileNil = errors.New("disk: state file is nil")

type DiskRepository struct {
	mu        *sync.RWMutex
	stateFile io.ReadWriter
}

func NewDiskRepository(stateFile io.ReadWriter) (*DiskRepository, error) {
	if stateFile == nil {
		return nil, ErrStateFileNil
	}

	return &DiskRepository{
		mu:        &sync.RWMutex{},
		stateFile: stateFile,
	}, nil
}

func (dr *DiskRepository) GetState(ctx context.Context) (*model.State, error) {
	dr.mu.RLock()
	defer dr.mu.RUnlock()
	var err error

	var state model.State
	if err = json.NewDecoder(dr.stateFile).Decode(&state); err == io.EOF {
	} else if err != nil {
		return nil, fmt.Errorf("disk: error decoding state: %w", err)
	}

	return &state, nil
}

func (dr *DiskRepository) SetState(ctx context.Context, state *model.State) error {
	dr.mu.Lock()
	defer dr.mu.Unlock()

	enc := json.NewEncoder(dr.stateFile)
	enc.SetIndent("", "  ")
	err := enc.Encode(state)
	if err != nil {
		return fmt.Errorf("disk: error encoding state: %w", err)
	}

	return nil
}
