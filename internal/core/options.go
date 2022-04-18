package core

// SyncServiceOption is a function that can be used to configure the SyncService
// following the Option pattern.
type SyncServiceOption func(*SyncService)

// WithIdentityProviderGroupsFilter is a SyncServiceOption that can be used to
// provide a filter for the groups that should be synced.
func WithIdentityProviderGroupsFilter(filter []string) SyncServiceOption {
	return func(ss *SyncService) {
		ss.provGroupsFilter = filter
	}
}

// WithIdentityProviderUsersFilter is a SyncServiceOption that can be used to
// provide a filter for the users that should be synced.
func WithIdentityProviderUsersFilter(filter []string) SyncServiceOption {
	return func(ss *SyncService) {
		ss.provUsersFilter = filter
	}
}
