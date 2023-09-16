package model

import (
	"errors"
)

var (
	// ErrIdentityProviderGroupsMembersNil is returned when the idp *GroupsMembersResult argument is nil
	ErrIdentityProviderGroupsMembersNil = errors.New("identity provider groups members is nil")

	// ErrSCIMGroupsMembersNil is returned when the scim *GroupsMembersResult argument is nil
	ErrSCIMGroupsMembersNil = errors.New("scim groups members is nil")

	// ErrIdentityProviderGroupsNil is returned when the idp *GroupsResult argument is nil
	ErrIdentityProviderGroupsNil = errors.New("identity provider groups is nil")

	// ErrSCIMGroupsNil is returned when the scim *GroupsResult argument is nil
	ErrSCIMGroupsNil = errors.New("scim groups is nil")

	// ErrIdentityProviderUsersNil is returned when the idp *UsersResult argument is nil
	ErrIdentityProviderUsersNil = errors.New("identity provider users is nil")

	// ErrSCIMUsersNil is returned when the scim *UsersResult argument is nil
	ErrSCIMUsersNil = errors.New("scim users is nil")
)

// MembersOperations returns datasets used to perform different operations over the SCIM side
// return 4 objet of GroupsMembersResult
// create: groups that exist in "idp" but not in "scim" or "state"
// update: groups that exist in "idp" and in "scim" or "state" but attributes changed in idp
// equal: groups that exist in both "idp" and "scim" or "state" and their attributes are equal
// remove: groups that exist in "scim" or "state" but not in "idp"
//
// also this extract the id from scim to fill the results
func MembersOperations(idp, scim *GroupsMembersResult) (create, equal, remove *GroupsMembersResult, err error) {
	if idp == nil {
		create, equal, remove, err = nil, nil, nil, ErrIdentityProviderGroupsMembersNil
		return
	}
	if scim == nil {
		create, equal, remove, err = nil, nil, nil, ErrSCIMGroupsMembersNil
		return
	}

	toCreate, toEqual, toRemove := membersDataSets(idp.Resources, scim.Resources)

	create = GroupsMembersResultBuilder().WithResources(toCreate).Build()
	equal = GroupsMembersResultBuilder().WithResources(toEqual).Build()
	remove = GroupsMembersResultBuilder().WithResources(toRemove).Build()

	return
}

// GroupsOperations returns the differences between the groups in the
// this use the Groups Name as the key.
// SCIM Groups cannot be updated.
// return 4 objet of GroupsResult
// create: groups that exist in "idp" but not in "scim" or "state"
// update: groups that exist in "idp" and in "scim" or "state" but attributes changed in idp
// equal: groups that exist in both "idp" and "scim" or "state" and their attributes are equal
// remove: groups that exist in "scim" or "state" but not in "idp"
//
// also this extract the id from scim to fill the results
func GroupsOperations(idp, scim *GroupsResult) (create, update, equal, remove *GroupsResult, err error) {
	if idp == nil {
		create, update, equal, remove, err = nil, nil, nil, nil, ErrIdentityProviderGroupsNil
		return
	}
	if scim == nil {
		create, update, equal, remove, err = nil, nil, nil, nil, ErrSCIMGroupsNil
		return
	}

	idpGroups := make(map[string]struct{})
	scimGroups := make(map[string]Group)

	toCreate := make([]*Group, 0)
	toUpdate := make([]*Group, 0)
	toEqual := make([]*Group, 0)
	toRemove := make([]*Group, 0)

	for _, gr := range idp.Resources {
		idpGroups[gr.Name] = struct{}{}
	}

	for _, gr := range scim.Resources {
		scimGroups[gr.Name] = *gr
	}

	// loop over idp to see what to create and what to update
	for _, group := range idp.Resources {
		if _, ok := scimGroups[group.Name]; !ok {
			toCreate = append(toCreate, group)
		} else {
			group.SCIMID = scimGroups[group.Name].SCIMID

			if group.IPID != scimGroups[group.Name].IPID {
				toUpdate = append(toUpdate, group)
			} else {
				toEqual = append(toEqual, group)
			}
		}
	}

	// loop over scim to see what to remove
	for _, group := range scim.Resources {
		if _, ok := idpGroups[group.Name]; !ok {
			toRemove = append(toRemove, group)
		}
	}

	create = GroupsResultBuilder().WithResources(toCreate).Build()
	update = GroupsResultBuilder().WithResources(toUpdate).Build()
	equal = GroupsResultBuilder().WithResources(toEqual).Build()
	remove = GroupsResultBuilder().WithResources(toRemove).Build()

	return
}

