// Package google is a thin, testable client for the Google Workspace Admin SDK
// Directory API.
//
// [NewService] builds an authenticated *admin.Service from service-account
// credentials using domain-wide delegation, impersonating the Workspace admin
// named in [DirectoryServiceConfig.UserEmail]. [NewDirectoryService] then wraps
// it in a [DirectoryService] offering only the calls this project needs.
//
// # Partial responses
//
// Every call requests an explicit field mask rather than whole objects, which
// keeps responses small on directories with many users. [WithSyncFieldSet]
// narrows the user mask further to just the attributes a given sync is
// configured to carry, so disabling an optional attribute also stops fetching
// it.
//
// # Membership
//
// [DirectoryService.ListGroupMembers] returns only ACTIVE members and can
// include derived (nested) membership. [DirectoryService.ListGroupMembersBatch]
// fans that out across many groups with a bounded concurrency limit, cancelling
// the remaining requests as soon as one fails.
//
// Usage:
//
//	svc, err := google.NewService(ctx, google.DirectoryServiceConfig{
//	    UserEmail:      "admin@example.com",
//	    ServiceAccount: credentialsJSON,
//	    Scopes:         scopes,
//	    UserAgent:      "idp-scim-sync/v1.0.0",
//	    Client:         httpClient,
//	})
//	if err != nil {
//	    return err
//	}
//	ds, err := google.NewDirectoryService(svc)
package google
