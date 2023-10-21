# State File example

This is an example of the `state file` the program will store in `AWS S3 Bucket`.

As you see this use `JSON` serialization to save the data and the data is about:

* Groups --> `groups` field
* Users --> `users` field
* Groups Members --> `groupsMembers` field

Also the `State file` contains some `metadata`:

* schemaVersion --> this could change if the `fields` of the `state file` change
* codeVersion --> this inform you about the version of the code that generated the `state file`
* lastSync --> this is the date and time when the `state file` was generated

and the `most important feature here` is the `hashCode` field, this is a `SHA256` hash of the each element of the `state file` content, and it is used to `save time in the operations` when we want to `detect changes`, also we can use that to checks `data integrity`.

```json
{
  "schemaVersion": "1.0.0",
  "codeVersion": "v0.1.0",
  "lastSync": "2023-10-21T18:48:49+02:00",
  "hashCode": "e72d58ac523af315fa6f3ed3329b8a174f2938c9e67a573ed45217f4a1a7b4e2",
  "resources": {
    "groups": {
      "items": 1,
      "hashCode": "15cf5de941f6eb2d96e037675ac6f85401911889e12651f58990573c9f1f84ba",
      "resources": [
        {
          "ipid": "00xvir7l2tu59gn",
          "scimid": "b295b414-e091-70f6-3981-df556957e68a",
          "name": "AWS-SSO-Administrators",
          "email": "aws-sso-administrators@slashdevops.com",
          "hashCode": "bcc54ec742946488860ec5f11eac4c958a178393a837abc878749fc0c40fefea"
        }
      ]
    },
    "users": {
      "items": 1,
      "hashCode": "bbbcf7f0ba3e94c811c03962ff986dcceffd97b1c95b0f6a50304df4d182380c",
      "resources": [
        {
          "hashCode": "4945a50f8b93337f5632dca20b49870f4507f0da28ee5d6d66add1f4b6df9045",
          "ipid": "100439965050892133351",
          "scimid": "2275b4a4-d031-70b1-1bb0-e5049d0a0689",
          "userName": "christian.gonzalez@slashdevops.com",
          "displayName": "Christian González Di Antonio",
          "title": "Chief Technology Officer",
          "userType": "admin#directory#user",
          "preferredLanguage": "en-GB",
          "emails": [
            {
              "value": "christian.gonzalez@slashdevops.com",
              "primary": true
            }
          ],
          "addresses": [
            {
              "formatted": "private address here",
            }
          ],
          "phoneNumbers": [
            {
              "value": "+55 555 555 555",
              "type": "work"
            }
          ],
          "name": {
            "formatted": "Christian González Di Antonio",
            "familyName": "González Di Antonio",
            "givenName": "Christian"
          },
          "enterpriseData": {
            "costCenter": "123654",
            "department": "IT"
          },
          "active": true
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
            "ipid": "00xvir7l2tu59gn",
            "scimid": "b295b414-e091-70f6-3981-df556957e68a",
            "name": "AWS-SSO-Administrators",
            "email": "aws-sso-administrators@slashdevops.com",
            "hashCode": "bcc54ec742946488860ec5f11eac4c958a178393a837abc878749fc0c40fefea"
          },
          "resources": [
            {
              "ipid": "100439965050892133351",
              "scimid": "2275b4a4-d031-70b1-1bb0-e5049d0a0689",
              "email": "christian.gonzalez@slashdevops.com",
              "status": "ACTIVE",
              "hashCode": "f78efeb7e034db070cf78c804174f8de32a6a823d80674bae4d012f0fbecaf1f"
            }
          ]
        }
      ]
    }
  }
}
```
