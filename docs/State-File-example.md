# State File Example

This document shows an example of the sync state file stored in S3.

The state file is an implementation detail used to make repeated synchronizations faster and to avoid unnecessary SCIM updates. It is useful for understanding how the project works, but it should not be treated as a stable external contract.

## What The State File Contains

The JSON representation contains:

* `groups`: synchronized group data
* `users`: synchronized user data
* `groupsMembers`: group membership data used for reconciliation
* `schemaVersion`: state schema version
* `codeVersion`: application version that produced the state
* `lastSync`: timestamp of the last successful synchronization
* `hashCode`: top-level hash used to detect changes efficiently

The current schema version in the codebase is `1.0.0`.

## How Change Detection Works

Every comparison the sync makes is a **hash comparison**, and understanding which fields participate
explains most of the system's behaviour.

Each `groups`, `users` and `groupsMembers` container carries its own `hashCode`. On each run the
identity provider's hash for a container is compared against the state's hash for the same container
([`internal/core/actions.go`](../internal/core/actions.go)). If they match, that entire resource kind
is skipped — no SCIM API calls at all. The three comparisons are independent, so users can be skipped
while groups are reconciled.

### What is hashed

Hashes are SHA-256 over a **`gob`** encoding, not over the JSON you see in the file. Each type in
`internal/model` hand-writes a `MarshalBinary` method that decides exactly which fields participate.

| Field | In the hash? | Why |
| --- | --- | --- |
| `ipid` | ✅ yes | The identity provider's own identifier |
| `name`, `email`, `userName`, `displayName`, … | ✅ yes | Identity-provider-owned attributes |
| `emails`, `addresses`, `phoneNumbers`, `enterpriseData` | ✅ yes | Synced attributes |
| `scimid` | ❌ **no** | Assigned by AWS, not the identity provider. Including it would make every first-run comparison differ, so nothing would ever be seen as unchanged. |
| `hashCode` | ❌ no | It is the output |
| `lastSync`, `codeVersion` | ❌ no | Metadata; a new run must not look like a data change |

### Ordering

Container hashes are computed over a **sorted copy** of their resources, so the hash does not depend
on the order elements arrived in. This matters because the identity-provider fan-outs build their
slices from maps, whose iteration order is randomised. Sorting uses `slices.SortStableFunc`, so
resources sharing a sort key keep a fixed relative order.

> [!NOTE]
> The **top-level** `hashCode` is written but never read back — only the three container hashes drive
> decisions. It is currently order-dependent, so it can differ between two runs over identical data.
> Don't use it as an external change indicator.

### Why this file is not a stable contract

The `MarshalBinary` methods are effectively a wire format. Changing field order, adding a hashed
field, or removing one changes **every** hash, which invalidates every deployed state file and causes
a one-off full re-sync. [`golden_hash_test.go`](../internal/model/golden_hash_test.go) pins the
current values so this cannot happen by accident, and releases that do change them say so in
[Whats-New.md](Whats-New.md).

## Deleting The State File Is Safe

The state file is a **cache, not a ledger**. If it is missing, empty or corrupt, delete it and let the
next run rebuild it — both `NoSuchKey` and an empty object are treated as "first run", not as errors.

On that first run the sync reconciles against the live SCIM provider instead of the state, and
**adopts** what is already there rather than recreating it: groups are matched by name, users by
primary email address, and only `externalId` is updated. That is the standard recovery procedure for
state corruption, and it is also why migrating from a different sync tool does not delete and
recreate every user.

The cost is one slower run that enumerates the SCIM population.

## Storage Location

The exact object key depends on how you deploy the application:

* The code default is `state.json`
* The AWS SAM template default is `data/state.json`

In both cases, the object is stored in the configured S3 state bucket.

## Example

```json
{
  "resources": {
    "groups": {
      "items": 1,
      "hashCode": "15cf5de941f6eb2d96e037675ac6f85401911889e12651f58990573c9f1f84ba",
      "resources": [
        {
          "ipid": "00examplegroup",
          "scimid": "b295b414-e091-70f6-3981-df556957e68a",
          "name": "AWS-Administrators",
          "email": "aws-administrators@example.com",
          "hashCode": "bcc54ec742946488860ec5f11eac4c958a178393a837abc878749fc0c40fefea"
        }
      ]
    },
    "users": {
      "items": 1,
      "hashCode": "bbbcf7f0ba3e94c811c03962ff986dcceffd97b1c95b0f6a50304df4d182380c",
      "resources": [
        {
          "ipid": "100000000000000000001",
          "scimid": "2275b4a4-d031-70b1-1bb0-e5049d0a0689",
          "userName": "alice@example.com",
          "displayName": "Alice Example",
          "title": "Platform Engineer",
          "userType": "admin#directory#user",
          "preferredLanguage": "en-US",
          "emails": [
            {
              "value": "alice@example.com",
              "primary": true
            }
          ],
          "addresses": [
            {
              "formatted": "123 Example Street"
            }
          ],
          "phoneNumbers": [
            {
              "value": "+1 555 0100",
              "type": "work"
            }
          ],
          "name": {
            "formatted": "Alice Example",
            "familyName": "Example",
            "givenName": "Alice"
          },
          "enterpriseData": {
            "costCenter": "ENG-001",
            "department": "Engineering"
          },
          "active": true,
          "hashCode": "4945a50f8b93337f5632dca20b49870f4507f0da28ee5d6d66add1f4b6df9045"
        }
      ]
    },
    "groupsMembers": {
      "items": 1,
      "hashCode": "72b7104a684c9cc04b04835c6f6e31deee272418440b3fd47c40a303c1fa3a02",
      "resources": [
        {
          "items": 1,
          "hashCode": "2b691179255bef46299eb3359433b5d019c6623904b90bf6fd032f4856ff7ded",
          "group": {
            "ipid": "00examplegroup",
            "scimid": "b295b414-e091-70f6-3981-df556957e68a",
            "name": "AWS-Administrators",
            "email": "aws-administrators@example.com",
            "hashCode": "bcc54ec742946488860ec5f11eac4c958a178393a837abc878749fc0c40fefea"
          },
          "resources": [
            {
              "ipid": "100000000000000000001",
              "scimid": "2275b4a4-d031-70b1-1bb0-e5049d0a0689",
              "email": "alice@example.com",
              "status": "ACTIVE",
              "hashCode": "f78efeb7e034db070cf78c804174f8de32a6a823d80674bae4d012f0fbecaf1f"
            }
          ]
        }
      ]
    }
  },
  "schemaVersion": "1.0.0",
  "codeVersion": "v0.44.0",
  "lastSync": "2026-04-03T10:15:00Z",
  "hashCode": "e72d58ac523af315fa6f3ed3329b8a174f2938c9e67a573ed45217f4a1a7b4e2"
}
```

## Related Documentation

* [Architecture.md](Architecture.md) — the sync algorithm and the hashing scheme in full
* [Implementation-Guide.md](Implementation-Guide.md) — state-file compatibility rules for contributors
* [Configuration.md](Configuration.md) — `aws_s3_bucket_name` and `aws_s3_bucket_key`
* [User-Manual.md](User-Manual.md) — inspecting and recovering the state file
* [AWS-SAM.md](AWS-SAM.md)
* [Demo.md](Demo.md)
* [README.md](../README.md)
