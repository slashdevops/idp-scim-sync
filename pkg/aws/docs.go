// Package aws provides AWS SSO SCIM API client implementation.
// It supports all standard SCIM 2.0 operations for users and groups
// as defined in the AWS SSO SCIM API specification.
//
// Basic usage:
//
//	client, err := aws.NewSCIMService(httpClient, endpoint, token)
//	if err != nil {
//	    return err
//	}
//
//	users, err := client.ListUsers(ctx, "")
package aws
