package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/pkg/errors"

	"github.com/slashdevops/idp-scim-sync/internal/model"
)

// consume io.ReadWriter

const (
	stateFileName = "state.json"
)

// ErrStateFileNil is returned when the state file is nil
var ErrStateFileNil = errors.New("disk: state file is nil")

// DiskRepository represents a disk based state repository and implement core.StateRepository interface
type DiskRepository struct {
	stateFile io.ReadWriter
}

// NewDiskRepository creates a new disk based state repository
func NewDiskRepository(stateFile io.ReadWriter) (*DiskRepository, error) {
	if stateFile == nil {
		return nil, ErrStateFileNil
	}

	return &DiskRepository{
		stateFile: stateFile,
	}, nil
}

// GetState returns the state from the state file
func (dr *DiskRepository) GetState(ctx context.Context) (*model.State, error) {
	var err error

	var state model.State
	if err = json.NewDecoder(dr.stateFile).Decode(&state); err == io.EOF {
	} else if err != nil {
		return nil, fmt.Errorf("disk: error decoding state: %w", err)
	}

	return &state, nil
}

// SetState sets the state in the state file
func (dr *DiskRepository) SetState(ctx context.Context, state *model.State) error {
	enc := json.NewEncoder(dr.stateFile)
	enc.SetIndent("", "  ")
	err := enc.Encode(state)
	if err != nil {
		return fmt.Errorf("disk: error encoding state: %w", err)
	}

	return nil
}
