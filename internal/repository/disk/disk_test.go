package disk

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/slashdevops/idp-scim-sync/internal/repository"
	"github.com/stretchr/testify/assert"
)

func Test_StateRepository_NewDiskRepository(t *testing.T) {
	t.Run("Call with nil parameters", func(t *testing.T) {
		repo, err := NewDiskRepository(nil)
		assert.Error(t, err)
		assert.Nil(t, repo)
	})
}

func Test_StateRepository_GetGroups(t *testing.T) {
	t.Run("Empty file", func(t *testing.T) {
		tmpDir := os.TempDir()
		defer os.Remove(tmpDir)

		groupsFile, err := ioutil.TempFile(tmpDir, repository.GroupsIndexName)
		if err != nil {
			t.Fatal(err)
		}
		groupsMetaFile, err := ioutil.TempFile(tmpDir, repository.GroupsMetaName)
		if err != nil {
			t.Fatal(err)
		}

		db := &DBFiles{
			groups:     groupsFile,
			groupsMeta: groupsMetaFile,
		}

		repo, err := NewDiskRepository(db)
		assert.NoError(t, err)
		assert.NotNil(t, repo)

		grps, err := repo.GetGroups()
		assert.NoError(t, err)
		assert.NotNil(t, grps)
	})

	t.Run("Golden file", func(t *testing.T) {
		groupsFile, err := os.OpenFile("testdata/"+repository.GroupsIndexName, os.O_RDWR, 0644)
		if err != nil {
			t.Fatal(err)
		}
		defer groupsFile.Close()

		groupsMetaFile, err := os.OpenFile("testdata/"+repository.GroupsMetaName, os.O_RDWR, 0644)
		if err != nil {
			t.Fatal(err)
		}
		defer groupsMetaFile.Close()

		db := &DBFiles{
			groups:     groupsFile,
			groupsMeta: groupsMetaFile,
		}

		repo, err := NewDiskRepository(db)
		assert.NoError(t, err)
		assert.NotNil(t, repo)

		grps, err := repo.GetGroups()
		assert.NoError(t, err)
		assert.NotNil(t, grps)
		assert.Equal(t, 2, grps.Items)
		assert.Equal(t, "123456789", grps.HashCode)

		assert.Equal(t, 2, len(grps.Resources))

		assert.Equal(t, "1", grps.Resources[0].ID)
		assert.Equal(t, "group 1", grps.Resources[0].Name)
		assert.Equal(t, "group.1@mail.com", grps.Resources[0].Email)
		assert.Equal(t, "123456789", grps.Resources[0].HashCode)

		assert.Equal(t, "2", grps.Resources[1].ID)
		assert.Equal(t, "group 2", grps.Resources[1].Name)
		assert.Equal(t, "group.2@mail.com", grps.Resources[1].Email)
		assert.Equal(t, "987654321", grps.Resources[1].HashCode)
	})
}

func Test_StateRepository_GetGroupsMeta(t *testing.T) {
	t.Run("Empty file", func(t *testing.T) {
		tmpDir := os.TempDir()
		defer os.Remove(tmpDir)

		groupsMetaFile, err := ioutil.TempFile(tmpDir, repository.GroupsMetaName)
		if err != nil {
			t.Fatal(err)
		}

		db := &DBFiles{
			groupsMeta: groupsMetaFile,
		}

		repo, err := NewDiskRepository(db)
		assert.NoError(t, err)
		assert.NotNil(t, repo)

		grpsMeta, err := repo.GetGroupsMeta()
		assert.NoError(t, err)
		assert.NotNil(t, grpsMeta)
		assert.Equal(t, 0, grpsMeta.Items)
		assert.Equal(t, "", grpsMeta.HashCode)
	})

	t.Run("Golden file", func(t *testing.T) {
		groupsMetaFile, err := os.OpenFile("testdata/"+repository.GroupsMetaName, os.O_RDWR, 0644)
		if err != nil {
			t.Fatal(err)
		}
		defer groupsMetaFile.Close()

		db := &DBFiles{
			groupsMeta: groupsMetaFile,
		}

		repo, err := NewDiskRepository(db)
		assert.NoError(t, err)
		assert.NotNil(t, repo)

		grpsMeta, err := repo.GetGroupsMeta()
		assert.NoError(t, err)
		assert.NotNil(t, grpsMeta)
		assert.Equal(t, 2, grpsMeta.Items)
		assert.Equal(t, "123456789", grpsMeta.HashCode)
	})
}