// UsersOperations returns datasets used to perform different operations over the SCIM side
// return 4 objet of UsersResult
// create: users that exist in "idp" but not in "scim" or "state"
// update: users that exist in "idp" and in "scim" or "state" but attributes changed in idp
// equal: users that exist in both "idp" and "scim" or "state" and their attributes are equal
// remove: users that exist in "scim" or "state" but not in "idp"
func UsersOperations(idp, scim *UsersResult) (create, update, equal, remove *UsersResult, err error) {
	if idp == nil {
		create, update, equal, remove, err = nil, nil, nil, nil, ErrIdentityProviderUsersNil
		return
	}
	if scim == nil {
		create, update, equal, remove, err = nil, nil, nil, nil, ErrSCIMUsersNil
		return
	}

	idpUsers := make(map[string]struct{})
	scimUsers := make(map[string]User)

	toCreate := make([]*User, 0)
	toUpdate := make([]*User, 0)
	toEqual := make([]*User, 0)
	toRemove := make([]*User, 0)

	for _, usr := range idp.Resources {
		idpUsers[usr.GetPrimaryEmailAddress()] = struct{}{}
	}

	for _, usr := range scim.Resources {
		scimUsers[usr.GetPrimaryEmailAddress()] = *usr
	}

	// new users and what equal to them
	for _, usr := range idp.Resources {
		if _, ok := scimUsers[usr.GetPrimaryEmailAddress()]; !ok {
			toCreate = append(toCreate, usr)
		} else {
			usr.SCIMID = scimUsers[usr.GetPrimaryEmailAddress()].SCIMID

			if usr.HashCode != scimUsers[usr.GetPrimaryEmailAddress()].HashCode {
				toUpdate = append(toUpdate, usr)
			} else {
				toEqual = append(toEqual, usr)
			}
		}
	}

	for _, usr := range scim.Resources {
		if _, ok := idpUsers[usr.GetPrimaryEmailAddress()]; !ok {
			toRemove = append(toRemove, usr)
		}
	}

	create = UsersResultBuilder().WithResources(toCreate).Build()
	update = UsersResultBuilder().WithResources(toUpdate).Build()
	equal = UsersResultBuilder().WithResources(toEqual).Build()
	remove = UsersResultBuilder().WithResources(toRemove).Build()

	return
}

// MergeGroupsResult merges n GroupsResult result
// NOTE: this function does not check the content of the GroupsResult, so
// the return could have duplicated groups
func MergeGroupsResult(grs ...*GroupsResult) (merged *GroupsResult) {
	groups := make([]*Group, 0)

	for _, gr := range grs {
		groups = append(groups, gr.Resources...)
	}

	merged = GroupsResultBuilder().WithResources(groups).Build()

	return
}

// MergeUsersResult merges n UsersResult result
// NOTE: this function does not check the content of the UsersResult, so
// the return could have duplicated users
func MergeUsersResult(urs ...*UsersResult) (merged *UsersResult) {
	users := make([]*User, 0)

	for _, u := range urs {
		users = append(users, u.Resources...)
	}

	merged = UsersResultBuilder().WithResources(users).Build()

	return
}

// MergeGroupsMembersResult merges n GroupMembers result
// NOTE: this function does not check the content of the GroupMembers, so
// the return could have duplicated groupsMembers
func MergeGroupsMembersResult(gms ...*GroupsMembersResult) (merged *GroupsMembersResult) {
	groupsMembers := make([]*GroupMembers, 0)

	for _, gm := range gms {
		groupsMembers = append(groupsMembers, gm.Resources...)
	}

	merged = GroupsMembersResultBuilder().WithResources(groupsMembers).Build()

	return
}

// UpdateGroupsMembersSCIMID updates the SCIMID of the group in the idp object
// this is necessary because during the sync process we can create users and groups and to add
// these users to the groups we need to have the SCIMID of the user and the group
func UpdateGroupsMembersSCIMID(idp *GroupsMembersResult, scimGroups *GroupsResult, scimUsers *UsersResult) *GroupsMembersResult {
	groups := make(map[string]Group)
	users := make(map[string]User)

	for _, group := range scimGroups.Resources {
		groups[group.Name] = *group
	}

	for _, user := range scimUsers.Resources {
		users[user.GetPrimaryEmailAddress()] = *user
	}

	gms := make([]*GroupMembers, 0)
	for _, groupMembers := range idp.Resources {
		mbs := make([]*Member, 0)

		g := GroupBuilder().
			WithIPID(groupMembers.Group.IPID).
			WithSCIMID(groups[groupMembers.Group.Name].SCIMID).
			WithName(groupMembers.Group.Name).
			WithEmail(groupMembers.Group.Email).
			Build()

		for _, member := range groupMembers.Resources {
			m := MemberBuilder().
				WithIPID(member.IPID).
				WithSCIMID(users[member.Email].SCIMID).
				WithEmail(member.Email).
				WithStatus(member.Status).
				Build()

			mbs = append(mbs, m)
		}

		gm := GroupMembersBuilder().
			WithGroup(g).
			WithResources(mbs).
			Build()

		gms = append(gms, gm)
	}

	gmr := GroupsMembersResultBuilder().WithResources(gms).Build()

	return gmr
}

