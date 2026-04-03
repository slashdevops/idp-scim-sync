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

* [README.md](../README.md)
* [Configuration.md](Configuration.md)
* [AWS-SAM.md](AWS-SAM.md)
* [Demo.md](Demo.md)
