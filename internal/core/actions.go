package core

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/slashdevops/idp-scim-sync/internal/model"
)

// scimSync executes the sync of the data on the SCIM side and
// returns the datasets synced
func scimSync(
	ctx context.Context,
	scim SCIMService,
	idpGroupsResult *model.GroupsResult,
	idpUsersResult *model.UsersResult,
	idpGroupsMembersResult *model.GroupsMembersResult,
) (*model.GroupsResult, *model.UsersResult, *model.GroupsMembersResult, error) {
	slog.Warn("reconciling the SCIM data with the Identity Provider data")

	var totalGroupsResult *model.GroupsResult
	var totalUsersResult *model.UsersResult
	var totalGroupsMembersResult *model.GroupsMembersResult

	slog.Info("getting SCIM Groups")
	scimGroupsResult, err := scim.GetGroups(ctx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error getting groups from the SCIM service: %w", err)
	}

	slog.Info("reconciling groups",
		"idp", idpGroupsResult.Items,
		"scim", scimGroupsResult.Items,
	)

	groupsCreate, groupsUpdate, groupsEqual, groupsDelete, err := model.GroupsOperations(idpGroupsResult, scimGroupsResult)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error reconciling groups: %w", err)
	}

	groupsCreated, groupsUpdated, err := reconcilingGroups(ctx, scim, groupsCreate, groupsUpdate, groupsDelete)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error reconciling groups: %w", err)
	}

	// groupsCreated + groupsUpdated + groupsEqual = groups total
	totalGroupsResult = model.MergeGroupsResult(groupsCreated, groupsUpdated, groupsEqual)

	slog.Info("getting SCIM Users")
	scimUsersResult, err := scim.GetUsers(ctx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error getting users from the SCIM service: %w", err)
	}

	slog.Info("reconciling users",
		"idp", idpUsersResult.Items,
		"scim", scimUsersResult.Items,
	)
	usersCreate, usersUpdate, usersEqual, usersDelete, err := model.UsersOperations(idpUsersResult, scimUsersResult)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error operating with users: %w", err)
	}

	usersCreated, usersUpdated, err := reconcilingUsers(ctx, scim, usersCreate, usersUpdate, usersDelete)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error reconciling users: %w", err)
	}

	// usersCreated + usersUpdated + usersEqual = users total
	totalUsersResult = model.MergeUsersResult(usersCreated, usersUpdated, usersEqual)

	slog.Info("getting SCIM Groups Members")
	// unfortunately, the SCIM service does not support the getGroupsMembers method in and efficient way
	// see: "Nor Supported" section in: https://docs.aws.amazon.com/singlesignon/latest/developerguide/listgroups.html
	// scimGroupsMembersResult, err := scim.GetGroupsMembers(ctx, &totalGroupsResult) // not supported yet
	scimGroupsMembersResult, err := scim.GetGroupsMembersBruteForce(ctx, totalGroupsResult, totalUsersResult)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error getting groups members from the SCIM service: %w", err)
	}

	slog.Info("reconciling groups members",
		"idp", idpGroupsMembersResult.Items,
		"scim", scimGroupsMembersResult.Items,
	)
	membersCreate, membersEqual, membersDelete, err := model.MembersOperations(idpGroupsMembersResult, scimGroupsMembersResult)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error reconciling groups members: %w", err)
	}

	membersCreated, err := reconcilingGroupsMembers(ctx, scim, membersCreate, membersDelete)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error reconciling groups members: %w", err)
	}

	// membersCreate + membersEqual = members total
	totalGroupsMembersResult = model.MergeGroupsMembersResult(membersCreated, membersEqual)

	return totalGroupsResult, totalUsersResult, totalGroupsMembersResult, nil
}

