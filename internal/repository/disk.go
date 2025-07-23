package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/slashdevops/idp-scim-sync/internal/model"
)

// consume io.ReadWriter

// DiskRepository represents a disk based state repository and implement core.StateRepository interface
type DiskRepository struct {
	stateFile io.ReadWriter
}

// NewDiskRepository creates a new disk based state repository
func NewDiskRepository(stateFile io.ReadWriter) (*DiskRepository, error) {
	if stateFile == nil {
		return nil, &ErrStateFileNil{Message: "state file cannot be nil"}
	}

	return &DiskRepository{
		stateFile: stateFile,
	}, nil
}

// GetState returns the state from the state file
func (dr *DiskRepository) GetState(_ context.Context) (*model.State, error) {
	var err error

	data, err := io.ReadAll(dr.stateFile)
	if err != nil {
		return nil, &ErrReadingStateFile{Message: fmt.Sprintf("error reading state file: %s", err)}
	}

	// if the state file is empty, create a new empty state
	// necessary to avoid error when unmarshalling empty state with pointers
	if len(data) == 0 {
		return nil, &ErrStateFileEmpty{Message: "state file is empty"}
	}

	var state model.State
	err = json.Unmarshal(data, &state)
	if err != nil {
		return nil, fmt.Errorf("disk: error unmarshalling state: %w", err)
	}
	return &state, nil
}

// SetState sets the state in the state file
func (dr *DiskRepository) SetState(_ context.Context, state *model.State) error {
	enc := json.NewEncoder(dr.stateFile)
	enc.SetIndent("", "  ")
	err := enc.Encode(state)
	if err != nil {
		return fmt.Errorf("disk: error encoding state: %w", err)
	}

	return nil
}

// ErrStateFileEmpty indicates that the state file is empty.
type ErrStateFileEmpty struct {
	Message string
}

func (e *ErrStateFileEmpty) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode(), e.ErrorMessage())
}

func (e *ErrStateFileEmpty) ErrorMessage() string {
	return e.Message
}
func (e *ErrStateFileEmpty) ErrorCode() string { return "ErrStateFileEmpty" }

// ErrReadingStateFile indicates an error occurred while reading the state file.
type ErrReadingStateFile struct {
	Message string
}

func (e *ErrReadingStateFile) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode(), e.ErrorMessage())
}

func (e *ErrReadingStateFile) ErrorMessage() string {
	return e.Message
}
func (e *ErrReadingStateFile) ErrorCode() string { return "ErrReadingStateFile" }

// ErrStateFileNil indicates that the state file is nil.
type ErrStateFileNil struct {
	Message string
}

func (e *ErrStateFileNil) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode(), e.ErrorMessage())
}

func (e *ErrStateFileNil) ErrorMessage() string {
	return e.Message
}
func (e *ErrStateFileNil) ErrorCode() string { return "ErrStateFileNil" }
