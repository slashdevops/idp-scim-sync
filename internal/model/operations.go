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

	create = &GroupsMembersResult{
		Items:     len(toCreate),
		Resources: toCreate,
	}
	create.SetHashCode()

	equal = &GroupsMembersResult{
		Items:     len(toEqual),
		Resources: toEqual,
	}
	equal.SetHashCode()

	remove = &GroupsMembersResult{
		Items:     len(toRemove),
		Resources: toRemove,
	}
	remove.SetHashCode()

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

	create = &GroupsResult{
		Items:     len(toCreate),
		Resources: toCreate,
	}
	create.SetHashCode()

	update = &GroupsResult{
		Items:     len(toUpdate),
		Resources: toUpdate,
	}
	update.SetHashCode()

	equal = &GroupsResult{
		Items:     len(toEqual),
		Resources: toEqual,
	}
	equal.SetHashCode()

	remove = &GroupsResult{
		Items:     len(toRemove),
		Resources: toRemove,
	}
	remove.SetHashCode()

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
		idpUsers[usr.Email] = struct{}{}
	}

	for _, usr := range scim.Resources {
		scimUsers[usr.Email] = *usr
	}

	// new users and what equal to them
	for _, usr := range idp.Resources {
		if _, ok := scimUsers[usr.Email]; !ok {
			toCreate = append(toCreate, usr)
		} else {
			usr.SCIMID = scimUsers[usr.Email].SCIMID

			// TODO: replace this check with the check of the hash code
			if usr.Name.FamilyName != scimUsers[usr.Email].Name.FamilyName ||
				usr.Name.GivenName != scimUsers[usr.Email].Name.GivenName ||
				usr.Active != scimUsers[usr.Email].Active || usr.IPID != scimUsers[usr.Email].IPID ||
				usr.DisplayName != scimUsers[usr.Email].DisplayName {
				toUpdate = append(toUpdate, usr)
			} else {
				toEqual = append(toEqual, usr)
			}
		}
	}

	for _, usr := range scim.Resources {
		if _, ok := idpUsers[usr.Email]; !ok {
			toRemove = append(toRemove, usr)
		}
	}

	create = &UsersResult{
		Items:     len(toCreate),
		Resources: toCreate,
	}
	create.SetHashCode()

	update = &UsersResult{
		Items:     len(toUpdate),
		Resources: toUpdate,
	}
	update.SetHashCode()

	equal = &UsersResult{
		Items:     len(toEqual),
		Resources: toEqual,
	}
	equal.SetHashCode()

	remove = &UsersResult{
		Items:     len(toRemove),
		Resources: toRemove,
	}
	remove.SetHashCode()

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

	merged = &GroupsResult{
		Items:     len(groups),
		Resources: groups,
	}
	merged.SetHashCode()

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

	merged = &UsersResult{
		Items:     len(users),
		Resources: users,
	}
	merged.SetHashCode()

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

	merged = &GroupsMembersResult{
		Items:     len(groupsMembers),
		Resources: groupsMembers,
	}
	merged.SetHashCode()

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
		users[user.Email] = *user
	}

	gms := make([]*GroupMembers, 0)
	for _, groupMembers := range idp.Resources {
		mbs := make([]*Member, 0)

		g := Group{
			IPID:   groupMembers.Group.IPID,
			SCIMID: groups[groupMembers.Group.Name].SCIMID,
			Name:   groupMembers.Group.Name,
			Email:  groupMembers.Group.Email,
		}
		g.SetHashCode()

		for _, member := range groupMembers.Resources {
			m := &Member{
				IPID:   member.IPID,
				SCIMID: users[member.Email].SCIMID,
				Email:  member.Email,
				Status: member.Status,
			}
			m.SetHashCode()
			mbs = append(mbs, m)
		}

		gm := &GroupMembers{
			Items:     len(mbs),
			Group:     g,
			Resources: mbs,
		}
		gm.SetHashCode()

		gms = append(gms, gm)
	}

	gmr := &GroupsMembersResult{
		Items:     idp.Items,
		Resources: gms,
	}
	gmr.SetHashCode()

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
		scimGroupsSet[grpMembers.Group.Name] = grpMembers.Group
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

			e := &GroupMembers{
				Items:     len(toC[grpMembers.Group.Name]),
				Group:     grpMembers.Group,
				Resources: toC[grpMembers.Group.Name],
			}
			e.SetHashCode()
			toCreate = append(toCreate, e)
		}

		if noMembers > 0 || len(toE[grpMembers.Group.Name]) > 0 {
			grpMembers.Group.SetHashCode()
			ee := &GroupMembers{
				Items:     len(toE[grpMembers.Group.Name]),
				Group:     grpMembers.Group,
				Resources: toE[grpMembers.Group.Name],
			}
			ee.SetHashCode()
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

			e := &GroupMembers{
				Items:     len(toD[grpMembers.Group.Name]),
				Group:     grpMembers.Group,
				Resources: toD[grpMembers.Group.Name],
			}
			e.SetHashCode()
			toRemove = append(toRemove, e)
		}
	}

	return toCreate, toEqual, toRemove
}
