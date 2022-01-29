package google

type getGroupMembersOptions struct {
	includeDerivedMembership bool
	maxResults               int64
	pageToken                string
	roles                    string
}

type GetGroupMembersOption func(*getGroupMembersOptions)

// WithIncludeDerivedMembership is a GetGroupMembersOption that can be used to provide a filter for members by include derived membership.
// includeDerivedMembership Whether to list indirect memberships. Default: false.
func WithIncludeDerivedMembership(include bool) GetGroupMembersOption {
	return func(ggmo *getGroupMembersOptions) {
		ggmo.includeDerivedMembership = include
	}
}

// WithMaxResults is a GetGroupMembersOption that can be used to provide a filter for members by max results.
// maxResults Maximum number of results to return. Max allowed value is 200.
func WithMaxResults(maxResults int64) GetGroupMembersOption {
	return func(ggmo *getGroupMembersOptions) {
		ggmo.maxResults = maxResults
	}
}

// WithPageToken is a GetGroupMembersOption that can be used to provide a filter for members by page token.
// pageToken to specify next page in the list.
func WithPageToken(pageToken string) GetGroupMembersOption {
	return func(ggmo *getGroupMembersOptions) {
		ggmo.pageToken = pageToken
	}
}

// WithRoles is a GetGroupMembersOption that can be used to provide a filter for members by roles.
// roles=one or more of OWNER,MANAGER,MEMBER separated by a comma
func WithRoles(role string) GetGroupMembersOption {
	return func(ggmo *getGroupMembersOptions) {
		ggmo.roles = role
	}
}
