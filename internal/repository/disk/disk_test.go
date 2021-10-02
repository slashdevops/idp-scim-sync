package disk

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_StateRepository_NewDiskRepository(t *testing.T) {
	t.Run("call with nil parameters", func(t *testing.T) {
		db := &DBFiles{
			state:       nil,
			groups:      nil,
			users:       nil,
			groupsUsers: nil,
		}

		repo, err := NewDiskRepository(db)
		assert.Error(t, err)
		assert.Nil(t, repo)
	})
}

func Test_StateRepository_GetState(t *testing.T) {
	t.Run("Empty files", func(t *testing.T) {
		var content []byte
		tmpDir := os.TempDir()
		defer os.Remove(tmpDir)

		stateFile, err := ioutil.TempFile(tmpDir, "state.json")
		if err != nil {
			t.Fatal(err)
		}

		groupsFile, err := ioutil.TempFile(tmpDir, "groups.json")
		if err != nil {
			t.Fatal(err)
		}
		usersFile, err := ioutil.TempFile(tmpDir, "users.json")
		if err != nil {
			t.Fatal(err)
		}

		groupsUsersFile, err := ioutil.TempFile(tmpDir, "groups_users.json")
		if err != nil {
			t.Fatal(err)
		}

		content = []byte(`{}`)
		if _, err = stateFile.Write(content); err != nil {
			t.Fatal("Failed to write to state file", err)
		}

		content = []byte("{}")
		if _, err = groupsFile.Write(content); err != nil {
			t.Fatal("Failed to write to groups file", err)
		}

		content = []byte("{}")
		if _, err = usersFile.Write(content); err != nil {
			t.Fatal("Failed to write to users file", err)
		}

		content = []byte("{}")
		if _, err = groupsUsersFile.Write(content); err != nil {
			t.Fatal("Failed to write to groups_users file", err)
		}

		db := &DBFiles{
			state:       stateFile,
			groups:      groupsFile,
			users:       usersFile,
			groupsUsers: groupsUsersFile,
		}

		repo, err := NewDiskRepository(db)
		assert.NoError(t, err)
		assert.NotNil(t, repo)

		state, err := repo.GetState()
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
		stateFile, err := os.OpenFile("testdata/state.golden", os.O_RDWR, 0644)
		if err != nil {
			t.Fatal(err)
		}
		defer stateFile.Close()

		groupsFile, err := os.OpenFile("testdata/groups.golden", os.O_RDWR, 0644)
		if err != nil {
			t.Fatal(err)
		}
		defer groupsFile.Close()

		usersFile, err := os.OpenFile("testdata/users.golden", os.O_RDWR, 0644)
		if err != nil {
			t.Fatal(err)
		}
		defer usersFile.Close()

		groupsUsersFile, err := os.OpenFile("testdata/groups_users.golden", os.O_RDWR, 0644)
		if err != nil {
			t.Fatal(err)
		}
		defer groupsUsersFile.Close()

		db := &DBFiles{
			state:       stateFile,
			groups:      groupsFile,
			users:       usersFile,
			groupsUsers: groupsUsersFile,
		}

		repo, err := NewDiskRepository(db)
		assert.NoError(t, err)

		state, err := repo.GetState()
		assert.NoError(t, err)
		assert.NotNil(t, state)

		assert.Equal(t, "2021-09-25T20:49:46+02:00", state.LastSync)
		assert.Equal(t, "123456789", state.HashCode)

		assert.Equal(t, 2, state.Resources.Groups.Items)
		assert.Equal(t, "123456789", state.Resources.Groups.HashCode)
		assert.Equal(t, 2, len(state.Resources.Groups.Resources))

		assert.Equal(t, 2, state.Resources.Users.Items)
		assert.Equal(t, "123456789", state.Resources.Users.HashCode)
		assert.Equal(t, 2, len(state.Resources.Users.Resources))

		assert.Equal(t, 2, state.Resources.GroupsUsers.Items)
		assert.Equal(t, "123456789", state.Resources.GroupsUsers.HashCode)
		assert.Equal(t, 2, len(state.Resources.GroupsUsers.Resources))
	})
}
