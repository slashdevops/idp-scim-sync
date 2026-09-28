// Package core orchestrates the synchronization between an identity provider
// and a SCIM provider.
//
// It is the only package that knows the shape of a sync run, and it is
// deliberately free of concrete AWS and Google types: everything it touches
// arrives through three interfaces it declares itself —
// [IdentityProviderService], [SCIMService] and [StateRepository] — which the
// internal/idp, internal/scim and internal/repository packages implement.
//
// # The sync
//
// [SyncService.SyncGroupsAndTheirMembers] is the entry point. It reads groups,
// their members, and the users behind those members from the identity provider,
// then branches on whether a previous run left a state file:
//
//   - First run (state.LastSync is empty): reconcile directly against the SCIM
//     provider, so an existing SCIM population — for example one created by a
//     different tool — is adopted rather than recreated.
//   - Subsequent runs: reconcile against the state file, which is far cheaper
//     because the AWS SCIM API does not have to be enumerated.
//
// Either way the run ends by writing a fresh state file.
//
// # Change detection
//
// Every comparison is a hash comparison. Each model result type carries a
// HashCode covering only identity-provider-owned fields, so a run whose
// upstream data has not changed short-circuits without issuing a single write.
// See internal/model for how those hashes are computed and why they must stay
// stable.
//
// # Reconciliation
//
// model.GroupsOperations, model.UsersOperations and model.MembersOperations
// partition the two sides into create, update, equal and remove sets. Groups are
// keyed by name and users by primary email address, because those are the only
// attributes both sides agree on before AWS has assigned its own identifiers.
// The reconciling* functions then apply each set through [SCIMService].
//
// For the full picture, including diagrams, see docs/Architecture.md.
package core