func Test_StateRepository_GetUsers(t *testing.T) {
	t.Run("Empty file", func(t *testing.T) {
		tmpDir := os.TempDir()
		defer os.Remove(tmpDir)

		usersFile, err := ioutil.TempFile(tmpDir, repository.UsersIndexName)
		if err != nil {
			t.Fatal(err)
		}
		usersMetaFile, err := ioutil.TempFile(tmpDir, repository.UsersMetaName)
		if err != nil {
			t.Fatal(err)
		}

		db := &DBFiles{
			users:     usersFile,
			usersMeta: usersMetaFile,
		}

		repo, err := NewDiskRepository(db)
		assert.NoError(t, err)
		assert.NotNil(t, repo)

		usrs, err := repo.GetUsers()
		assert.NoError(t, err)
		assert.NotNil(t, usrs)
	})

	t.Run("Golden file", func(t *testing.T) {
		usersFile, err := os.OpenFile("testdata/"+repository.UsersIndexName, os.O_RDWR, 0644)
		if err != nil {
			t.Fatal(err)
		}
		defer usersFile.Close()

		usersMetaFile, err := os.OpenFile("testdata/"+repository.UsersMetaName, os.O_RDWR, 0644)
		if err != nil {
			t.Fatal(err)
		}
		defer usersMetaFile.Close()

		db := &DBFiles{
			users:     usersFile,
			usersMeta: usersMetaFile,
		}

		repo, err := NewDiskRepository(db)
		assert.NoError(t, err)
		assert.NotNil(t, repo)

		usrs, err := repo.GetUsers()
		assert.NoError(t, err)
		assert.NotNil(t, usrs)
		assert.Equal(t, 2, usrs.Items)
		assert.Equal(t, "123456789", usrs.HashCode)

		assert.Equal(t, 2, len(usrs.Resources))

		assert.Equal(t, "1", usrs.Resources[0].ID)
		assert.Equal(t, "1", usrs.Resources[0].Name.FamilyName)
		assert.Equal(t, "user", usrs.Resources[0].Name.GivenName)
		assert.Equal(t, "user.1@mail.com", usrs.Resources[0].Email)
		assert.Equal(t, "123456789", usrs.Resources[0].HashCode)

		assert.Equal(t, "2", usrs.Resources[1].ID)
		assert.Equal(t, "2", usrs.Resources[1].Name.FamilyName)
		assert.Equal(t, "user", usrs.Resources[1].Name.GivenName)
		assert.Equal(t, "user.2@mail.com", usrs.Resources[1].Email)
		assert.Equal(t, "987654321", usrs.Resources[1].HashCode)
	})
}

func Test_StateRepository_GetUsersMeta(t *testing.T) {
	t.Run("Empty file", func(t *testing.T) {
		tmpDir := os.TempDir()
		defer os.Remove(tmpDir)

		usersMetaFile, err := ioutil.TempFile(tmpDir, repository.UsersMetaName)
		if err != nil {
			t.Fatal(err)
		}

		db := &DBFiles{
			usersMeta: usersMetaFile,
		}

		repo, err := NewDiskRepository(db)
		assert.NoError(t, err)
		assert.NotNil(t, repo)

		usrsMeta, err := repo.GetUsersMeta()
		assert.NoError(t, err)
		assert.NotNil(t, usrsMeta)
		assert.Equal(t, 0, usrsMeta.Items)
		assert.Equal(t, "", usrsMeta.HashCode)
	})

	t.Run("Golden file", func(t *testing.T) {
		usersMetaFile, err := os.OpenFile("testdata/"+repository.UsersMetaName, os.O_RDWR, 0644)
		if err != nil {
			t.Fatal(err)
		}
		defer usersMetaFile.Close()

		db := &DBFiles{
			usersMeta: usersMetaFile,
		}

		repo, err := NewDiskRepository(db)
		assert.NoError(t, err)
		assert.NotNil(t, repo)

		usersMeta, err := repo.GetUsersMeta()
		assert.NoError(t, err)
		assert.NotNil(t, usersMeta)
		assert.Equal(t, 2, usersMeta.Items)
		assert.Equal(t, "123456789", usersMeta.HashCode)
	})
}

