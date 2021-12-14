package core

import (
	"errors"

	"github.com/slashdevops/idp-scim-sync/internal/model"
)

var (
	ErrIdentityProviderGroupsMembersNil = errors.New("identity provider groups members is nil")
	ErrSCIMGroupsMembersNil             = errors.New("scim groups members is nil")
	ErrIdentityProviderGroupsNil        = errors.New("identity provider groups is nil")
	ErrSCIMGroupsNil                    = errors.New("scim groups is nil")
	ErrIdentityProviderUsersNil         = errors.New("identity provider users is nil")
	ErrSCIMUsersNil                     = errors.New("scim users is nil")
)

// groupsOperations returns datasets used to perform differents operations over the SCIM side
// return 4 objest of GroupsMembersResult
// create: groups that exist in "idp" but not in "scim" or "state"
// update: groups that exist in "idp" and in "scim" or "state" but attributes changed in idp
// equal: groups that exist in both "idp" and "scim" or "state" and their attributes are equal
// delete: groups that exist in "scim" or "state" but not in "idp"
//
// also this extract the id from scim to fill the resutls
func membersOperations(idp, scim *model.GroupsMembersResult) (create *model.GroupsMembersResult, equal *model.GroupsMembersResult, delete *model.GroupsMembersResult, err error) {
	if idp == nil {
		create, equal, delete, err = nil, nil, nil, ErrIdentityProviderGroupsMembersNil
		return
	}
	if scim == nil {
		create, equal, delete, err = nil, nil, nil, ErrSCIMGroupsMembersNil
		return
	}

	idpMemberSet := make(map[string]map[string]model.Member)  // [idp.GroupMembers.Group.Name] -> idp.member.Email -> idp.member
	scimMemberSet := make(map[string]map[string]model.Member) // [scim.Group.Name] -> [scim.member.Email] -> scim.member
	scimGroupsSet := make(map[string]model.Group)             // [scim.Group.Name] -> [scim.Group]

	toCreate := make([]*model.GroupMembers, 0)
	toEqual := make([]*model.GroupMembers, 0)
	toDelete := make([]*model.GroupMembers, 0)

	for _, grpMembers := range idp.Resources {
		idpMemberSet[grpMembers.Group.Name] = make(map[string]model.Member)
		for _, member := range grpMembers.Resources {
			idpMemberSet[grpMembers.Group.Name][member.Email] = *member
		}
	}

	for _, grpMembers := range scim.Resources {
		scimGroupsSet[grpMembers.Group.Name] = grpMembers.Group
		scimMemberSet[grpMembers.Group.Name] = make(map[string]model.Member)
		for _, member := range grpMembers.Resources {
			scimMemberSet[grpMembers.Group.Name][member.Email] = *member
		}
	}

	for _, grpMembers := range idp.Resources {
		toC := make(map[string][]*model.Member)
		toE := make(map[string][]*model.Member)

		toC[grpMembers.Group.Name] = make([]*model.Member, 0)
		toE[grpMembers.Group.Name] = make([]*model.Member, 0)

		// count when both side have members == 0
		noMembers := 0

		// groups equals both sides without members
		if _, ok := scimMemberSet[grpMembers.Group.Name]; ok {
			if len(scimMemberSet[grpMembers.Group.Name]) == 0 && len(idpMemberSet[grpMembers.Group.Name]) == 0 {
				noMembers += 1
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

			e := &model.GroupMembers{
				Items:     len(toC[grpMembers.Group.Name]),
				Group:     grpMembers.Group,
				Resources: toC[grpMembers.Group.Name],
			}
			e.SetHashCode()

			toCreate = append(toCreate, e)
		}

		if noMembers > 0 || len(toE[grpMembers.Group.Name]) > 0 {
			grpMembers.Group.SetHashCode()
			ee := &model.GroupMembers{
				Items:     len(toE[grpMembers.Group.Name]),
				Group:     grpMembers.Group,
				Resources: toE[grpMembers.Group.Name],
			}
			ee.SetHashCode()

			toEqual = append(toEqual, ee)
		}
	}

	for _, grpMembers := range scim.Resources {
		toD := make(map[string][]*model.Member)
		toD[grpMembers.Group.Name] = make([]*model.Member, 0)

		for _, member := range grpMembers.Resources {
			if _, ok := idpMemberSet[grpMembers.Group.Name][member.Email]; !ok {
				toD[grpMembers.Group.Name] = append(toD[grpMembers.Group.Name], member)
			}
		}

		if len(toD[grpMembers.Group.Name]) > 0 {
			grpMembers.Group.SetHashCode()

			e := &model.GroupMembers{
				Items:     len(toD[grpMembers.Group.Name]),
				Group:     grpMembers.Group,
				Resources: toD[grpMembers.Group.Name],
			}
			e.SetHashCode()

			toDelete = append(toDelete, e)
		}
	}

	create = &model.GroupsMembersResult{
		Items:     len(toCreate),
		Resources: toCreate,
	}
	create.SetHashCode()

	equal = &model.GroupsMembersResult{
		Items:     len(toEqual),
		Resources: toEqual,
	}
	equal.SetHashCode()

	delete = &model.GroupsMembersResult{
		Items:     len(toDelete),
		Resources: toDelete,
	}
	delete.SetHashCode()

	return
}

// groupsOperations returns the differences between the groups in the
// this use the Groups Name as the key.
// SCIM Groups cannot be updated.
// return 4 objest of GroupsResult
// create: groups that exist in "idp" but not in "scim" or "state"
// update: groups that exist in "idp" and in "scim" or "state" but attributes changed in idp
// equal: groups that exist in both "idp" and "scim" or "state" and their attributes are equal
// delete: groups that exist in "scim" or "state" but not in "idp"
//
// also this extract the id from scim to fill the resutls
func groupsOperations(idp, scim *model.GroupsResult) (create *model.GroupsResult, update *model.GroupsResult, equal *model.GroupsResult, delete *model.GroupsResult, err error) {
	if idp == nil {
		create, update, equal, delete, err = nil, nil, nil, nil, ErrIdentityProviderGroupsNil
		return
	}
	if scim == nil {
		create, update, equal, delete, err = nil, nil, nil, nil, ErrSCIMGroupsNil
		return
	}

	idpGroups := make(map[string]struct{})     // [idp.Group.Name ] -> struct{}{}
	scimGroups := make(map[string]model.Group) // [scim.Group.Name] -> scim.Group

	toCreate := make([]*model.Group, 0)
	toUpdate := make([]*model.Group, 0)
	toEqual := make([]*model.Group, 0)
	toDelete := make([]*model.Group, 0)

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

	// loop over scim to see what to delete
	for _, group := range scim.Resources {
		if _, ok := idpGroups[group.Name]; !ok {
			toDelete = append(toDelete, group)
		}
	}

	create = &model.GroupsResult{
		Items:     len(toCreate),
		Resources: toCreate,
	}
	create.SetHashCode()

	update = &model.GroupsResult{
		Items:     len(toUpdate),
		Resources: toUpdate,
	}
	update.SetHashCode()

	equal = &model.GroupsResult{
		Items:     len(toEqual),
		Resources: toEqual,
	}
	equal.SetHashCode()

	delete = &model.GroupsResult{
		Items:     len(toDelete),
		Resources: toDelete,
	}
	delete.SetHashCode()

	return
}

// usersOperations returns datasets used to perform differents operations over the SCIM side
// return 4 objest of UsersResult
// create: users that exist in "idp" but not in "scim" or "state"
// update: users that exist in "idp" and in "scim" or "state" but attributes changed in idp
// equal: users that exist in both "idp" and "scim" or "state" and their attributes are equal
// delete: users that exist in "scim" or "state" but not in "idp"
func usersOperations(idp, scim *model.UsersResult) (create *model.UsersResult, update *model.UsersResult, equal *model.UsersResult, delete *model.UsersResult, err error) {
	if idp == nil {
		create, update, equal, delete, err = nil, nil, nil, nil, ErrIdentityProviderUsersNil
		return
	}
	if scim == nil {
		create, update, equal, delete, err = nil, nil, nil, nil, ErrSCIMUsersNil
		return
	}

	idpUsers := make(map[string]struct{})
	scimUsers := make(map[string]model.User)

	toCreate := make([]*model.User, 0)
	toUpdate := make([]*model.User, 0)
	toEqual := make([]*model.User, 0)
	toDelete := make([]*model.User, 0)

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

			if usr.Name.FamilyName != scimUsers[usr.Email].Name.FamilyName ||
				usr.Name.GivenName != scimUsers[usr.Email].Name.GivenName ||
				usr.Active != scimUsers[usr.Email].Active ||
				usr.IPID != scimUsers[usr.Email].IPID {
				toUpdate = append(toUpdate, usr)
			} else {
				toEqual = append(toEqual, usr)
			}
		}
	}

	for _, usr := range scim.Resources {
		if _, ok := idpUsers[usr.Email]; !ok {
			toDelete = append(toDelete, usr)
		}
	}

	create = &model.UsersResult{
		Items:     len(toCreate),
		Resources: toCreate,
	}
	create.SetHashCode()

	update = &model.UsersResult{
		Items:     len(toUpdate),
		Resources: toUpdate,
	}
	update.SetHashCode()

	equal = &model.UsersResult{
		Items:     len(toEqual),
		Resources: toEqual,
	}
	equal.SetHashCode()

	delete = &model.UsersResult{
		Items:     len(toDelete),
		Resources: toDelete,
	}
	delete.SetHashCode()

	return
}

