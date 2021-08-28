package core

type SyncRepository interface {
	StoreGroups(gr *GroupsResult) (StoreGroupsResult, error)
	StoreUsers(ur *UsersResult) (StoreUsersResult, error)
	StoreGroupsMembers(gr *GroupsMembersResult) (StoreGroupsMembersResult, error)
	StoreState(state *SyncState) (StoreStateResult, error)
	GetState() (SyncState, error)
	GetGroups(place string) (*GroupsResult, error)
	GetUsers(place string) (*UsersResult, error)
	GetGroupsMembers(place string) (*GroupsMembersResult, error)
}
