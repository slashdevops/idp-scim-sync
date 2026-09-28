package model

import "testing"

// The *Result.SetHashCode methods sort a copy of Resources before hashing. The
// comparators dereference each element, so a nil entry panics. Nil entries do
// reach these slices: internal/idp.buildUser and internal/scim.buildUser both
// return nil to reject a malformed record, and their callers store the result.
//
// SetHashCode must therefore tolerate nil elements rather than crashing the
// sync mid-run.
func TestUsersResult_SetHashCode_nilResource(t *testing.T) {
	ur := &UsersResult{
		Items: 3,
		Resources: []*User{
			{UserName: "b@example.com", HashCode: "bbb"},
			nil,
			{UserName: "a@example.com", HashCode: "aaa"},
		},
	}

	ur.SetHashCode()

	if ur.HashCode == "" {
		t.Error("SetHashCode() left HashCode empty")
	}
}

func TestGroupsResult_SetHashCode_nilResource(t *testing.T) {
	gr := &GroupsResult{
		Items: 3,
		Resources: []*Group{
			{Name: "beta", HashCode: "bbb"},
			nil,
			{Name: "alpha", HashCode: "aaa"},
		},
	}

	gr.SetHashCode()

	if gr.HashCode == "" {
		t.Error("SetHashCode() left HashCode empty")
	}
}

func TestMembersResult_SetHashCode_nilResource(t *testing.T) {
	mr := &MembersResult{
		Items: 3,
		Resources: []*Member{
			{IPID: "2", Email: "b@example.com"},
			nil,
			{IPID: "1", Email: "a@example.com"},
		},
	}

	mr.SetHashCode()

	if mr.HashCode == "" {
		t.Error("SetHashCode() left HashCode empty")
	}
}

func TestGroupsMembersResult_SetHashCode_nilResource(t *testing.T) {
	gmr := &GroupsMembersResult{
		Items: 3,
		Resources: []*GroupMembers{
			{Group: &Group{Name: "beta"}, HashCode: "bbb"},
			nil,
			{Group: &Group{Name: "alpha"}, HashCode: "aaa"},
		},
	}

	gmr.SetHashCode()

	if gmr.HashCode == "" {
		t.Error("SetHashCode() left HashCode empty")
	}
}

// GroupMembers.SetHashCode sorts its members by Email, so a nil member panics
// there too.
func TestGroupMembers_SetHashCode_nilResource(t *testing.T) {
	gm := &GroupMembers{
		Group: &Group{Name: "alpha"},
		Items: 3,
		Resources: []*Member{
			{Email: "b@example.com"},
			nil,
			{Email: "a@example.com"},
		},
	}

	gm.SetHashCode()

	if gm.HashCode == "" {
		t.Error("SetHashCode() left HashCode empty")
	}
}