// mergeGroupsResult merges n GroupsResult result
// NOTE: this function does not check the content of the GroupsResult, so
// the return could have duplicated groups
func mergeGroupsResult(grs ...*model.GroupsResult) (merged model.GroupsResult) {
	groups := make([]*model.Group, 0)

	for _, gr := range grs {
		groups = append(groups, gr.Resources...)
	}

	merged = model.GroupsResult{
		Items:     len(groups),
		Resources: groups,
	}
	merged.SetHashCode()

	return
}

// mergeUsersResult merges n UsersResult result
// NOTE: this function does not check the content of the UsersResult, so
// the return could have duplicated users
func mergeUsersResult(urs ...*model.UsersResult) (merged model.UsersResult) {
	users := make([]*model.User, 0)

	for _, u := range urs {
		users = append(users, u.Resources...)
	}

	merged = model.UsersResult{
		Items:     len(users),
		Resources: users,
	}
	merged.SetHashCode()

	return
}

// mergeGroupsMembersResult merges n GroupMembers result
// NOTE: this function does not check the content of the GroupMembers, so
// the return could have duplicated groupsMembers
func mergeGroupsMembersResult(gms ...*model.GroupsMembersResult) (merged model.GroupsMembersResult) {
	groupsMembers := make([]*model.GroupMembers, 0)

	for _, gm := range gms {
		groupsMembers = append(groupsMembers, gm.Resources...)
	}

	merged = model.GroupsMembersResult{
		Items:     len(groupsMembers),
		Resources: groupsMembers,
	}
	merged.SetHashCode()

	return
}

