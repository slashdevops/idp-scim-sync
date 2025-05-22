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
		return nil, nil, nil, fmt.Errorf("error operating with groups: %w", err)
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
	slog.Info("reconciling the state data with the Identity Provider data")

	lastSyncTime, err := time.Parse(time.RFC3339, state.LastSync)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error parsing last sync time: %w", err)
	}

	slog.Info("syncing from state",
		"lastsync", state.LastSync,
		"since", time.Since(lastSyncTime).String(),
	)

	totalGroupsResult, err := syncGroupsFromState(ctx, scim, idpGroupsResult, state.Resources.Groups)
	if err != nil {
		return nil, nil, nil, err
	}

	totalUsersResult, err := syncUsersFromState(ctx, scim, idpUsersResult, state.Resources.Users)
	if err != nil {
		return nil, nil, nil, err
	}

	totalGroupsMembersResult, err := syncGroupsMembersFromState(ctx, scim, idpGroupsMembersResult, state.Resources.GroupsMembers, totalGroupsResult, totalUsersResult)
	if err != nil {
		return nil, nil, nil, err
	}

	return totalGroupsResult, totalUsersResult, totalGroupsMembersResult, nil
}

func syncGroupsFromState(
	ctx context.Context,
	scim SCIMService,
	idpGroupsResult *model.GroupsResult,
	stateGroupsResult *model.GroupsResult,
) (*model.GroupsResult, error) {
	if idpGroupsResult.HashCode == stateGroupsResult.HashCode {
		slog.Info("provider groups and state groups are the same, nothing to do with groups")
		return stateGroupsResult, nil
	}

	slog.Warn("provider groups and state groups are different")
	slog.Info("reconciling groups",
		"idp", idpGroupsResult.Items,
		"state", stateGroupsResult.Items,
	)
	groupsCreate, groupsUpdate, groupsEqual, groupsDelete, err := model.GroupsOperations(idpGroupsResult, stateGroupsResult)
	if err != nil {
		return nil, fmt.Errorf("error reconciling groups: %w", err)
	}

	groupsCreated, groupsUpdated, err := reconcilingGroups(ctx, scim, groupsCreate, groupsUpdate, groupsDelete)
	if err != nil {
		return nil, fmt.Errorf("error reconciling groups: %w", err)
	}

	return model.MergeGroupsResult(groupsCreated, groupsUpdated, groupsEqual), nil
}

func syncUsersFromState(
	ctx context.Context,
	scim SCIMService,
	idpUsersResult *model.UsersResult,
	stateUsersResult *model.UsersResult,
) (*model.UsersResult, error) {
	if idpUsersResult.HashCode == stateUsersResult.HashCode {
		slog.Info("provider users and state users are the same, nothing to do with users")
		return stateUsersResult, nil
	}

	slog.Warn("provider users and state users are different")
	slog.Info("reconciling users",
		"idp", idpUsersResult.Items,
		"state", stateUsersResult.Items,
	)
	usersCreate, usersUpdate, usersEqual, usersDelete, err := model.UsersOperations(idpUsersResult, stateUsersResult)
	if err != nil {
		return nil, fmt.Errorf("error operating with users: %w", err)
	}

	usersCreated, usersUpdated, err := reconcilingUsers(ctx, scim, usersCreate, usersUpdate, usersDelete)
	if err != nil {
		return nil, fmt.Errorf("error reconciling users: %w", err)
	}

	return model.MergeUsersResult(usersCreated, usersUpdated, usersEqual), nil
}

func syncGroupsMembersFromState(
	ctx context.Context,
	scim SCIMService,
	idpGroupsMembersResult *model.GroupsMembersResult,
	stateGroupsMembersResult *model.GroupsMembersResult,
	totalGroupsResult *model.GroupsResult,
	totalUsersResult *model.UsersResult,
) (*model.GroupsMembersResult, error) {
	if idpGroupsMembersResult.HashCode == stateGroupsMembersResult.HashCode {
		slog.Info("provider groups-members and state groups-members are the same, nothing to do with groups-members")
		return stateGroupsMembersResult, nil
	}

	slog.Warn("provider groups-members and state groups-members are different")

	groupsMembers := model.UpdateGroupsMembersSCIMID(idpGroupsMembersResult, totalGroupsResult, totalUsersResult)

	slog.Info("reconciling groups members",
		"idp", idpGroupsMembersResult.Items,
		"state", stateGroupsMembersResult.Items,
	)

	membersCreate, membersEqual, membersDelete, err := model.MembersOperations(groupsMembers, stateGroupsMembersResult)
	if err != nil {
		return nil, fmt.Errorf("error reconciling groups members: %w", err)
	}

	membersCreated, err := reconcilingGroupsMembers(ctx, scim, membersCreate, membersDelete)
	if err != nil {
		return nil, fmt.Errorf("error reconciling groups members: %w", err)
	}

	return model.MergeGroupsMembersResult(membersCreated, membersEqual), nil
}
