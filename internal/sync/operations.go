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

	// new groups and what equal to them
	for _, usr := range state.Resources {
		if _, ok := idpUsers[usr.Email]; !ok {

			// Check if the group email changed
			// Id changed happen when the group delete and create again with the same name and email I guest
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