func Test_StateRepository_GetGroupsUsers(t *testing.T) {
	t.Run("Empty file", func(t *testing.T) {
		tmpDir := os.TempDir()
		defer os.Remove(tmpDir)

		groupsUsersFile, err := ioutil.TempFile(tmpDir, repository.GroupsUsersIndexName)
		if err != nil {
			t.Fatal(err)
		}
		groupsUsersMetaFile, err := ioutil.TempFile(tmpDir, repository.GroupsUsersMetaName)
		if err != nil {
			t.Fatal(err)
		}

		db := &DBFiles{
			groupsUsers:     groupsUsersFile,
			groupsUsersMeta: groupsUsersMetaFile,
		}

		repo, err := NewDiskRepository(db)
		assert.NoError(t, err)
		assert.NotNil(t, repo)

		grpsUsrs, err := repo.GetGroupsUsers()
		assert.NoError(t, err)
		assert.NotNil(t, grpsUsrs)
	})

	t.Run("Golden file", func(t *testing.T) {
		groupsUsersFile, err := os.OpenFile("testdata/"+repository.GroupsUsersIndexName, os.O_RDWR, 0644)
		if err != nil {
			t.Fatal(err)
		}
		defer groupsUsersFile.Close()

		groupsUsersMetaFile, err := os.OpenFile("testdata/"+repository.GroupsUsersMetaName, os.O_RDWR, 0644)
		if err != nil {
			t.Fatal(err)
		}
		defer groupsUsersMetaFile.Close()

		db := &DBFiles{
			groupsUsers:     groupsUsersFile,
			groupsUsersMeta: groupsUsersMetaFile,
		}

		repo, err := NewDiskRepository(db)
		assert.NoError(t, err)
		assert.NotNil(t, repo)

		grpsUsrs, err := repo.GetGroupsUsers()
		assert.NoError(t, err)
		assert.NotNil(t, grpsUsrs)
		assert.Equal(t, 2, grpsUsrs.Items)
		assert.Equal(t, "123456789", grpsUsrs.HashCode)

		assert.Equal(t, 2, len(grpsUsrs.Resources))

		assert.Equal(t, 1, grpsUsrs.Resources[0].Items)
		assert.Equal(t, "123456789", grpsUsrs.Resources[0].HashCode)
		assert.Equal(t, "1", grpsUsrs.Resources[0].Group.ID)
		assert.Equal(t, "group 1", grpsUsrs.Resources[0].Group.Name)
		assert.Equal(t, "group.1@mail.com", grpsUsrs.Resources[0].Group.Email)
		assert.Equal(t, "123456789", grpsUsrs.Resources[0].Group.HashCode)

		assert.Equal(t, "1", grpsUsrs.Resources[0].Resources[0].ID)
		assert.Equal(t, "1", grpsUsrs.Resources[0].Resources[0].Name.FamilyName)
		assert.Equal(t, "user", grpsUsrs.Resources[0].Resources[0].Name.GivenName)
		assert.Equal(t, "user 1", grpsUsrs.Resources[0].Resources[0].DisplayName)
		assert.Equal(t, true, grpsUsrs.Resources[0].Resources[0].Active)
		assert.Equal(t, "user.1@mail.com", grpsUsrs.Resources[0].Resources[0].Email)
		assert.Equal(t, "123456789", grpsUsrs.Resources[0].Resources[0].HashCode)

		assert.Equal(t, "2", grpsUsrs.Resources[1].Resources[0].ID)
		assert.Equal(t, "2", grpsUsrs.Resources[1].Resources[0].Name.FamilyName)
		assert.Equal(t, "user", grpsUsrs.Resources[1].Resources[0].Name.GivenName)
		assert.Equal(t, "user 2", grpsUsrs.Resources[1].Resources[0].DisplayName)
		assert.Equal(t, true, grpsUsrs.Resources[1].Resources[0].Active)
		assert.Equal(t, "user.2@mail.com", grpsUsrs.Resources[1].Resources[0].Email)
		assert.Equal(t, "987654321", grpsUsrs.Resources[1].Resources[0].HashCode)
	})
}

