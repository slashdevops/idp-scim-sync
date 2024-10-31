package core

import (
	"reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"
	mocks "github.com/slashdevops/idp-scim-sync/mocks/core"
)

func TestWithIdentityProviderGroupsFilter(t *testing.T) {
	t.Run("validate the return type", func(t *testing.T) {
		var sso SyncServiceOption
		filter := []string{"group1", "group2"}

		got := WithIdentityProviderGroupsFilter(filter)

		if reflect.TypeOf(got) != reflect.TypeOf(sso) {
			t.Errorf("WithIdentityProviderGroupsFilter() return %T, different type than %T", got, sso)
		}
	})

	t.Run("validate the return values", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		filter := []string{"group1", "group2"}
		prov := mocks.NewMockIdentityProviderService(mockCtrl)
		scim := mocks.NewMockSCIMService(mockCtrl)
		repo := mocks.NewMockStateRepository(mockCtrl)

		got, _ := NewSyncService(prov, scim, repo, WithIdentityProviderGroupsFilter(filter))

		want := &SyncService{
			prov:             prov,
			provGroupsFilter: filter,
			provUsersFilter:  []string{},
			scim:             scim,
			repo:             repo,
		}

		// test length
		if len(got.provGroupsFilter) != len(want.provGroupsFilter) {
			t.Errorf("len(got.provGroupsFilter) != len(want.provGroupsFilter), got %v, want %v", len(got.provGroupsFilter), len(want.provGroupsFilter))
		}

		// test values
		for i := range got.provGroupsFilter {
			if got.provGroupsFilter[i] != want.provGroupsFilter[i] {
				t.Errorf("got.provGroupsFilter[%d] != want.provGroupsFilter[%d], got %v, want %v", i, i, got.provGroupsFilter[i], want.provGroupsFilter[i])
			}
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("NewSyncService() got = %s, want %s", got, want)
		}
	})
}

func TestWithIdentityProviderUsersFilter(t *testing.T) {
	t.Run("validate the return type", func(t *testing.T) {
		var sso SyncServiceOption
		filter := []string{"user1", "user2"}

		got := WithIdentityProviderUsersFilter(filter)

		if reflect.TypeOf(got) != reflect.TypeOf(sso) {
			t.Errorf("WithIdentityProviderUsersFilter() return %T, different type than %T", got, sso)
		}
	})

	t.Run("validate the return values", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		filter := []string{"user1", "user2"}
		prov := mocks.NewMockIdentityProviderService(mockCtrl)
		scim := mocks.NewMockSCIMService(mockCtrl)
		repo := mocks.NewMockStateRepository(mockCtrl)

		got, _ := NewSyncService(prov, scim, repo, WithIdentityProviderUsersFilter(filter))

		want := &SyncService{
			prov:             prov,
			provGroupsFilter: []string{},
			provUsersFilter:  filter,
			scim:             scim,
			repo:             repo,
		}

		// test length
		if len(got.provGroupsFilter) != len(want.provGroupsFilter) {
			t.Errorf("len(got.provGroupsFilter) != len(want.provGroupsFilter), got %v, want %v", len(got.provUsersFilter), len(want.provUsersFilter))
		}

		// test values
		for i := range got.provUsersFilter {
			if got.provUsersFilter[i] != want.provUsersFilter[i] {
				t.Errorf("got.provUsersFilter[%d] != want.provUsersFilter[%d], got %v, want %v", i, i, got.provUsersFilter[i], want.provUsersFilter[i])
			}
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("NewSyncService() got = %s, want %s", got, want)
		}
	})
}
