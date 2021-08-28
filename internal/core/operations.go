package core

import "github.com/slashdevops/idp-scim-sync/internal/model"

func createSyncState(sgr *model.StoreGroupsResult, sgmr *model.StoreGroupsMembersResult, sur *model.StoreUsersResult) (model.SyncState, error) {
	return model.SyncState{
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
func groupsDifferences(idp, state *model.GroupsResult) (create *model.GroupsResult, update *model.GroupsResult, equal *model.GroupsResult, delete *model.GroupsResult) {
	idpGroups := make(map[string]*model.Group)
	stateGroups := make(map[string]struct{})

	toCreate := make([]*model.Group, 0)
	toUpdate := make([]*model.Group, 0)
	toEqual := make([]*model.Group, 0)
	toDelete := make([]*model.Group, 0)

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
			if gr.Email != idpGroups[gr.Name].Email || gr.ID != idpGroups[gr.Name].ID {
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

	create = &model.GroupsResult{
		Items:     len(toCreate),
		Resources: toCreate,
	}

	update = &model.GroupsResult{
		Items:     len(toUpdate),
		Resources: toUpdate,
	}

	equal = &model.GroupsResult{
		Items:     len(toEqual),
		Resources: toEqual,
	}

	delete = &model.GroupsResult{
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
func usersDifferences(idp, state *model.UsersResult) (create *model.UsersResult, update *model.UsersResult, equal *model.UsersResult, delete *model.UsersResult) {
	idpUsers := make(map[string]*model.User)
	stateUsers := make(map[string]struct{})

	toCreate := make([]*model.User, 0)
	toUpdate := make([]*model.User, 0)
	toEqual := make([]*model.User, 0)
	toDelete := make([]*model.User, 0)

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
				usr.ID != idpUsers[usr.Email].ID {
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

	create = &model.UsersResult{
		Items:     len(toCreate),
		Resources: toCreate,
	}

	update = &model.UsersResult{
		Items:     len(toUpdate),
		Resources: toUpdate,
	}

	equal = &model.UsersResult{
		Items:     len(toEqual),
		Resources: toEqual,
	}

	delete = &model.UsersResult{
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
func groupsMembersDifferences(idp, state *model.GroupsMembersResult) (create *model.GroupsMembersResult, equal *model.GroupsMembersResult, delete *model.GroupsMembersResult) {
	idpMembers := make(map[string]map[string]struct{})
	stateMembers := make(map[string]map[string]struct{})

	toCreate := make(model.GroupsMembers)
	toEqual := make(model.GroupsMembers)
	toDelete := make(model.GroupsMembers)

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

	create = &model.GroupsMembersResult{
		Items:     len(toCreate),
		Resources: &toCreate,
	}

	equal = &model.GroupsMembersResult{
		Items:     len(toEqual),
		Resources: &toEqual,
	}

	delete = &model.GroupsMembersResult{
		Items:     len(toDelete),
		Resources: &toDelete,
	}

	return
}