// updateSCIMID updates the SCIMID of the group in the idp object
// this is necessary because during the sync process we can create users and groups and to add
// these users to the groups we need to have the SCIMID of the user and the group
func updateSCIMID(idp *model.GroupsMembersResult, scimGroups *model.GroupsResult, scimUsers *model.UsersResult) *model.GroupsMembersResult {
	groups := make(map[string]model.Group)
	users := make(map[string]model.User)

	for _, group := range scimGroups.Resources {
		groups[group.Name] = *group
	}

	for _, user := range scimUsers.Resources {
		users[user.Email] = *user
	}

	gms := make([]*model.GroupMembers, 0)
	for _, groupMembers := range idp.Resources {
		mbs := make([]*model.Member, 0)

		g := model.Group{
			IPID:   groupMembers.Group.IPID,
			SCIMID: groups[groupMembers.Group.Name].SCIMID,
			Name:   groupMembers.Group.Name,
			Email:  groupMembers.Group.Email,
		}
		g.SetHashCode()

		for _, member := range groupMembers.Resources {
			m := &model.Member{
				IPID:   member.IPID,
				SCIMID: users[member.Email].SCIMID,
				Email:  member.Email,
			}
			m.SetHashCode()
			mbs = append(mbs, m)
		}

		gm := &model.GroupMembers{
			Items:     len(mbs),
			Group:     g,
			Resources: mbs,
		}
		gm.SetHashCode()

		gms = append(gms, gm)
	}

	gmr := &model.GroupsMembersResult{
		Items:     idp.Items,
		Resources: gms,
	}
	gmr.SetHashCode()

	return gmr
}
