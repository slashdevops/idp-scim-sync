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
  "codeVersion": "v0.0.1",
  "lastSync": "2022-01-15T17:55:36+01:00",
  "hashCode": "1be3807fa0a69ea22251f4ed71aee4c4c1bf0aa9a163a5c24f8ac4425e6f0d69",
  "resources": {
    "groups": {
      "items": 3,
      "hashCode": "251033ebdbcecb157b6d989c2847f1912c3460eab7a56eb4081f4d912e5145b1",
      "resources": [
        {
          "ipid": "00nmf14n0mfn3n3",
          "scimid": "90675b464e-5251235c-1f8c-4607-906e-015d9efc29a2",
          "name": "Administrators",
          "email": "administrators@<your domain here>",
          "hashCode": "f87938591c76e34cce90f79d86e14f9280a8bb2052c92d577f3705c3b681aefa"
        },
        {
          "ipid": "00ihv63633k64om",
          "scimid": "90675b464e-0880a5e4-6601-4138-91a8-b599aedf7a83",
          "name": "AWS Administrators",
          "email": "aws-administrators@<your domain here>",
          "hashCode": "eb3e4b4061c3781aac2ba3228b3c0d9a763909c326aeb7106460d72eb062657c"
        },
        {
          "ipid": "019c6y180i470k3",
          "scimid": "90675b464e-79914545-790b-4171-9142-36a55acf5a39",
          "name": "AWS DevOps",
          "email": "aws-devops@<your domain here>",
          "hashCode": "ca1462c1188f8e583fa0e79ae9f4a06651f188c4e3a6809387a29bbe243ab38f"
        }
      ]
    },
    "users": {
      "items": 1,
      "hashCode": "82da67ab6b9c0576b727fe4053c4d7c3ba0c3a7f0e88115920f76706411901f8",
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
          "email": "christian@<your domain here>",
          "hashCode": "7a0c55f4941d4512b9cdb8880e5e61130d089b96e7547d6df9ad2e38c2932f49"
        }
      ]
    },
    "groupsMembers": {
      "items": 3,
      "hashCode": "c7ed0e0b713310d4fb503e6cc8d4bde397766a9dab486e18df3a47d3881c237e",
      "resources": [
        {
          "items": 1,
          "hashCode": "8f2632a7e5c1fba4360066bdf69729f213f12037b825d25be29f682844751644",
          "group": {
            "ipid": "00nmf14n0mfn3n3",
            "scimid": "90675b464e-5251235c-1f8c-4607-906e-015d9efc29a2",
            "name": "Administrators",
            "email": "administrators@<your domain here>",
            "hashCode": "f87938591c76e34cce90f79d86e14f9280a8bb2052c92d577f3705c3b681aefa"
          },
          "resources": [
            {
              "ipid": "100439965050892133351",
              "scimid": "90675b464e-11025ca4-0a49-480e-afd5-5eda1ae3fc3c",
              "email": "christian@<your domain here>",
              "hashCode": "0563204b5acd6ce1f481e86b29ea4b4b5feab0cf84799f143c191bfe912ec571"
            }
          ]
        },
        {
          "items": 1,
          "hashCode": "80902e649cd990205eeed09e3e2d3714ea273c28f2ccd7af976f04606f2153c8",
          "group": {
            "ipid": "00ihv63633k64om",
            "scimid": "90675b464e-0880a5e4-6601-4138-91a8-b599aedf7a83",
            "name": "AWS Administrators",
            "email": "aws-administrators@<your domain here>",
            "hashCode": "eb3e4b4061c3781aac2ba3228b3c0d9a763909c326aeb7106460d72eb062657c"
          },
          "resources": [
            {
              "ipid": "100439965050892133351",
              "scimid": "90675b464e-11025ca4-0a49-480e-afd5-5eda1ae3fc3c",
              "email": "christian@<your domain here>",
              "hashCode": "0563204b5acd6ce1f481e86b29ea4b4b5feab0cf84799f143c191bfe912ec571"
            }
          ]
        },
        {
          "items": 0,
          "hashCode": "9e27fbd3e3ae6ba1a115fe95e01ae2009bdc0a9d953ba83139911903cc5e34d9",
          "group": {
            "ipid": "019c6y180i470k3",
            "scimid": "90675b464e-79914545-790b-4171-9142-36a55acf5a39",
            "name": "AWS DevOps",
            "email": "aws-devops@<your domain here>",
            "hashCode": "ca1462c1188f8e583fa0e79ae9f4a06651f188c4e3a6809387a29bbe243ab38f"
          },
          "resources": []
        }
      ]
    }
  }
}
```
