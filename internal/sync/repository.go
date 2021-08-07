package sync

type SyncRepository interface {
	StoreGroups(gr *GroupsResult) error
	StoreUsers(ur *UsersResult) error
	StoreGroupsMembers(gr *GroupsMembersResult) error
}
