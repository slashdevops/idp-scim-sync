package repository

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/slashdevops/idp-scim-sync/internal/model"
	"github.com/stretchr/testify/assert"
)

const stateFileName = "state.json"

func TestStateRepository_NewDiskRepository(t *testing.T) {
	t.Run("Call with nil parameters", func(t *testing.T) {
		repo, err := NewDiskRepository(nil)
		assert.Error(t, err)
		assert.Nil(t, repo)
	})
}

func TestStateRepository_GetState(t *testing.T) {
	t.Run("Empty file", func(t *testing.T) {
		tmpDir := os.TempDir()
		defer os.Remove(tmpDir)

		stateFile, err := ioutil.TempFile(tmpDir, stateFileName)
		if err != nil {
			t.Fatal(err)
		}

		repo, err := NewDiskRepository(stateFile)
		assert.NoError(t, err)
		assert.NotNil(t, repo)

		state, err := repo.GetState(context.TODO())
		assert.Error(t, err)
		assert.Nil(t, state)
	})

	t.Run("Golden files", func(t *testing.T) {
		stateFile, err := os.OpenFile("testdata/"+stateFileName, os.O_RDWR, 0o644)
		if err != nil {
			t.Fatal(err)
		}
		defer stateFile.Close()

		repo, err := NewDiskRepository(stateFile)
		assert.NoError(t, err)

		state, err := repo.GetState(context.TODO())
		assert.NoError(t, err)
		assert.NotNil(t, state)

		assert.Equal(t, "2021-09-25T20:49:46+02:00", state.LastSync)
		assert.Equal(t, "hashCode", state.HashCode)

		assert.Equal(t, 1, state.Resources.Groups.Items)
		assert.Equal(t, "123456789", state.Resources.Groups.HashCode)
		assert.Equal(t, 1, len(state.Resources.Groups.Resources))

		assert.Equal(t, 1, state.Resources.Users.Items)
		assert.Equal(t, "hashCode", state.Resources.Users.HashCode)
		assert.Equal(t, 1, len(state.Resources.Users.Resources))
		assert.Equal(t, "1", state.Resources.Users.Resources[0].IPID)
		assert.Equal(t, "user 1", state.Resources.Users.Resources[0].DisplayName)
	})
}

func TestStateRepository_SetState(t *testing.T) {
	t.Run("Empty file", func(t *testing.T) {
		tmpDir := os.TempDir()

		stateFile, err := ioutil.TempFile(tmpDir, stateFileName)
		if err != nil {
			t.Fatal(err)
		}

		repo, err := NewDiskRepository(stateFile)
		assert.NoError(t, err)
		assert.NotNil(t, repo)

		stateDef := &model.State{
			LastSync: "2021-09-25T20:49:46+02:00",
			HashCode: "hashCode",
			Resources: &model.StateResources{
				Groups: &model.GroupsResult{
					Items:    1,
					HashCode: "1234567890",
					Resources: []*model.Group{
						{
							IPID:     "1",
							Name:     "group 1",
							Email:    "group.1@mail.com",
							HashCode: "123456789",
						},
					},
				},
				Users: &model.UsersResult{
					Items:    1,
					HashCode: "hashCode",
					Resources: []*model.User{
						{
							IPID:        "1",
							Name:        model.Name{FamilyName: "1", GivenName: "user"},
							DisplayName: "user 1",
							Emails:      []model.Email{{Value: "user.1@mail.com", Type: "work", Primary: true}},
							HashCode:    "123456789",
						},
					},
				},
			},
		}

		err = repo.SetState(context.TODO(), stateDef)
		assert.NoError(t, err)

		stateFileRO, err := os.OpenFile(stateFile.Name(), os.O_RDONLY, 0o644)
		if err != nil {
			t.Fatal(err)
		}
		defer stateFileRO.Close()

		repoRO, err := NewDiskRepository(stateFileRO)
		assert.NoError(t, err)
		assert.NotNil(t, repoRO)

		state, err := repoRO.GetState(context.TODO())
		assert.NoError(t, err)
		assert.NotNil(t, state)
		assert.Equal(t, "2021-09-25T20:49:46+02:00", state.LastSync)
		assert.Equal(t, "hashCode", state.HashCode)

		assert.Equal(t, 1, state.Resources.Groups.Items)
		assert.Equal(t, "1234567890", state.Resources.Groups.HashCode)
		assert.Equal(t, 1, len(state.Resources.Groups.Resources))

		assert.Equal(t, 1, state.Resources.Users.Items)
		assert.Equal(t, "hashCode", state.Resources.Users.HashCode)
		assert.Equal(t, 1, len(state.Resources.Users.Resources))
		assert.Equal(t, "1", state.Resources.Users.Resources[0].IPID)
		assert.Equal(t, "user 1", state.Resources.Users.Resources[0].DisplayName)

		os.Remove(tmpDir)
		stateFile.Close()
	})
}
