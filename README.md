# 🆔 idp-scim-sync

[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/5348/badge)](https://bestpractices.coreinfrastructure.org/projects/5348)
[![Build](https://github.com/slashdevops/idp-scim-sync/actions/workflows/build.yml/badge.svg)](https://github.com/slashdevops/idp-scim-sync/actions/workflows/build.yml)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/slashdevops/idp-scim-sync?style=plastic)
[![Go Report Card](https://goreportcard.com/badge/github.com/slashdevops/idp-scim-sync)](https://goreportcard.com/report/github.com/slashdevops/idp-scim-sync)
[![license](https://img.shields.io/github/license/slashdevops/idp-scim-sync.svg)](https://github.com/slashdevops/idp-scim-sync/blob/main/LICENSE)
[![Release](https://github.com/slashdevops/idp-scim-sync/actions/workflows/release.yml/badge.svg)](https://github.com/slashdevops/idp-scim-sync/actions/workflows/release.yml)
[![release](https://img.shields.io/github/release/slashdevops/idp-scim-sync/all.svg)](https://github.com/slashdevops/idp-scim-sync/releases)
[![codecov](https://codecov.io/gh/slashdevops/idp-scim-sync/branch/main/graph/badge.svg?token=H72NWJGHZ0)](https://codecov.io/gh/slashdevops/idp-scim-sync)

Keep your [AWS IAM Identity Center](https://aws.amazon.com/iam/identity-center/) (formerly AWS SSO) in sync with your [Google Workspace](https://workspace.google.com/) directory using an [AWS Lambda function](https://aws.amazon.com/lambda/). 🚀

![On AWS](https://raw.githubusercontent.com/slashdevops/idp-scim-sync/main/docs/images/diagrams/ipd-scim-sync.drawio.png)

## ✨ Features

* ✅ **Extended Attribute Support**: Syncs extended AWS SSO SCIM API fields as described in the [official documentation](https://docs.aws.amazon.com/singlesignon/latest/developerguide/limitations.html).
* ✅ **Configurable User Fields**: Choose which optional user attributes (phone numbers, addresses, enterprise data, etc.) to sync. See [Configurable User Fields](#configurable-user-fields) for details.
* ✅ **Efficient Data Retrieval**: Uses [partial responses](https://cloud.google.com/storage/docs/json_api#partial-response) from the Google Workspace API to fetch only the data you need.
* ✅ **Nested Groups Support**: Supports nested groups in Google Workspace thanks to the `includeDerivedMembership` API query parameter.
* ✅ **Multiple Deployment Options**: Can be deployed via the `AWS Serverless Application Repository`, as a `Container Image`, or as a `CLI`.
* ✅ **Incremental Sync**: Drastically reduces the number of requests to the AWS SSO SCIM API by using a [state file](docs/State-File-example.md) to track changes.

## Compatibility

This project is compatible with the latest AWS Lambda runtimes. Since version `v0.0.19`, it uses the `provided.al2` runtime and `arm64` architecture.

| Version Range        | AWS Lambda Runtime | Architecture       | Deprecation Date |
| -------------------- | ------------------ | ------------------ | ---------------- |
| `<= v0.0.18`         | Go 1.x             | amd64 (Intel)      | 2023-12-31       |
| `>= v0.0.19 < v0.31.0` | provided.al2       | arm64 (Graviton 2) | 2026-06-30       |
| `>= v0.31.0`         | provided.al2023    | arm64 (Graviton 2) | 2029-06-30       |

## ⚙️ How It Works

The AWS Lambda function is triggered by a CloudWatch event rule (every 15 minutes by default). It syncs your AWS IAM Identity Center with your Google Workspace directory using their respective APIs.

During the first sync, the data of your Groups and Users is stored in an AWS S3 bucket as a [state file](docs/State-File-example.md). This state file is a custom implementation to save time and requests to the AWS SSO SCIM API, and to mitigate some of its limitations.

This project is developed using the [Go language](https://go.dev/) and [AWS SAM](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-specification.html).

For more details on the resources created by the CloudFormation template, please check the [AWS SAM Template documentation](docs/AWS-SAM-Template.md).

> **Note:** If this is your first time implementing AWS IAM Identity Center, please read [Using SSO](docs/Using-SSO.md).

## 🚀 Getting Started

The easiest way to deploy and use this project is through the [AWS Serverless Application Repository](https://serverlessrepo.aws.amazon.com/applications/us-east-1/889836709304/idp-scim-sync).

### Credentials

You will need to configure credentials for both Google Workspace and AWS.

* **Google Workspace API Credentials**
  * Follow the [Google Workspace documentation](https://developers.google.com/workspace/guides/create-credentials) to create credentials.
  * You will need to create a **Service Account** and delegate **domain-wide authority** to it with the following scopes:
    * `https://www.googleapis.com/auth/admin.directory.group.readonly`
    * `https://www.googleapis.com/auth/admin.directory.user.readonly`
    * `https://www.googleapis.com/auth/admin.directory.group.member.readonly`

* **AWS SSO SCIM API Credentials**
  * Configure these credentials in the [AWS IAM Identity Center](https://aws.amazon.com/iam/identity-center/) service by following the [Automatic provisioning guide](https://docs.aws.amazon.com/singlesignon/latest/userguide/provision-automatically.html).

## 🛠️ Usage

You have several options to use this project:

### In AWS

* **AWS Serverless Application Repository (Recommended)**
  * Deploy the application directly from the [AWS Serverless Application Repository](https://serverlessrepo.aws.amazon.com/applications/us-east-1/889836709304/idp-scim-sync).

* **AWS SAM**
  * Build and deploy the Lambda function from your local machine.
  * **Requirements:**
    * [Git](https://git-scm.com/)
    * [Go](https://go.dev/learn/)
    * [AWS SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)
  * **Commands:**

```bash
# Set your AWS CLI profile and region
export AWS_PROFILE=<profile_name>
export AWS_REGION=<region>

# Validate the template
sam validate

# Build the project
sam build

# Deploy with a guided process
sam deploy --guided --capabilities CAPABILITY_IAM CAPABILITY_NAMED_IAM
```

### Locally

* **Build from Source**
  * **Requirements:**
    * [Git](https://git-scm.com/)
    * [Go](https://go.dev/learn/)
    * [Make](https://www.gnu.org/software/make/)
  * **Commands:**

```bash
# Compile for your operating system
make

# Cross-compile for Windows, macOS, and Linux
make build-dist
```

* **Pre-built Binaries**
  * Download the binaries from the [GitHub Releases](https://github.com/slashdevops/idp-scim-sync/releases).

* **Docker Image**
  * Pull the image from one of the public repositories.

## 📦 Repositories

* 📦 [AWS Serverless Application Repository](https://serverlessrepo.aws.amazon.com/applications/us-east-1/889836709304/idp-scim-sync)
* 📦 [AWS ECR Public Gallery](https://gallery.ecr.aws/l2n7y5s7/slashdevops/idp-scim-sync)
* 📦 [GitHub Packages](https://github.com/slashdevops/idp-scim-sync/pkgs/container/idp-scim-sync)
* 📦 [Docker Hub](https://hub.docker.com/r/slashdevops/idp-scim-sync)

## Configurable User Fields

By default, all optional user attributes are synced from Google Workspace to AWS SSO SCIM. You can control which optional fields are included using the `sync_user_fields` configuration option.

**Available fields:**

| Field               | Description                                                                 |
| ------------------- | --------------------------------------------------------------------------- |
| `phoneNumbers`      | User's phone numbers                                                        |
| `addresses`         | User's addresses (street, city, region, postal code, country)               |
| `title`             | User's job title                                                            |
| `preferredLanguage` | User's preferred language                                                   |
| `locale`            | User's locale                                                               |
| `timezone`          | User's timezone                                                             |
| `nickName`          | User's nickname                                                             |
| `profileURL`        | User's profile URL                                                          |
| `userType`          | User type attribute                                                         |
| `enterpriseData`    | Enterprise extension (employeeNumber, costCenter, organization, department, division, manager) |

**Required fields** (always synced, not configurable): `name`, `userName`, `displayName`, `emails`, `active`.

### Configuration Examples

**Config file (.idpscim.yaml):**

```yaml
# Sync only phone numbers, addresses, and enterprise data
sync_user_fields:
  - phoneNumbers
  - addresses
  - enterpriseData
```

**Environment variable (Lambda / SAM):**

```bash
IDPSCIM_SYNC_USER_FIELDS=phoneNumbers,addresses,enterpriseData
```

**CLI flag:**

```bash
idpscim --sync-user-fields phoneNumbers,addresses,enterpriseData
```

**SAM template parameter:**

Set the `SyncUserFields` parameter when deploying:

```bash
sam deploy --parameter-overrides SyncUserFields=phoneNumbers,addresses,enterpriseData
```

### Behavior Notes

* **Default (empty or not set):** When `sync_user_fields` is empty or not configured, all optional fields are synced. This preserves backward compatibility with existing deployments.
* **Specifying fields:** Only the listed fields will be synced. For example, setting `sync_user_fields: [phoneNumbers]` will sync only phone numbers; addresses, enterprise data, and other optional attributes will not be sent to AWS SSO SCIM.
* **Invalid field names:** If an invalid field name is provided, the application will fail at startup with a clear error message listing the unrecognized field.
* **Changing on an existing deployment:** The first sync after modifying this configuration will detect all users as "changed" (due to hash differences) and update them in AWS SSO. This is expected behavior — it will clear the excluded fields from SCIM.

## ⚠️ Limitations

* **Group Limit**: The AWS SSO SCIM API has a limit of 50 groups per request. Please support the feature request on the [AWS Support site](https://repost.aws/questions/QUqqnVkIo_SYyF_SlX5LcUjg/aws-sso-scim-api-pagination-for-methods) to help get this limit increased.
* **Throttling**: With a large number of users and groups, you may encounter a `ThrottlingException` from the AWS SSO SCIM API. This project uses a [retryable HTTP client](https://github.com/hashicorp/go-retryablehttp) to mitigate this, but it's still a possibility.
* **User Status**: The Google Workspace API doesn't differentiate between normal and guest users except for their status. This project only syncs `ACTIVE` users.

## For `ssosync` Users

If you are coming from the [awslabs/ssosync](https://github.com/awslabs/ssosync) project, please note the following:

* This project only implements the `--sync-method groups`.
* This project only implements filtering for Google Workspace Groups, not Users.
* The flag names are different.
* Not all features of `ssosync` are implemented here, and they may not be in the future.

## 🧩 Components

* **idpscim**: A program for keeping AWS IAM Identity Center groups and users synced with your Google Workspace directory. See the [idpscim documentation](docs/idpscim.md) for more details.
* **idpscimcli**: A command-line tool to check and validate some of the functionalities implemented in `idpscim`. See the [idpscimcli documentation](docs/idpscimcli.md) for more details.

## 📄 License

This project is released under the Apache License 2.0. See the [LICENSE](LICENSE) file for more details.