// membersDataSets returns the data sets of the members of the groups
// given an idp and a scim groups members this function
// this function performs the comparison between the idp and the scim data
// and returns the data sets of the members that need to be created, equal and removed
func membersDataSets(idp, scim []*GroupMembers) (create, equal, remove []*GroupMembers) {
	idpMemberSet := make(map[string]map[string]Member)
	scimMemberSet := make(map[string]map[string]Member)
	scimGroupsSet := make(map[string]Group)

	for _, grpMembers := range idp {
		idpMemberSet[grpMembers.Group.Name] = make(map[string]Member)
		for _, member := range grpMembers.Resources {
			idpMemberSet[grpMembers.Group.Name][member.Email] = *member
		}
	}

	for _, grpMembers := range scim {
		scimGroupsSet[grpMembers.Group.Name] = *grpMembers.Group
		scimMemberSet[grpMembers.Group.Name] = make(map[string]Member)
		for _, member := range grpMembers.Resources {
			scimMemberSet[grpMembers.Group.Name][member.Email] = *member
		}
	}

	toCreate := make([]*GroupMembers, 0)
	toEqual := make([]*GroupMembers, 0)
	toRemove := make([]*GroupMembers, 0)

	for _, grpMembers := range idp {
		toC := make(map[string][]*Member)
		toE := make(map[string][]*Member)

		toC[grpMembers.Group.Name] = make([]*Member, 0)
		toE[grpMembers.Group.Name] = make([]*Member, 0)

		// count when both side have members == 0
		noMembers := 0

		// groups equals both sides without members
		if _, ok := scimMemberSet[grpMembers.Group.Name]; ok {
			if len(scimMemberSet[grpMembers.Group.Name]) == 0 && len(idpMemberSet[grpMembers.Group.Name]) == 0 {
				noMembers++
			}
		}

		// this case is when the groups is not new in scim
		if grpMembers.Group.SCIMID == "" {
			if _, ok := scimGroupsSet[grpMembers.Group.Name]; ok {
				grpMembers.Group.SCIMID = scimGroupsSet[grpMembers.Group.Name].SCIMID
			}
		}

		for _, member := range grpMembers.Resources {
			if _, ok := scimMemberSet[grpMembers.Group.Name][member.Email]; !ok {
				toC[grpMembers.Group.Name] = append(toC[grpMembers.Group.Name], member)
			} else {
				// check if the groups has the same members before adding to equal
				// TODO: check if the groups has the same members before adding to equal, what happens if some members are different?
				for grpMemberEmail := range scimMemberSet[grpMembers.Group.Name] {
					if grpMemberEmail == member.Email {
						member.SCIMID = scimMemberSet[grpMembers.Group.Name][member.Email].SCIMID
						toE[grpMembers.Group.Name] = append(toE[grpMembers.Group.Name], member)
					}
				}
			}
		}

		if len(toC[grpMembers.Group.Name]) > 0 {
			grpMembers.Group.SetHashCode()

			e := GroupMembersBuilder().
				WithGroup(grpMembers.Group).
				WithResources(toC[grpMembers.Group.Name]).
				Build()

			toCreate = append(toCreate, e)
		}

		if noMembers > 0 || len(toE[grpMembers.Group.Name]) > 0 {
			grpMembers.Group.SetHashCode()

			ee := GroupMembersBuilder().
				WithGroup(grpMembers.Group).
				WithResources(toE[grpMembers.Group.Name]).
				Build()

			toEqual = append(toEqual, ee)
		}
	}

	for _, grpMembers := range scim {
		toD := make(map[string][]*Member)
		toD[grpMembers.Group.Name] = make([]*Member, 0)

		for _, member := range grpMembers.Resources {
			if _, ok := idpMemberSet[grpMembers.Group.Name][member.Email]; !ok {
				toD[grpMembers.Group.Name] = append(toD[grpMembers.Group.Name], member)
			}
		}

		if len(toD[grpMembers.Group.Name]) > 0 {
			grpMembers.Group.SetHashCode()

			e := GroupMembersBuilder().
				WithGroup(grpMembers.Group).
				WithResources(toD[grpMembers.Group.Name]).
				Build()

			toRemove = append(toRemove, e)
		}
	}

	return toCreate, toEqual, toRemove
}
