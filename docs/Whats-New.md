# What's New

This document tracks notable changes, new features, and bug fixes across releases.

## v0.40.1

### Improved HTTP Retry Library

Replaced the `httpretrier` library with [httpx](https://github.com/slashdevops/httpx), a zero-dependency HTTP client with built-in retry support.

**Why:** The previous library did not properly handle HTTP `429 Too Many Requests` responses, which caused issues with AWS SSO SCIM API throttling under high load.

**What changed:**

* The `httpx` library automatically retries on `429` and `5xx` responses with configurable backoff strategies.
* AWS SCIM API calls now use **jitter backoff** instead of simple exponential backoff, reducing the chance of thundering herd effects during rate limiting.
* Google Workspace API calls use **exponential backoff** for reliable retries.
* The `httpx` library has zero external dependencies and integrates with Go's `slog` logging.

### AWS SCIM Client Improvements (`pkg/aws`)

Several code quality improvements and bug fixes in the AWS SCIM client:

* **Bug fix:** `CreateOrGetUser` used `reflect.DeepEqual` to compare a `*CreateUserRequest` with a `*GetUserResponse` — different types, so the comparison always returned `false`, causing unnecessary PUT updates on every 409 conflict. Replaced with a typed `usersEqual` function that compares only sync-relevant attributes.
* **Removed `pkg/errors` dependency:** Replaced with stdlib `errors` and `fmt` packages. Sentinel errors now use `errors.New` instead of `errors.Errorf`.
* **Go 1.26 `errors.AsType`:** Migrated all `errors.As` calls to the generic `errors.AsType[T]` for compile-time type safety and better performance.
* **Fixed `String()` methods:** `User.String()` and `Group.String()` no longer call `os.Exit(1)` on marshal failure. They return a safe fallback string instead.
* **Eliminated double JSON decode:** `GetUserByUserName` and `GetGroupByDisplayName` no longer marshal a resource to JSON and re-parse it. They use direct type conversion instead.
* **Fixed decode error fallback:** `CreateGroup` and `CreateOrGetGroup` no longer attempt to read an already-consumed response body on decode failure.
* **Removed redundant context set:** `do()` no longer calls `req.WithContext(ctx)` since the request is already created with `http.NewRequestWithContext`.
* **Simplified type conversions:** `CreateOrGetUser` and `CreateOrGetGroup` use type conversions instead of manual field-by-field struct copies.

### Go 1.26 Modernization

Applied Go 1.26 best practices across the codebase:

* **Removed `github.com/pkg/errors` dependency:** Replaced all `errors.Wrap` and `errors.Errorf` with stdlib `fmt.Errorf` (with `%w`) and `errors.New` in `internal/setup`, `internal/repository`, and `pkg/aws`.
* **`errors.AsType[T]`:** Migrated `errors.As` calls to the generic `errors.AsType[T]` in `internal/core/sync.go` for type safety and performance.
* **Fixed `os.Exit` in `Hash()`:** `internal/model.Hash()` no longer calls `os.Exit(1)` on nil input or encoding failure. It panics instead (appropriate for programming errors, recoverable, produces stack trace).

## v0.44.0

### Configurable User Fields

You can now choose which optional user attributes are synced from Google Workspace to AWS SSO SCIM using the new `sync_user_fields` configuration option.

For example, sync only phone numbers and enterprise data while excluding addresses, locale, or timezone. When not configured, all fields are synced as before (fully backward compatible).

**Available fields:** `phoneNumbers`, `addresses`, `title`, `preferredLanguage`, `locale`, `timezone`, `nickName`, `profileURL`, `userType`, `enterpriseData`.

See [Configurable User Fields](../README.md#configurable-user-fields) for configuration examples and behavior notes.

### Bug Fix: Unnecessary member re-syncs

Fixed a bug where group members were re-synced on every Lambda execution even when nothing changed in Google Workspace.

**Root cause:** `MergeGroupsMembersResult` was not consolidating entries for the same group when merging "created" and "equal" member sets. This produced duplicate group entries in the state file, causing the groups-members hash to never match the IDP data on subsequent syncs.

**Impact:** After upgrading, the first sync will reconcile the state file automatically. Subsequent syncs will correctly skip member reconciliation when no changes are detected.

### Performance Improvements

* **Concurrent user fetching:** `GetUsersByGroupsMembers` now fetches user details from the Google Workspace API concurrently (up to 10 parallel requests) instead of sequentially. For deployments with 100+ users, this reduces the user retrieval phase from minutes to seconds.

* **Optimized member comparison:** Removed a redundant O(m) inner loop in `membersDataSets` that iterated over the entire SCIM member set to find an email already confirmed by a direct map lookup. Benchmarks show ~16-19% improvement for large groups.

* **Goroutine leak safety:** Concurrent operations are verified with `synctest.Test` (Go 1.26) to ensure no goroutine leaks in both success and error paths.
