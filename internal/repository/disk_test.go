package repository

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/slashdevops/idp-scim-sync/internal/model"
	"github.com/stretchr/testify/assert"
)

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
		assert.NoError(t, err)
		assert.NotNil(t, state)
		assert.Equal(t, "", state.LastSync)
		assert.Equal(t, "", state.HashCode)

		assert.Equal(t, 0, state.Resources.Groups.Items)
		assert.Equal(t, "", state.Resources.Groups.HashCode)
		assert.Equal(t, 0, len(state.Resources.Groups.Resources))

		assert.Equal(t, 0, state.Resources.Users.Items)
		assert.Equal(t, "", state.Resources.Users.HashCode)
		assert.Equal(t, 0, len(state.Resources.Users.Resources))

		assert.Equal(t, 0, state.Resources.GroupsUsers.Items)
		assert.Equal(t, "", state.Resources.GroupsUsers.HashCode)
		assert.Equal(t, 0, len(state.Resources.GroupsUsers.Resources))
	})

	t.Run("Golden files", func(t *testing.T) {
		stateFile, err := os.OpenFile("testdata/"+stateFileName, os.O_RDWR, 0644)
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

		assert.Equal(t, 1, state.Resources.GroupsUsers.Items)
		assert.Equal(t, "123456789", state.Resources.GroupsUsers.HashCode)
		assert.Equal(t, 1, len(state.Resources.GroupsUsers.Resources))
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
			Resources: model.StateResources{
				Groups: model.GroupsResult{
					Items:    1,
					HashCode: "1234567890",
					Resources: []model.Group{
						{
							IPID:     "1",
							Name:     "group 1",
							Email:    "group.1@mail.com",
							HashCode: "123456789",
						},
					},
				},
				Users: model.UsersResult{
					Items:    1,
					HashCode: "hashCode",
					Resources: []model.User{
						{
							IPID:        "1",
							Name:        model.Name{FamilyName: "1", GivenName: "user"},
							DisplayName: "user 1",
							Email:       "user.1@mail.com",
							HashCode:    "123456789",
						},
					},
				},
				GroupsUsers: model.GroupsUsersResult{
					Items:    1,
					HashCode: "123456789",
					Resources: []model.GroupUsers{
						{
							Items:    1,
							HashCode: "123456789",
							Group: model.Group{
								IPID:     "1",
								Name:     "group 1",
								Email:    "group.1@mail.com",
								HashCode: "123456789",
							},
							Resources: []model.User{
								{
									IPID:        "1",
									Name:        model.Name{FamilyName: "1", GivenName: "user"},
									DisplayName: "user 1",
									Email:       "user.1@mail.com",
									HashCode:    "123456789",
								},
							},
						},
					},
				},
			},
		}

		err = repo.SetState(context.TODO(), stateDef)
		assert.NoError(t, err)

		stateFileRO, err := os.OpenFile(stateFile.Name(), os.O_RDONLY, 0644)
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

		assert.Equal(t, 1, state.Resources.GroupsUsers.Items)
		assert.Equal(t, "123456789", state.Resources.GroupsUsers.HashCode)
		assert.Equal(t, 1, len(state.Resources.GroupsUsers.Resources))

		os.Remove(tmpDir)
		stateFile.Close()
	})
}