// stateSync executes the sync of the data on the state side and
// returns the datasets synced
func stateSync(
	ctx context.Context,
	state *model.State,
	scim SCIMService,
	idpGroupsResult *model.GroupsResult,
	idpUsersResult *model.UsersResult,
	idpGroupsMembersResult *model.GroupsMembersResult,
) (*model.GroupsResult, *model.UsersResult, *model.GroupsMembersResult, error) {
	var totalGroupsResult *model.GroupsResult
	var totalUsersResult *model.UsersResult
	var totalGroupsMembersResult *model.GroupsMembersResult
	slog.Info("reconciling the state data with the Identity Provider data")

	lastSyncTime, err := time.Parse(time.RFC3339, state.LastSync)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error parsing last sync time: %w", err)
	}

	slog.Info("syncing from state",
		"lastsync", state.LastSync,
		"since", time.Since(lastSyncTime).String(),
	)

	if idpGroupsResult.HashCode == state.Resources.Groups.HashCode {
		slog.Info("provider groups and state groups are the same, nothing to do with groups")

		totalGroupsResult = state.Resources.Groups
	} else {
		slog.Warn("provider groups and state groups are different")
		// now here we have the google fresh data and the last sync data state
		// we need to compare the data and decide what to do
		// see differences between the two datasets
		slog.Info("reconciling groups",
			"idp", idpGroupsResult.Items,
			"state", state.Resources.Groups.Items,
		)
		groupsCreate, groupsUpdate, groupsEqual, groupsDelete, err := model.GroupsOperations(idpGroupsResult, state.Resources.Groups)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("error reconciling groups: %w", err)
		}

		groupsCreated, groupsUpdated, err := reconcilingGroups(ctx, scim, groupsCreate, groupsUpdate, groupsDelete)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("error reconciling groups: %w", err)
		}

		// merge in only one data structure the groups created, updated amd equals who has the SCIMID
		totalGroupsResult = model.MergeGroupsResult(groupsCreated, groupsUpdated, groupsEqual)
	}

	if idpUsersResult.HashCode == state.Resources.Users.HashCode {
		slog.Info("provider users and state users are the same, nothing to do with users")

		totalUsersResult = state.Resources.Users
	} else {
		slog.Warn("provider users and state users are different")

		slog.Info("reconciling users",
			"idp", idpUsersResult.Items,
			"state", state.Resources.Users.Items,
		)
		usersCreate, usersUpdate, usersEqual, usersDelete, err := model.UsersOperations(idpUsersResult, state.Resources.Users)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("error operating with users: %w", err)
		}

		usersCreated, usersUpdated, err := reconcilingUsers(ctx, scim, usersCreate, usersUpdate, usersDelete)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("error reconciling users: %w", err)
		}

		// usersCreated + usersUpdated + usersEqual = users total
		totalUsersResult = model.MergeUsersResult(usersCreated, usersUpdated, usersEqual)
	}

	if idpGroupsMembersResult.HashCode == state.Resources.GroupsMembers.HashCode {
		slog.Info("provider groups-members and state groups-members are the same, nothing to do with groups-members")

		totalGroupsMembersResult = state.Resources.GroupsMembers
	} else {
		slog.Warn("provider groups-members and state groups-members are different")

		// if we create a group or user during the sync, we need the scimid of these new groups/users
		// because to add members to a group the scim api needs that.
		// so this function will fill the scimid of the new groups/users
		groupsMembers := model.UpdateGroupsMembersSCIMID(idpGroupsMembersResult, totalGroupsResult, totalUsersResult)

		slog.Info("reconciling groups members",
			"idp", idpGroupsMembersResult.Items,
			"state", state.Resources.GroupsMembers.Items,
		)

		membersCreate, _, membersDelete, err := model.MembersOperations(groupsMembers, state.Resources.GroupsMembers)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("error reconciling groups members: %w", err)
		}

		_, err = reconcilingGroupsMembers(ctx, scim, membersCreate, membersDelete)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("error reconciling groups members: %w", err)
		}

		totalGroupsMembersResult = model.MergeGroupsMembersResult(groupsMembers)
	}
	return totalGroupsResult, totalUsersResult, totalGroupsMembersResult, nil
}
