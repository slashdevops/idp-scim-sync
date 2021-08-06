package sync

type SyncServiceOption func(*syncService)

func WithIdentityProviderGroupsFilter(filter []string) SyncServiceOption {
	return func(ss *syncService) {
		ss.provGroupsFilter = filter
	}
}

func WithIdentityProviderUsersFilter(filter []string) SyncServiceOption {
	return func(ss *syncService) {
		ss.provUsersFilter = filter
	}
}
