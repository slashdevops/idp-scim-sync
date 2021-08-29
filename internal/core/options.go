package core

type SyncServiceOption func(*SyncService)

func WithIdentityProviderGroupsFilter(filter []string) SyncServiceOption {
	return func(ss *SyncService) {
		ss.provGroupsFilter = filter
	}
}

func WithIdentityProviderUsersFilter(filter []string) SyncServiceOption {
	return func(ss *SyncService) {
		ss.provUsersFilter = filter
	}
}
