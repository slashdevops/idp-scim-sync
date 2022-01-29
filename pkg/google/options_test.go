package google

import (
	"reflect"
	"testing"
)

func TestWithIncludeDerivedMembership(t *testing.T) {
	t.Run("validate the return type", func(t *testing.T) {
		var ggmo GetGroupMembersOption

		got := WithIncludeDerivedMembership(false)

		if reflect.TypeOf(got) != reflect.TypeOf(ggmo) {
			t.Errorf("WithIncludeDerivedMembership() return %T, different type than %T", got, ggmo)
		}
	})

	t.Run("validate the return values", func(t *testing.T) {
		opt := WithIncludeDerivedMembership(true)
		got := getGroupMembersOptions{}
		opt(&got)

		want := getGroupMembersOptions{
			includeDerivedMembership: true,
		}

		if !reflect.DeepEqual(got.includeDerivedMembership, want.includeDerivedMembership) {
			t.Errorf("got = %v, want %v", got.includeDerivedMembership, want.includeDerivedMembership)
		}
	})
}

func TestWithMaxResults(t *testing.T) {
	t.Run("validate the return type", func(t *testing.T) {
		var ggmo GetGroupMembersOption
		got := WithMaxResults(100)

		if reflect.TypeOf(got) != reflect.TypeOf(ggmo) {
			t.Errorf("WithMaxResults() return %T, different type than %T", got, ggmo)
		}
	})

	t.Run("validate the return values", func(t *testing.T) {
		opt := WithMaxResults(2)
		got := getGroupMembersOptions{}
		opt(&got)

		want := getGroupMembersOptions{
			maxResults: 2,
		}

		if !reflect.DeepEqual(got.maxResults, want.maxResults) {
			t.Errorf("got = %v, want %v", got.maxResults, want.maxResults)
		}
	})
}

func TestWithPageToken(t *testing.T) {
	t.Run("validate the return type", func(t *testing.T) {
		var ggmo GetGroupMembersOption
		got := WithPageToken("thisIsAToken")

		if reflect.TypeOf(got) != reflect.TypeOf(ggmo) {
			t.Errorf("WithPageToken() return %T, different type than %T", got, ggmo)
		}
	})

	t.Run("validate the return values", func(t *testing.T) {
		opt := WithPageToken("thisIsAToken")
		got := getGroupMembersOptions{}
		opt(&got)

		want := getGroupMembersOptions{
			pageToken: "thisIsAToken",
		}

		if !reflect.DeepEqual(got.pageToken, want.pageToken) {
			t.Errorf("got = %s, want %s", got.pageToken, want.pageToken)
		}
	})
}

func TestWithRoles(t *testing.T) {
	t.Run("validate the return type", func(t *testing.T) {
		var ggmo GetGroupMembersOption
		got := WithRoles("OWNER,MANAGER,MEMBER")

		if reflect.TypeOf(got) != reflect.TypeOf(ggmo) {
			t.Errorf("WithRoles() return %T, different type than %T", got, ggmo)
		}
	})

	t.Run("validate the return values", func(t *testing.T) {
		opt := WithRoles("OWNER,MANAGER,MEMBER")
		got := getGroupMembersOptions{}
		opt(&got)

		want := getGroupMembersOptions{
			roles: "OWNER,MANAGER,MEMBER",
		}

		if !reflect.DeepEqual(got.roles, want.roles) {
			t.Errorf("got = %v, want %v", got.roles, want.roles)
		}
	})
}
