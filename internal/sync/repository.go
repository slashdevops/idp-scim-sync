package sync

type SyncRepository interface {
	StoreGroups(gr *GroupsResult) (StoreGroupsResult, error)
	StoreUsers(ur *UsersResult) (StoreUsersResult, error)
	StoreGroupsMembers(gr *GroupsMembersResult) (StoreGroupsMembersResult, error)
	StoreState(state *SyncState) (StoreStateResult, error)
}

type StoreGroupsResult struct {
	Place string
}

type StoreUsersResult struct {
	Place string
}

type StoreGroupsMembersResult struct {
	Place string
}

type StoreStateResult struct {
	Place string
}
