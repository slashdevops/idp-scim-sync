package core

import (
	log "github.com/sirupsen/logrus"
	"github.com/slashdevops/idp-scim-sync/internal/hash"
	"github.com/slashdevops/idp-scim-sync/internal/model"
	"github.com/slashdevops/idp-scim-sync/internal/utils"
)

// groupsOperations returns the differences between the groups in the
// this use the Groups Name as the key.
// SCIM Groups cannot be updated.
// return 4 objest of GroupsResult
// create: groups that exist in "state" but not in "idp"
// update: groups that exist in "idp" and in "state" but attributes changed in idp
// equal: groups that exist in both "idp" and "state" and their attributes are equal
// delete: groups that exist in "state" but not in "idp"
func groupsOperations(idp, state *model.GroupsResult) (create *model.GroupsResult, update *model.GroupsResult, equal *model.GroupsResult, delete *model.GroupsResult) {
	idpGroups := make(map[string]struct{})
	stateGroups := make(map[string]model.Group)

	toCreate := make([]model.Group, 0)
	toUpdate := make([]model.Group, 0)
	toEqual := make([]model.Group, 0)
	toDelete := make([]model.Group, 0)

	for _, gr := range idp.Resources {
		idpGroups[gr.Name] = struct{}{}
	}

	for _, gr := range state.Resources {
		stateGroups[gr.Name] = gr
	}

	// new groups and what equal to them
	for _, gr := range idp.Resources {
		if _, ok := stateGroups[gr.Name]; !ok {
			toCreate = append(toCreate, gr)
		} else {
			// Check if the group IPID is not the same
			if gr.IPID != stateGroups[gr.Name].IPID {

				// log.WithFields(log.Fields{
				// 	"idp":  gr.IPID,
				// 	"scim": gr.SCIMID,
				// }).Debugf("group to update, payload: %s", utils.ToJSON(gr))

				gr.SCIMID = stateGroups[gr.Name].SCIMID

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

// usersOperations returns the differences between the users in the
// Users Email as the key.
// return 4 objest of UsersResult
// create: users that exist in "state" but not in "idp"
// update: users that exist in "idp" and in "state" but attributes changed in idp
// equal: users that exist in both "idp" and "state" and their attributes are equal
// delete: users that exist in "state" but not in "idp"
func usersOperations(idp, state *model.UsersResult) (create *model.UsersResult, update *model.UsersResult, equal *model.UsersResult, delete *model.UsersResult) {
	idpUsers := make(map[string]struct{})
	stateUsers := make(map[string]model.User)

	toCreate := make([]model.User, 0)
	toUpdate := make([]model.User, 0)
	toEqual := make([]model.User, 0)
	toDelete := make([]model.User, 0)

	// log.Debugf("payloads: user idp: %s, user scim: %s", utils.ToJSON(idp), utils.ToJSON(state))

	for _, usr := range idp.Resources {
		idpUsers[usr.Email] = struct{}{}
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
				usr.IPID != stateUsers[usr.Email].IPID {

				// log.WithFields(log.Fields{
				// 	"idp":  usr.IPID,
				// 	"scim": usr.SCIMID,
				// }).Debugf("user to update, payload: %s", utils.ToJSON(usr))

				usr.SCIMID = stateUsers[usr.Email].SCIMID

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

// groupsUsersOperationsSCIM returns the differences between the users in the
// Users Email as the key.
// return 3 objest of GroupsUsersResult
// create: users that exist in "state" but not in "idp"
// equal: users that exist in both "idp" and "state" and their attributes are equal
// delete: users that exist in "state" but not in "idp"
//
// NOTE: state has the groups and the users that belong to them in the SCIM side, in case the groups doesn't have users,
// the resources of the groups are empty.
func groupsUsersOperationsSCIM(idp, state *model.GroupsUsersResult, ur *model.UsersResult) (create *model.GroupsUsersResult, equal *model.GroupsUsersResult, delete *model.GroupsUsersResult) {
	idpUsers := make(map[string]map[string]model.User)
	stateUsers := make(map[string]map[string]model.User)
	stateGroups := make(map[string]model.Group)

	// This is necessary to fill the scimid of the users when the existing goups in scim side dosen't have members
	// and we need to create these users
	scimUsers := make(map[string]model.User)

	toCreate := make([]model.GroupUsers, 0)
	toEqual := make([]model.GroupUsers, 0)
	toDelete := make([]model.GroupUsers, 0)

	// fill wilth the email of the users because id the unique identifier of the user and because
	// at this point the IDPID is empty
	for _, uSCIM := range ur.Resources {
		scimUsers[uSCIM.Email] = uSCIM
	}

	// log.Debugf("scimUsers: %s", utils.ToJSON(scimUsers))

	for _, grpUsrs := range idp.Resources {
		idpUsers[grpUsrs.Group.IPID] = make(map[string]model.User)

		for _, usr := range grpUsrs.Resources {

			if _, ok := scimUsers[usr.Email]; ok {
				usr.SCIMID = scimUsers[usr.Email].SCIMID
			}

			idpUsers[grpUsrs.Group.IPID][usr.IPID] = usr

		}
	}

	for _, grpUsrs := range state.Resources {
		stateGroups[grpUsrs.Group.IPID] = grpUsrs.Group
		stateUsers[grpUsrs.Group.IPID] = make(map[string]model.User)
		for _, usr := range grpUsrs.Resources {
			stateUsers[grpUsrs.Group.IPID][usr.IPID] = usr
		}
	}

	// log.Debugf("idpUsers: %s", utils.ToJSON(idpUsers))
	// log.Debugf("stateUsers: %s", utils.ToJSON(stateUsers))
	// log.Debugf("state: %s", utils.ToJSON(state))

	// map[group.ID][]*model.User
	toC := make(map[string][]model.User)
	toE := make(map[string][]model.User)
	toD := make(map[string][]model.User)

	for _, grpUsrs := range idp.Resources {
		toC[grpUsrs.Group.IPID] = make([]model.User, 0)
		toE[grpUsrs.Group.IPID] = make([]model.User, 0)

		grpUsrs.Group.SCIMID = stateGroups[grpUsrs.Group.IPID].SCIMID

		for _, usr := range grpUsrs.Resources {
			usr.SCIMID = idpUsers[grpUsrs.Group.IPID][usr.IPID].SCIMID

			if _, ok := stateUsers[grpUsrs.Group.IPID][usr.IPID]; !ok {
				toC[grpUsrs.Group.IPID] = append(toC[grpUsrs.Group.IPID], usr)
			} else {
				toE[grpUsrs.Group.IPID] = append(toE[grpUsrs.Group.IPID], usr)
			}
		}

		if len(toC[grpUsrs.Group.IPID]) > 0 {
			toCreate = append(toCreate, model.GroupUsers{
				Items:     len(toC[grpUsrs.Group.IPID]),
				Group:     grpUsrs.Group,
				Resources: toC[grpUsrs.Group.IPID],
			})
		}

		if len(toE[grpUsrs.Group.IPID]) > 0 {
			toEqual = append(toEqual, model.GroupUsers{
				Items:     len(toE[grpUsrs.Group.IPID]),
				Group:     grpUsrs.Group,
				Resources: toE[grpUsrs.Group.IPID],
			})
		}
	}

	for _, grpUsrs := range state.Resources {
		toD[grpUsrs.Group.IPID] = make([]model.User, 0)

		for _, usr := range grpUsrs.Resources {
			if _, ok := idpUsers[grpUsrs.Group.IPID][usr.IPID]; !ok {
				toD[grpUsrs.Group.IPID] = append(toD[grpUsrs.Group.IPID], usr)
			}
		}

		if len(toD[grpUsrs.Group.IPID]) > 0 {
			toDelete = append(toDelete, model.GroupUsers{
				Items:     len(toD[grpUsrs.Group.IPID]),
				Group:     grpUsrs.Group,
				Resources: toD[grpUsrs.Group.IPID],
			})
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

// groupsUsersOperationsState returns the differences between the users in the
// Users Email as the key.
// return 3 objest of GroupsUsersResult
// create: users that exist in "state" but not in "idp"
// equal: users that exist in both "idp" and "state" and their attributes are equal
// delete: users that exist in "state" but not in "idp"
//
// NOTE: state has the groups and the users that belong to them in the SCIM side, in case the groups doesn't have users,
// the resources of the groups are empty.
func groupsUsersOperationsState(idp, state *model.GroupsUsersResult) (create *model.GroupsUsersResult, equal *model.GroupsUsersResult, delete *model.GroupsUsersResult) {
	idpUsers := make(map[string]map[string]model.User)
	stateUsers := make(map[string]map[string]model.User)
	stateGroups := make(map[string]model.Group)

	toCreate := make([]model.GroupUsers, 0)
	toEqual := make([]model.GroupUsers, 0)
	toDelete := make([]model.GroupUsers, 0)

	for _, grpUsrs := range idp.Resources {
		idpUsers[grpUsrs.Group.IPID] = make(map[string]model.User)

		for _, usr := range grpUsrs.Resources {
			idpUsers[grpUsrs.Group.IPID][usr.IPID] = usr
		}
	}

	for _, grpUsrs := range state.Resources {
		stateGroups[grpUsrs.Group.IPID] = grpUsrs.Group
		stateUsers[grpUsrs.Group.IPID] = make(map[string]model.User)
		for _, usr := range grpUsrs.Resources {
			stateUsers[grpUsrs.Group.IPID][usr.IPID] = usr
		}
	}

	// log.Debugf("idpUsers: %s", utils.ToJSON(idpUsers))
	// log.Debugf("stateUsers: %s", utils.ToJSON(stateUsers))
	log.Debugf("state: %s", utils.ToJSON(state))

	// map[group.ID][]*model.User
	toC := make(map[string][]model.User)
	toE := make(map[string][]model.User)
	toD := make(map[string][]model.User)

	for _, grpUsrs := range idp.Resources {
		toC[grpUsrs.Group.IPID] = make([]model.User, 0)
		toE[grpUsrs.Group.IPID] = make([]model.User, 0)

		grpUsrs.Group.SCIMID = stateGroups[grpUsrs.Group.IPID].SCIMID

		for _, usr := range grpUsrs.Resources {
			usr.SCIMID = stateUsers[grpUsrs.Group.IPID][usr.IPID].SCIMID
			log.Debugf("scim user: %s,idp user: %s", utils.ToJSON(stateUsers[grpUsrs.Group.IPID]), utils.ToJSON(usr))

			if _, ok := stateUsers[grpUsrs.Group.IPID][usr.IPID]; !ok {
				toC[grpUsrs.Group.IPID] = append(toC[grpUsrs.Group.IPID], usr)
			} else {
				toE[grpUsrs.Group.IPID] = append(toE[grpUsrs.Group.IPID], usr)
			}
		}

		if len(toC[grpUsrs.Group.IPID]) > 0 {
			toCreate = append(toCreate, model.GroupUsers{
				Items:     len(toC[grpUsrs.Group.IPID]),
				Group:     grpUsrs.Group,
				Resources: toC[grpUsrs.Group.IPID],
			})
		}

		if len(toE[grpUsrs.Group.IPID]) > 0 {
			toEqual = append(toEqual, model.GroupUsers{
				Items:     len(toE[grpUsrs.Group.IPID]),
				Group:     grpUsrs.Group,
				Resources: toE[grpUsrs.Group.IPID],
			})
		}
	}

	for _, grpUsrs := range state.Resources {
		toD[grpUsrs.Group.IPID] = make([]model.User, 0)

		for _, usr := range grpUsrs.Resources {
			if _, ok := idpUsers[grpUsrs.Group.IPID][usr.IPID]; !ok {
				toD[grpUsrs.Group.IPID] = append(toD[grpUsrs.Group.IPID], usr)
			}
		}

		if len(toD[grpUsrs.Group.IPID]) > 0 {
			toDelete = append(toDelete, model.GroupUsers{
				Items:     len(toD[grpUsrs.Group.IPID]),
				Group:     grpUsrs.Group,
				Resources: toD[grpUsrs.Group.IPID],
			})
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

func mergeGroupsResult(grs ...*model.GroupsResult) (merged model.GroupsResult) {
	groups := make([]model.Group, 0)

	for _, gr := range grs {
		groups = append(groups, gr.Resources...)
	}

	merged = model.GroupsResult{
		Items:     len(groups),
		Resources: groups,
	}
	if len(merged.Resources) > 0 {
		merged.HashCode = hash.Get(merged)
	}

	return
}

func mergeUsersResult(urs ...*model.UsersResult) (merged model.UsersResult) {
	users := make([]model.User, 0)

	for _, u := range urs {
		users = append(users, u.Resources...)
	}

	merged = model.UsersResult{
		Items:     len(users),
		Resources: users,
	}
	if len(merged.Resources) > 0 {
		merged.HashCode = hash.Get(merged)
	}

	return
}

func mergeGroupsUsersResult(gurs ...*model.GroupsUsersResult) (merged model.GroupsUsersResult) {
	groupsUsers := make([]model.GroupUsers, 0)

	for _, gu := range gurs {
		groupsUsers = append(groupsUsers, gu.Resources...)
	}

	merged = model.GroupsUsersResult{
		Items:     len(groupsUsers),
		Resources: groupsUsers,
	}
	if len(merged.Resources) > 0 {
		merged.HashCode = hash.Get(merged)
	}

	return
}
