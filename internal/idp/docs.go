// Package idp implements core.IdentityProviderService against Google Workspace.
//
// It wraps the Directory API client in pkg/google and converts its responses
// into the domain types in internal/model, applying the rules AWS requires of a
// SCIM user along the way.
//
// # Records can be rejected
//
// buildUser returns nil when a Google record cannot become a valid AWS SCIM
// user — no name, no given name, no family name, or no primary email — and logs
// the reason. Callers skip those records and continue; one unusable directory
// entry must not stop the rest of a sync.
//
// # Nested groups
//
// Group membership is fetched with includeDerivedMembership, so members of
// nested groups are included while the nested group itself is skipped. AWS IAM
// Identity Center has no concept of a group inside a group, so the hierarchy is
// flattened.
//
// # Concurrency
//
// User lookups and per-group member lookups fan out through errgroup with a
// bounded limit, sharing one context so the first failure cancels the rest
// rather than spending Google API quota on results that will be discarded.
//
// See docs/Architecture.md for the full data flow.
package idp
