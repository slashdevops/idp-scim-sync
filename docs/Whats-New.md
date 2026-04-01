# What's New

This document tracks notable changes, new features, and bug fixes across releases.

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
