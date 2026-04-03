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

## 🆕 What's New

For a detailed list of new features, improvements, and bug fixes in each release, see the [What's New](docs/Whats-New.md) page.

## Compatibility

This project is compatible with the latest AWS Lambda runtimes. Since version `v0.0.19`, it uses the `provided.al2` runtime and `arm64` architecture.

| Version Range          | AWS Lambda Runtime | Architecture       | Deprecation Date |
| ---------------------- | ------------------ | ------------------ | ---------------- |
| `<= v0.0.18`           | Go 1.x             | amd64 (Intel)      | 2023-12-31       |
| `>= v0.0.19 < v0.31.0` | provided.al2       | arm64 (Graviton 2) | 2026-06-30       |
| `>= v0.31.0`           | provided.al2023    | arm64 (Graviton 2) | 2029-06-30       |

## ⚙️ How It Works

The AWS Lambda function is triggered by a CloudWatch event rule (every 15 minutes by default). It syncs your AWS IAM Identity Center with your Google Workspace directory using their respective APIs.

During the first sync, the data of your Groups and Users is stored in an AWS S3 bucket as a [state file](docs/State-File-example.md). This state file is a custom implementation to save time and requests to the AWS SSO SCIM API, and to mitigate some of its limitations.

This project is developed using the [Go language](https://go.dev/) and [AWS SAM](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-specification.html).

For more details on the resources created by the CloudFormation template, please check the [AWS SAM Template documentation](docs/AWS-SAM-Template.md).

> **Note:** If this is your first time implementing AWS IAM Identity Center, please read [Using SSO](docs/Using-SSO.md).

## Programs

This repository builds two binaries from the `cmd/` directory:

| Program | Source | Purpose |
| ------- | ------ | ------- |
| `idpscim` | `cmd/idpscim` | Main synchronization program that runs as the Lambda function, a local CLI, or a container command |
| `idpscimcli` | `cmd/idpscimcli` | Helper CLI used to inspect AWS SCIM and Google Workspace data while validating configuration |

After `make build`, the binaries are available in `build/`:

```bash
./build/idpscim --help
./build/idpscimcli --help
```

## Documentation Map

The repository documentation is organized as follows:

| Document | Purpose |
| ------- | ------- |
| [docs/idpscim.md](docs/idpscim.md) | Main program reference for the `idpscim` sync executable |
| [docs/idpscimcli.md](docs/idpscimcli.md) | Command reference for the `idpscimcli` validation and inspection CLI |
| [docs/Configuration.md](docs/Configuration.md) | Configuration sources, examples, and environment variable usage |
| [docs/AWS-SAM.md](docs/AWS-SAM.md) | Source deployment, Serverless Application Repository update flow, and maintainer publishing workflow |
| [docs/AWS-SAM-Template.md](docs/AWS-SAM-Template.md) | Template parameters, generated resources, and Lambda environment mapping |
| [docs/Development.md](docs/Development.md) | Local development workflow, build steps, tests, and SAM-based cloud testing |
| [docs/Using-SSO.md](docs/Using-SSO.md) | Practical rollout guidance for AWS IAM Identity Center and Google Workspace group design |
| [docs/State-File-example.md](docs/State-File-example.md) | Example state file structure and notes about how sync state is stored |
| [docs/Demo.md](docs/Demo.md) | Visual walkthrough screenshots of the sync process and resulting AWS and Google Workspace data |
| [docs/Release.md](docs/Release.md) | Maintainer release flow based on semantic version tags and GitHub Actions |
| [docs/Whats-New.md](docs/Whats-New.md) | Release notes and notable changes across versions |

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
  * To update an existing deployment to a newer published version, reuse the same original application name that you entered when you first deployed it. Do not use the generated `serverlessrepo-...` stack name. See [docs/AWS-SAM.md](docs/AWS-SAM.md) for the full update flow.

* **AWS SAM**
  * Build and deploy the Lambda function from your local machine.
  * Quick start:

```bash
export AWS_PROFILE=<profile_name>
export AWS_REGION=<region>
GIT_VERSION=dev sam build
sam deploy --guided --stack-name idp-scim-sync --capabilities CAPABILITY_IAM CAPABILITY_NAMED_IAM
```

* For full validation, source deployment, publish, and update guidance, see [docs/AWS-SAM.md](docs/AWS-SAM.md) and [docs/AWS-SAM-Template.md](docs/AWS-SAM-Template.md).

### Locally

* **Build from Source**
  * Quick start:

```bash
make
./build/idpscim --help
./build/idpscimcli --help
```

* **Run the programs**
  * **idpscim** runs the actual synchronization logic.
  * **idpscimcli** helps you validate your AWS SCIM and Google Workspace configuration before enabling automated sync.
  * See [docs/idpscim.md](docs/idpscim.md), [docs/idpscimcli.md](docs/idpscimcli.md), [docs/Configuration.md](docs/Configuration.md), and [docs/Development.md](docs/Development.md) for examples, flags, and the full local workflow.

* **Pre-built Binaries**
  * Download the binaries from the [GitHub Releases](https://github.com/slashdevops/idp-scim-sync/releases).

* **Container Image**
  * Pull the image from the [GitHub Container Registry](https://github.com/slashdevops/idp-scim-sync/pkgs/container/idp-scim-sync).
  * Container build and execution details are documented in [docs/idpscim.md](docs/idpscim.md), [docs/idpscimcli.md](docs/idpscimcli.md), and [docs/Development.md](docs/Development.md).

## Configurable User Fields

By default, all optional user attributes are synced from Google Workspace to AWS SSO SCIM. You can control which optional fields are included using the `sync_user_fields` configuration option.

Supported optional fields include `phoneNumbers`, `addresses`, `title`, `preferredLanguage`, `locale`, `timezone`, `nickName`, `profileURL`, `userType`, and `enterpriseData`.

Required fields are always synchronized: `name`, `userName`, `displayName`, `emails`, and `active`.

For config file examples, environment variable usage, CLI flags, SAM parameter usage, and behavior notes, see [docs/Configuration.md](docs/Configuration.md) and [docs/idpscim.md](docs/idpscim.md).

## 📦 Repositories

* 📦 [AWS Serverless Application Repository](https://serverlessrepo.aws.amazon.com/applications/us-east-1/889836709304/idp-scim-sync)
* 📦 [GitHub Container Registry](https://github.com/slashdevops/idp-scim-sync/pkgs/container/idp-scim-sync)

## ⚠️ Limitations

* **Group Limit**: The AWS SSO SCIM API has a limit of 50 groups per request. Please support the feature request on the [AWS Support site](https://repost.aws/questions/QUqqnVkIo_SYyF_SlX5LcUjg/aws-sso-scim-api-pagination-for-methods) to help get this limit increased.
* **Throttling**: With a large number of users and groups, you may encounter a `ThrottlingException` from the AWS SSO SCIM API. This project uses the [httpx](https://github.com/slashdevops/httpx) library with automatic retry and jitter backoff to mitigate this, but it's still a possibility.
* **User Status**: The Google Workspace API doesn't differentiate between normal and guest users except for their status. This project only syncs `ACTIVE` users.

## For `ssosync` Users

If you are coming from the [awslabs/ssosync](https://github.com/awslabs/ssosync) project, please note the following:

* This project only implements the `--sync-method groups`.
* This project only implements filtering for Google Workspace Groups, not Users.
* This project supports selecting which optional user attributes to sync via `--sync-user-fields` (e.g., phone numbers, addresses, enterprise data).
* The flag names are different.

## 📄 License

This project is released under the Apache License 2.0. See the [LICENSE](LICENSE) file for more details.
