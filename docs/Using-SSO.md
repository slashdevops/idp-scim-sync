# Using AWS IAM Identity Center

This document is a practical rollout guide for using `idp-scim-sync` with [AWS IAM Identity Center](https://aws.amazon.com/iam/identity-center/) and Google Workspace.

It is intentionally opinionated in one area: stable group design matters more than short-term convenience.

## Before You Deploy

Before enabling automatic synchronization, make sure you have all of the following:

* AWS IAM Identity Center enabled in the target AWS account or organization
* SCIM automatic provisioning enabled so you have a SCIM endpoint and access token
* A Google Workspace service account with domain-wide delegation
* Agreement on which Google Workspace groups should be synchronized

For the actual deployment mechanics, see [AWS-SAM.md](AWS-SAM.md). For command-level validation, see [idpscimcli.md](idpscimcli.md).

## Recommended Operating Model

The project works best when you synchronize groups and let users be derived from those groups.

Recommended practices:

* Use groups and group membership as the source of truth
* Avoid managing individual users separately when group membership already defines access
* Keep the synchronized scope narrow with explicit Google Workspace group filters
* Establish naming conventions before the first production sync
* Prefer predictable group names over ad hoc per-team naming

## Why Group Naming Matters

After the first successful sync, your AWS IAM Identity Center groups will typically receive permission set assignments. At that point, renaming groups becomes operationally expensive.

For this project, the most important stable identifiers are:

| Entity | Attribute you should keep stable |
| --- | --- |
| Groups | Display name |
| Users | Email address |

If a group name changes in the identity provider, AWS IAM Identity Center can no longer match it to the previous synchronized object in the way you expect for access governance. That can translate into lost permission set associations and manual cleanup.

The safest rule is simple:

* Finalize your group naming convention before the first production sync

## Start With A Filter Strategy

Filtering at the source is usually safer than synchronizing everything and trying to narrow the result later.

Example approach:

| Group Name | Group Email |
| --- | --- |
| AWS Administrators | `aws-administrators@example.com` |
| AWS DevOps | `aws-devops@example.com` |
| AWS Developers | `aws-developers@example.com` |

With a naming pattern like that, `idpscim` can use filters such as:

```bash
--gws-groups-filter 'name:AWS*'
```

That gives you a stable boundary for future growth. You can add or remove members and even create new `AWS...` groups without needing to rework the deployment parameters every time.

## Rollout Sequence

Recommended rollout order:

1. Plan the target group naming convention.
2. Create or clean up the Google Workspace groups you want to synchronize.
3. Validate filters and credentials with `idpscimcli`.
4. Deploy the application in a non-production AWS account first.
5. Verify synchronized groups, users, and permission assignments.
6. Move to production once the naming and filter strategy is stable.

## Validation Commands

Before enabling scheduled synchronization, run explicit validation commands.

Check Google Workspace groups:

```bash
./build/idpscimcli gws groups list \
  --gws-service-account-file credentials.json \
  --gws-user-email admin@example.com \
  --gws-groups-filter 'name=AWS*'
```

Check AWS IAM Identity Center SCIM connectivity:

```bash
./build/idpscimcli aws service config \
  --aws-scim-endpoint https://example.awsapps.com/scim/v2/ \
  --aws-scim-access-token "$SCIM_ACCESS_TOKEN"
```

## Related Documentation

* [Configuration.md](Configuration.md)
* [AWS-SAM.md](AWS-SAM.md)
* [idpscim.md](idpscim.md)
* [idpscimcli.md](idpscimcli.md)
