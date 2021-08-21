package sync

func createSyncState(sgr *StoreGroupsResult, sgmr *StoreGroupsMembersResult, sur *StoreUsersResult) (SyncState, error) {
	return SyncState{
		Version:  "1.0.0",
		Checksum: "TBD",
	}, nil
}

// groupsDifferences returns the differences between the groups in the
// this use the Groups Name as the key.
// SCIM Groups cannot be updated.
// return 4 objest of GroupsResult
// create: groups that exist in "state" but not in "idp"
// update: groups that exist in "idp" and in "state" but attributes changed in idp
// equal: groups that exist in both "idp" and "state" and their attributes are equal
// delete: groups that exist in "state" but not in "idp"
func groupsDifferences(idp, state *GroupsResult) (create *GroupsResult, update *GroupsResult, equal *GroupsResult, delete *GroupsResult) {

	idpGroups := make(map[string]*Group)
	stateGroups := make(map[string]struct{})

	toCreate := make([]*Group, 0)
	toUpdate := make([]*Group, 0)
	toEqual := make([]*Group, 0)
	toDelete := make([]*Group, 0)

	for _, gr := range idp.Resources {
		idpGroups[gr.Name] = gr
	}

	for _, gr := range state.Resources {
		stateGroups[gr.Name] = struct{}{}
	}

	// new groups and what equal to them
	for _, gr := range state.Resources {
		if _, ok := idpGroups[gr.Name]; !ok {

			// Check if the group email changed
			// Id changed happen when the group delete and create again with the same name and email I guest
			if gr.Email != idpGroups[gr.Name].Email || gr.Id != idpGroups[gr.Name].Id {
				toUpdate = append(toUpdate, gr)
			} else {
				toCreate = append(toCreate, gr)
			}

		} else {
			toEqual = append(toEqual, gr)
		}
	}

	for _, gr := range idp.Resources {
		if _, ok := stateGroups[gr.Name]; !ok {
			toDelete = append(toDelete, gr)
		}
	}

	create = &GroupsResult{
		Items:     len(toCreate),
		Resources: toCreate,
	}

	update = &GroupsResult{
		Items:     len(toUpdate),
		Resources: toUpdate,
	}

	equal = &GroupsResult{
		Items:     len(toEqual),
		Resources: toEqual,
	}

	delete = &GroupsResult{
		Items:     len(toDelete),
		Resources: toDelete,
	}

	return
}

// usersDifferences returns the differences between the users in the
// Users Email as the key.
// return 4 objest of UsersResult
// create: users that exist in "state" but not in "idp"
// update: users that exist in "idp" and in "state" but attributes changed in idp
// equal: users that exist in both "idp" and "state" and their attributes are equal
// delete: users that exist in "state" but not in "idp"
func usersDifferences(idp, state *UsersResult) (create *UsersResult, update *UsersResult, equal *UsersResult, delete *UsersResult) {

	idpUsers := make(map[string]*User)
	stateUsers := make(map[string]struct{})

	toCreate := make([]*User, 0)
	toUpdate := make([]*User, 0)
	toEqual := make([]*User, 0)
	toDelete := make([]*User, 0)

	for _, usr := range idp.Resources {
		idpUsers[usr.Email] = usr
	}

	for _, usr := range state.Resources {
		stateUsers[usr.Email] = struct{}{}
	}

	// new users and what equal to them
	for _, usr := range state.Resources {
		if _, ok := idpUsers[usr.Email]; !ok {

			// Check if the user fields changed
			if usr.Name.FamilyName != idpUsers[usr.Email].Name.FamilyName ||
				usr.Name.GivenName != idpUsers[usr.Email].Name.GivenName ||
				usr.Active != idpUsers[usr.Email].Active ||
				usr.Id != idpUsers[usr.Email].Id {
				toUpdate = append(toUpdate, usr)
			} else {
				toCreate = append(toCreate, usr)
			}

		} else {
			toEqual = append(toEqual, usr)
		}
	}

	for _, usr := range idp.Resources {
		if _, ok := stateUsers[usr.Email]; !ok {
			toDelete = append(toDelete, usr)
		}
	}

	create = &UsersResult{
		Items:     len(toCreate),
		Resources: toCreate,
	}

	update = &UsersResult{
		Items:     len(toUpdate),
		Resources: toUpdate,
	}

	equal = &UsersResult{
		Items:     len(toEqual),
		Resources: toEqual,
	}

	delete = &UsersResult{
		Items:     len(toDelete),
		Resources: toDelete,
	}

	return
}

// groupsMembersDifferences returns the differences between the members in the
// Users Email as the key.
// return 4 objest of GroupsMembersResult
// create: users that exist in "state" but not in "idp"
// equal: users that exist in both "idp" and "state" and their attributes are equal
// delete: users that exist in "state" but not in "idp"
func groupsMembersDifferences(idp, state *GroupsMembersResult) (create *GroupsMembersResult, equal *GroupsMembersResult, delete *GroupsMembersResult) {

	idpMembers := make(map[string]map[string]struct{})
	stateMembers := make(map[string]map[string]struct{})

	toCreate := make(GroupsMembers)
	toEqual := make(GroupsMembers)
	toDelete := make(GroupsMembers)

	for grpName, mbrs := range *idp.Resources {
		idpMembers[grpName] = make(map[string]struct{})
		for _, mbr := range mbrs {
			idpMembers[grpName][mbr.Email] = struct{}{}
		}
	}

	for grpName, mbrs := range *state.Resources {
		stateMembers[grpName] = make(map[string]struct{})
		for _, mbr := range mbrs {
			stateMembers[grpName][mbr.Email] = struct{}{}
		}
	}

	for grpName, mbrs := range *state.Resources {
		for _, mbr := range mbrs {
			if _, ok := idpMembers[grpName][mbr.Email]; !ok {
				toDelete[grpName] = append(toDelete[grpName], mbr)
			} else {
				toEqual[grpName] = append(toEqual[grpName], mbr)
			}
		}
	}

	for grpName, mbrs := range *idp.Resources {
		for _, mbr := range mbrs {
			if _, ok := stateMembers[grpName][mbr.Email]; !ok {
				toCreate[grpName] = append(toCreate[grpName], mbr)
			}
		}
	}

	create = &GroupsMembersResult{
		Items:     len(toCreate),
		Resources: &toCreate,
	}

	equal = &GroupsMembersResult{
		Items:     len(toEqual),
		Resources: &toEqual,
	}

	delete = &GroupsMembersResult{
		Items:     len(toDelete),
		Resources: &toDelete,
	}

	return
}
