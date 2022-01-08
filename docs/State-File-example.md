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

and the `most important feature here` is the `hashCode` field, this is a `SHA1` hash of the each element of the `state file` content, and it is used to `save time in the operations` when we want to `detect changes`, also we can use that to checks `data integrity`.

```json
{
  "schemaVersion": "1.0.0",
  "codeVersion": "0.0.1",
  "lastSync": "2022-01-08T18:35:16Z",
  "hashCode": "152814f8e5a3401080a03a91c3f1ec9d245be05b",
  "resources": {
    "groups": {
      "items": 2,
      "hashCode": "84f8a793162a18e5ea66c4f15f9d4ec1b6441e45",
      "resources": [
        {
          "ipid": "00ihv63633k64om",
          "scimid": "90675b464e-0880a5e4-6601-4138-91a8-b599aedf7a83",
          "name": "AWS Administrators",
          "email": "aws-administrators@<your email domain>",
          "hashCode": "aa7d1b4f76148dfc0ca6513ed1181dd1f73801ba"
        },
        {
          "ipid": "019c6y180i470k3",
          "scimid": "90675b464e-79914545-790b-4171-9142-36a55acf5a39",
          "name": "AWS DevOps",
          "email": "aws-devops@<your email domain>",
          "hashCode": "a0950ba2b70c751a8f5e34260d54b3b61734f990"
        }
      ]
    },
    "users": {
      "items": 3,
      "hashCode": "1aabe861b2d7f76d841e4cd17918062d9ecf74ef",
      "resources": [
        {
          "ipid": "100439965050892133351",
          "scimid": "90675b464e-11025ca4-0a49-480e-afd5-5eda1ae3fc3c",
          "name": {
            "familyName": "González Di Antonio",
            "givenName": "Christian"
          },
          "displayName": "Christian González Di Antonio",
          "active": true,
          "email": "administrator@<your email domain>",
          "hashCode": "76cd8b44628099358d2d6e6b616e306e1d04412a"
        },
        {
          "ipid": "113635714534969451687",
          "scimid": "90675b464e-803cfb36-4b44-47c2-b2bc-30dbb26c435d",
          "name": {
            "familyName": "user 2",
            "givenName": "test"
          },
          "displayName": "test user 2",
          "active": true,
          "email": "test.user2@<your email domain>",
          "hashCode": "af65f4d493bc415bf8dbd9db01fd45aa3895ccfb"
        },
        {
          "ipid": "106605753848140623644",
          "scimid": "90675b464e-f52d738f-4140-4c25-a904-7f70a4aa7e15",
          "name": {
            "familyName": "user 1",
            "givenName": "test"
          },
          "displayName": "test user 1",
          "active": true,
          "email": "test@<your email domain>",
          "hashCode": "4c274b0b85ffc0d6daba552dc03c5b2ff13113d6"
        }
      ]
    },
    "groupsMembers": {
      "items": 2,
      "hashCode": "4c6206f888592fa8c8fab9a08445b8cd5743cbab",
      "resources": [
        {
          "items": 2,
          "hashCode": "e46da18f2f3915727585de60bf131c44e068632f",
          "group": {
            "ipid": "00ihv63633k64om",
            "scimid": "90675b464e-0880a5e4-6601-4138-91a8-b599aedf7a83",
            "name": "AWS Administrators",
            "email": "aws-administrators@<your email domain>",
            "hashCode": "aa7d1b4f76148dfc0ca6513ed1181dd1f73801ba"
          },
          "resources": [
            {
              "ipid": "100439965050892133351",
              "scimid": "90675b464e-11025ca4-0a49-480e-afd5-5eda1ae3fc3c",
              "email": "administrator@<your email domain>",
              "hashCode": "abd437534e201e974335561543b5a2b8084cbc8a"
            },
            {
              "ipid": "106605753848140623644",
              "scimid": "90675b464e-f52d738f-4140-4c25-a904-7f70a4aa7e15",
              "email": "test@<your email domain>",
              "hashCode": "9c72d5480a4f5a9a849f43036508244f616aa5f6"
            }
          ]
        },
        {
          "items": 1,
          "hashCode": "066a9818c508574523245424a557471f8175ad4f",
          "group": {
            "ipid": "019c6y180i470k3",
            "scimid": "90675b464e-79914545-790b-4171-9142-36a55acf5a39",
            "name": "AWS DevOps",
            "email": "aws-devops@<your email domain>",
            "hashCode": "a0950ba2b70c751a8f5e34260d54b3b61734f990"
          },
          "resources": [
            {
              "ipid": "113635714534969451687",
              "scimid": "90675b464e-803cfb36-4b44-47c2-b2bc-30dbb26c435d",
              "email": "test.user2@<your email domain>",
              "hashCode": "50dcf8623cefa264ef484c8aa6cced324d783585"
            }
          ]
        }
      ]
    }
  }
}
```
