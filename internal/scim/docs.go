// Package scim implements core.SCIMService against the AWS IAM Identity Center
// SCIM API.
//
// It is the adapter between the domain types in internal/model and the wire
// types in pkg/aws: it translates in both directions, batches group membership
// patches to stay inside the API's per-request limit, and works around the
// parts of the AWS SCIM implementation that do not behave like the
// specification.
//
// # Membership has to be inferred
//
// AWS does not populate the members array on either ListGroups or GetGroup
// responses. [Provider.GetGroupsMembers] therefore inverts the question: for
// each user it asks which groups contain them, using the
// `members.value eq "<id>"` filter with cursor-based pagination, and rebuilds
// the group-to-members mapping from the answers. That is one request per user
// rather than one per (group, user) pair.
//
// # Conflicts
//
// Creating a user or group that already exists returns 409. pkg/aws handles
// that by fetching the existing record instead, so a create is effectively
// upsert-like; see pkg/aws.SCIMService.CreateOrGetUser.
//
// See docs/Architecture.md for diagrams of both flows.
package scim