func Test_StateRepository_GetGroupsUsersMeta(t *testing.T) {
	t.Run("Empty file", func(t *testing.T) {
		tmpDir := os.TempDir()
		defer os.Remove(tmpDir)

		grpsUsersMetaFile, err := ioutil.TempFile(tmpDir, repository.GroupsUsersMetaName)
		if err != nil {
			t.Fatal(err)
		}

		db := &DBFiles{
			groupsUsersMeta: grpsUsersMetaFile,
		}

		repo, err := NewDiskRepository(db)
		assert.NoError(t, err)
		assert.NotNil(t, repo)

		grpsUsrsMeta, err := repo.GetGroupsUsersMeta()
		assert.NoError(t, err)
		assert.NotNil(t, grpsUsrsMeta)
		assert.Equal(t, 0, grpsUsrsMeta.Items)
		assert.Equal(t, "", grpsUsrsMeta.HashCode)
	})

	t.Run("Golden file", func(t *testing.T) {
		grpsUsersMetaFile, err := os.OpenFile("testdata/"+repository.GroupsUsersMetaName, os.O_RDWR, 0644)
		if err != nil {
			t.Fatal(err)
		}
		defer grpsUsersMetaFile.Close()

		db := &DBFiles{
			groupsUsersMeta: grpsUsersMetaFile,
		}

		repo, err := NewDiskRepository(db)
		assert.NoError(t, err)
		assert.NotNil(t, repo)

		grpsUsrsMeta, err := repo.GetGroupsUsersMeta()
		assert.NoError(t, err)
		assert.NotNil(t, grpsUsrsMeta)
		assert.Equal(t, 2, grpsUsrsMeta.Items)
		assert.Equal(t, "123456789", grpsUsrsMeta.HashCode)
	})
}

func Test_StateRepository_GetState(t *testing.T) {
	t.Run("Empty files", func(t *testing.T) {
		tmpDir := os.TempDir()
		defer os.Remove(tmpDir)

		stateFile, err := ioutil.TempFile(tmpDir, repository.StateIndexName)
		if err != nil {
			t.Fatal(err)
		}

		groupsFile, err := ioutil.TempFile(tmpDir, repository.GroupsIndexName)
		if err != nil {
			t.Fatal(err)
		}

		groupsMetaFile, err := ioutil.TempFile(tmpDir, repository.GroupsMetaName)
		if err != nil {
			t.Fatal(err)
		}

		usersFile, err := ioutil.TempFile(tmpDir, repository.UsersIndexName)
		if err != nil {
			t.Fatal(err)
		}

		usersMetaFile, err := ioutil.TempFile(tmpDir, repository.UsersMetaName)
		if err != nil {
			t.Fatal(err)
		}

		groupsUsersFile, err := ioutil.TempFile(tmpDir, repository.GroupsUsersIndexName)
		if err != nil {
			t.Fatal(err)
		}

		groupsUsersMetaFile, err := ioutil.TempFile(tmpDir, repository.GroupsUsersMetaName)
		if err != nil {
			t.Fatal(err)
		}

		db := &DBFiles{
			state:           stateFile,
			groups:          groupsFile,
			groupsMeta:      groupsMetaFile,
			users:           usersFile,
			usersMeta:       usersMetaFile,
			groupsUsers:     groupsUsersFile,
			groupsUsersMeta: groupsUsersMetaFile,
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
		stateFile, err := os.OpenFile("testdata/"+repository.StateIndexName, os.O_RDWR, 0644)
		if err != nil {
			t.Fatal(err)
		}
		defer stateFile.Close()

		groupsFile, err := os.OpenFile("testdata/"+repository.GroupsIndexName, os.O_RDWR, 0644)
		if err != nil {
			t.Fatal(err)
		}
		defer groupsFile.Close()

		groupsMetaFile, err := os.OpenFile("testdata/"+repository.GroupsMetaName, os.O_RDWR, 0644)
		if err != nil {
			t.Fatal(err)
		}
		defer groupsMetaFile.Close()

		usersFile, err := os.OpenFile("testdata/"+repository.UsersIndexName, os.O_RDWR, 0644)
		if err != nil {
			t.Fatal(err)
		}
		defer usersFile.Close()

		usersMetaFile, err := os.OpenFile("testdata/"+repository.UsersMetaName, os.O_RDWR, 0644)
		if err != nil {
			t.Fatal(err)
		}
		defer usersMetaFile.Close()

		groupsUsersFile, err := os.OpenFile("testdata/"+repository.GroupsUsersIndexName, os.O_RDWR, 0644)
		if err != nil {
			t.Fatal(err)
		}
		defer groupsUsersFile.Close()

		groupsUsersMetaFile, err := os.OpenFile("testdata/"+repository.GroupsUsersMetaName, os.O_RDWR, 0644)
		if err != nil {
			t.Fatal(err)
		}
		defer groupsUsersMetaFile.Close()

		db := &DBFiles{
			state:           stateFile,
			groups:          groupsFile,
			groupsMeta:      groupsMetaFile,
			users:           usersFile,
			usersMeta:       usersMetaFile,
			groupsUsers:     groupsUsersFile,
			groupsUsersMeta: groupsUsersMetaFile,
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
