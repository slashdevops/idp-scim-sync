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

var (
	// ErrStateFileNil is returned when the state file is nil
	ErrStateFileNil = errors.New("disk: state file is nil")

	// ErrReadingStateFile is returned when the state file is not found
	ErrReadingStateFile = errors.New("disk: error reading state file")
)

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

	data, err := io.ReadAll(dr.stateFile)
	if err != nil {
		return nil, ErrReadingStateFile
	}

	// if the state file is empty, create a new empty state
	// necessary to avoid error when unmarshalling empty state with pointers
	if len(data) == 0 {
		return nil, fmt.Errorf("disk: error reading state: %w", err)
	}

	var state model.State
	err = json.Unmarshal(data, &state)
	if err != nil {
		return nil, fmt.Errorf("disk: error unmarshalling state: %w", err)
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
