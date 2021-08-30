package core

import (
	"github.com/slashdevops/idp-scim-sync/internal/model"
)

func createSyncState(sgr *model.StoreGroupsResult, sgmr *model.StoreGroupsUsersResult, sur *model.StoreUsersResult) (model.SyncState, error) {
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
	stateGroups := make(map[string]*model.Group)

	toCreate := make([]*model.Group, 0)
	toUpdate := make([]*model.Group, 0)
	toEqual := make([]*model.Group, 0)
	toDelete := make([]*model.Group, 0)

	for _, gr := range idp.Resources {
		idpGroups[gr.Name] = gr
	}

	for _, gr := range state.Resources {
		stateGroups[gr.Name] = gr
	}

	// new groups and what equal to them
	for _, gr := range idp.Resources {
		if _, ok := stateGroups[gr.Name]; !ok {
			toCreate = append(toCreate, gr)
		} else {
			// Check if the group email or ID changed
			// Id changed happen when the group delete and create again with the same name and email I guest
			if gr.Email != stateGroups[gr.Name].Email || gr.ID != stateGroups[gr.Name].ID {
				toUpdate = append(toUpdate, gr)
			} else {
				toEqual = append(toEqual, gr)
			}
		}
	}

	for _, gr := range state.Resources {
		if _, ok := idpGroups[gr.Name]; !ok {
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
	stateUsers := make(map[string]*model.User)

	toCreate := make([]*model.User, 0)
	toUpdate := make([]*model.User, 0)
	toEqual := make([]*model.User, 0)
	toDelete := make([]*model.User, 0)

	for _, usr := range idp.Resources {
		idpUsers[usr.Email] = usr
	}

	for _, usr := range state.Resources {
		stateUsers[usr.Email] = usr
	}

	// new users and what equal to them
	for _, usr := range idp.Resources {
		if _, ok := stateUsers[usr.Email]; !ok {
			toCreate = append(toCreate, usr)
		} else {
			// Check if the user fields changed
			if usr.Name.FamilyName != stateUsers[usr.Email].Name.FamilyName ||
				usr.Name.GivenName != stateUsers[usr.Email].Name.GivenName ||
				usr.Active != stateUsers[usr.Email].Active ||
				usr.ID != stateUsers[usr.Email].ID {
				toUpdate = append(toUpdate, usr)
			} else {
				toEqual = append(toEqual, usr)
			}
		}
	}

	for _, usr := range state.Resources {
		if _, ok := idpUsers[usr.Email]; !ok {
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

// groupsUsersDifferences returns the differences between the users in the
// Users Email as the key.
// return 3 objest of GroupsUsersResult
// create: users that exist in "state" but not in "idp"
// equal: users that exist in both "idp" and "state" and their attributes are equal
// delete: users that exist in "state" but not in "idp"
func groupsUsersDifferences(idp, state *model.GroupsUsersResult) (create *model.GroupsUsersResult, equal *model.GroupsUsersResult, delete *model.GroupsUsersResult) {
	idpUsers := make(map[string]map[string]*model.User)
	stateUsers := make(map[string]map[string]*model.User)

	toCreate := make([]*model.GroupUsers, 0)
	toEqual := make([]*model.GroupUsers, 0)
	toDelete := make([]*model.GroupUsers, 0)

	for _, grpUsrs := range idp.Resources {
		idpUsers[grpUsrs.Group.ID] = make(map[string]*model.User)
		for _, usr := range grpUsrs.Resources {
			idpUsers[grpUsrs.Group.ID][usr.ID] = usr
		}
	}

	for _, grpUsrs := range state.Resources {
		stateUsers[grpUsrs.Group.ID] = make(map[string]*model.User)
		for _, usr := range grpUsrs.Resources {
			stateUsers[grpUsrs.Group.ID][usr.ID] = usr
		}
	}

	for _, grpUsrs := range idp.Resources {
		if _, ok := stateUsers[grpUsrs.Group.ID]; !ok {
			toCreate = append(toCreate, grpUsrs)
		} else {
			toEqual = append(toEqual, grpUsrs)
		}
	}

	for _, grpUsrs := range state.Resources {
		if _, ok := idpUsers[grpUsrs.Group.ID]; !ok {
			toDelete = append(toDelete, grpUsrs)
		}
	}

	create = &model.GroupsUsersResult{
		Items:     len(toCreate),
		Resources: toCreate,
	}

	equal = &model.GroupsUsersResult{
		Items:     len(toEqual),
		Resources: toEqual,
	}

	delete = &model.GroupsUsersResult{
		Items:     len(toDelete),
		Resources: toDelete,
	}

	return
}
