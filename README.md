# idp-scim-sync

[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/5348/badge)](https://bestpractices.coreinfrastructure.org/projects/5348)
[![CodeQL](https://github.com/slashdevops/idp-scim-sync/actions/workflows/codeql.yml/badge.svg?branch=main)](https://github.com/slashdevops/idp-scim-sync/actions/workflows/codeql.yml)
[![Gosec](https://github.com/slashdevops/idp-scim-sync/actions/workflows/gosec.yml/badge.svg?branch=main)](https://github.com/slashdevops/idp-scim-sync/actions/workflows/gosec.yml)
[![Build](https://github.com/slashdevops/idp-scim-sync/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/slashdevops/idp-scim-sync/actions/workflows/build.yml)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/slashdevops/idp-scim-sync?style=plastic)
[![Go Report Card](https://goreportcard.com/badge/github.com/slashdevops/idp-scim-sync)](https://goreportcard.com/report/github.com/slashdevops/idp-scim-sync)
[![license](https://img.shields.io/github/license/slashdevops/idp-scim-sync.svg)](https://github.com/slashdevops/idp-scim-sync/blob/main/LICENSE)
[![Release](https://github.com/slashdevops/idp-scim-sync/actions/workflows/release.yml/badge.svg)](https://github.com/slashdevops/idp-scim-sync/actions/workflows/release.yml)
[![release](https://img.shields.io/github/release/slashdevops/idp-scim-sync/all.svg)](https://github.com/slashdevops/idp-scim-sync/releases)
[![Maintainability](https://api.codeclimate.com/v1/badges/8f88180aebaca6fc4923/maintainability)](https://codeclimate.com/github/slashdevops/idp-scim-sync/maintainability)
[![codecov](https://codecov.io/gh/slashdevops/idp-scim-sync/branch/main/graph/badge.svg?token=H72NWJGHZ0)](https://codecov.io/gh/slashdevops/idp-scim-sync)

Keep your [AWS IAM Identity Center (Successor to AWS Single Sign-On)](https://aws.amazon.com/iam/identity-center/) in sync with your [Google Workspace directory](https://workspace.google.com/) using and [AWS Lambda function](https://aws.amazon.com/lambda/).

![On AWS](https://raw.githubusercontent.com/slashdevops/idp-scim-sync/main/docs/images/diagrams/ipd-scim-sync.drawio.png)

As the image above shows, the [AWS Lambda function](https://aws.amazon.com/lambda) is triggered by a [CloudWatch event rule](https://docs.aws.amazon.com/eventbridge/latest/userguide/eb-create-rule-schedule.html), the event rule is configured to run every [15 minutes (default in the cfn template)](template.yaml), and `sync` the [AWS IAM Identity Center (Successor to AWS Single Sign-On)](https://aws.amazon.com/iam/identity-center/) with `Google Workspace directory` using their respective `APIs`.  During the `first sync`, the data of the `Groups and Users` are stored in the `AWS S3 bucket` as [the State file](docs/State-File-example.md)

[The State file](docs/State-File-example.md) is a custom implementation to save time and requests to the [AWS SSO SCIM API](https://docs.aws.amazon.com/singlesignon/latest/developerguide/what-is-scim.html), also mitigate some limitations of this.

This project is developed using the [Go language](https://go.dev/) and [AWS SAM](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-specification.html), a tool for creating, [publishing](https://aws.amazon.com/serverless/serverlessrepo) and deploying `AWS Serverless Applications` in an easy way.

If you want to know what creates the [CloudFormation Template](template.yaml), please check the [AWS SAM Template](docs/AWS-SAM-Template.md)

__First time implementing [AWS IAM Identity Center (Successor to AWS Single Sign-On)](https://aws.amazon.com/iam/identity-center/)? please read [Using SSO](docs/Using-SSO.md)__

The best way to to deploy and use this is through the [AWS Serverless public repository - slashdevops/idp-scim-sync](https://serverlessrepo.aws.amazon.com/applications/us-east-1/889836709304/idp-scim-sync)

## Compatibility

AWS recently announced `AWS Lambda Deprecates Go Runtime 1.x` and posted this article [Migrating AWS Lambda functions from the Go1.x runtime to the custom runtime on Amazon Linux 2](https://aws.amazon.com/blogs/compute/migrating-aws-lambda-functions-from-the-go1-x-runtime-to-the-custom-runtime-on-amazon-linux-2/) to help customers with the migration.

This project is already migrated since version `v0.0.19` to the `provided.al2` runtime and `arm64` architecture, so you can use it without any problem.

|   version   | AWS Lambda Runtime | Architecture       | Deprecation Date |
|-------------|--------------------|--------------------|------------------|
| <= v0.0.18  | Go 1.x             | amd64 (Intel)      | 2023-12-31       |
| >= v0.0.19  | provided.al2       | arm64 (Graviton 2) | ----------       |

## Features

* Efficient data retrieval from Google Workspace API using [Partial response](https://cloud.google.com/storage/docs/json_api#partial-response)
* Supported nested groups in Google Workspace thanks to [includeDerivedMembership](https://developers.google.com/admin-sdk/directory/reference/rest/v1/members/list#query-parameters) API Query Parameter
* Could be used or deployed via `AWS Serverless repository (Public)`, `Container Image` or `CLI`. See [Repositories](#Repositories)
* Incremental changes, drastically reduced the number of requests to the [AWS SSO SCIM API](https://docs.aws.amazon.com/singlesignon/latest/developerguide/what-is-scim.html) thanks to the implementation of [State file](docs/State-File-example.md)

## Important

The documentation is a __WIP__ and you can contribute!

## Credentials

* __Google Workspace API credentials__
  This application will need [Google Workspace Directory API](https://developers.google.com/admin-sdk/directory/v1/guides) to retrieve `Users`, `Groups` and `Members` data.  Configure this is a little `bit tricky`, but it is well [documented by Google](https://developers.google.com/workspace/guides/create-credentials).
  The Authorization/Authorization needed is [OAuth 2.0 for Server to Server Applications](https://developers.google.com/identity/protocols/oauth2/service-account) and this require to:
   1. [Create a Service Account](https://developers.google.com/identity/protocols/oauth2/service-account#creatinganaccount) on Google Cloud Platform
   2. [Delegate domain-wide authority to the service account](https://developers.google.com/identity/protocols/oauth2/service-account#delegatingauthority), the scope needed are:
      1. `https://www.googleapis.com/auth/admin.directory.group.readonly`
      2. `https://www.googleapis.com/auth/admin.directory.user.readonly`
      3. `https://www.googleapis.com/auth/admin.directory.group.member.readonly`
* __AWS SSO SCIM API credentials__
  This credentials is configured in the [AWS IAM Identity Center (Successor to AWS Single Sign-On)](https://aws.amazon.com/iam/identity-center/) service following the [Automatic provisioning](https://docs.aws.amazon.com/singlesignon/latest/userguide/provision-automatically.html) guide.

## Repositories

* [AWS Serverless public repository - slashdevops/idp-scim-sync](https://serverlessrepo.aws.amazon.com/applications/us-east-1/889836709304/idp-scim-sync)
* [AWS ECR public repository - slashdevops/idp-scim-sync](https://gallery.ecr.aws/l2n7y5s7/slashdevops/idp-scim-sync)
* [GitHub public repository - slashdevops/idp-scim-sync](https://github.com/slashdevops/idp-scim-sync/pkgs/container/idp-scim-sync)
* [Docker Hub public repository - slashdevops/idp-scim-sync](https://hub.docker.com/r/slashdevops/idp-scim-sync)

## Limitations

Most of the limitations of this project are due to [AWS SSO SCIM API Limitations](https://docs.aws.amazon.com/singlesignon/latest/developerguide/what-is-scim.html).

* Use less than 50 Groups -->  [AWS SSO SCIM API (ListGroups)](https://docs.aws.amazon.com/singlesignon/latest/developerguide/listgroups.html#Constraints) has a limit of 50 Groups per request.  I created these tickets in AWS Support site [AWS SSO SCIM API pagination for methods](https://repost.aws/questions/QUqqnVkIo_SYyF_SlX5LcUjg/aws-sso-scim-api-pagination-for-methods) and [AWS SSO SCIM API ListGroups members](https://repost.aws/questions/QURqsaKxH9SqWYsBJ9UDdAPg/aws-sso-scim-api-list-groups-members), `please consider supporting this ticket with your` 👍.
* Too much Users and Groups could generate a `ThrottlingException` of the some [AWS SSO SCIM API methods](https://docs.aws.amazon.com/singlesignon/latest/developerguide/what-is-scim.html)
* Google Workspace API doesn't separate normal and guest users expect for status (guest miss status), so only `ACTIVE` users are collected to model as group members. Logically all users who are wanted (and capable of) to sign in are `ACTIVE`.

NOTES:

1. The use of the [The State file](docs/State-File-example.md) could mitigate the number `1`, but I recommend you be cautious of these limitations as well.
2. The project implements a [well-known HTTP Retryable client (/go-retryablehttp)](https://github.com/hashicorp/go-retryablehttp) to mitigate the number `2`, but I recommend you be cautious of these limitations as well.

### Users that come from the project [SSO Sync](https://github.com/awslabs/ssosync)

* This project only implements the `--sync-method` `groups`, so if you are using the `--sync-method` `users_groups` you can't use it, because this is going to delete and recreate your data in the AWS SSO side.
* This project only implements the `filter` for the `Google Workspace Groups`, so if you are using the `filter` for the `Google Workspace Users`, you can't use it. Please see [Using SSO](docs/Using-SSO.md) for more information.
* The flags names of this project are different from the ones of the [SSO Sync](https://github.com/awslabs/ssosync)
* Not "all the features" of the [SSO Sync](https://github.com/awslabs/ssosync) are implemented here, and maybe will not be.

## Components

1. [idpscim](docs/idpscim.md) is a program for keeping [AWS Single Sign-On (SSO) groups and users](https://aws.amazon.com/single-sign-on/) synced with [Google Workspace directory service](https://workspace.google.com/) using the [AWS SSO SCIM API](https://docs.aws.amazon.com/singlesignon/latest/developerguide/what-is-scim.html). Details [here](docs/idpscim.md).
2. [idpscimcli](docs/idpscimcli.md) is a command-line tool to check and validate some functionalities implemented in `idpscim`. Details [here](docs/idpscimcli.md).

## Requirements

* [AWS Single Sign-On -> Connect to your external identity provider](https://docs.aws.amazon.com/singlesignon/latest/userguide/manage-your-identity-source-idp.html)
* [AWS Single Sign-On -> Automatic provisioning](https://docs.aws.amazon.com/singlesignon/latest/userguide/provision-automatically.html)

## How to use

To use this project you have different options, and depending on your needs you can use the following

### In AWS

There are two ways to use this project in AWS and described below.

#### Using AWS Serverless Repository

This is the easy way, this project is deployed as an [AWS Serverless Application](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-apps-overview.html) in [AWS Serverless Application Repository](https://aws.amazon.com/es/serverless/serverlessrepo/).

The public repository of the project is [slashdevops/idp-scim-sync](https://serverlessrepo.aws.amazon.com/applications/us-east-1/889836709304/idp-scim-sync)

NOTE: The repository depends on your `AWS Region`.

#### Using AWS SAM

This is the way if you want to build and deploy the lambda function from your local machine.

Requirements:

1. [git](https://git-scm.com/)
2. [Go](https://go.dev/learn/)
3. [AWS SAM Cli](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)

Validate, Build and Deploy:

```bash
# your AWS Cli Profile and Region
export AWS_PROFILE=<profile name here>
export AWS_REGION=<region here>

# validate
aws cloudformation validate-template --template-body file://template.yaml 1>/dev/null --profile $AWS_PROFILE
sam validate --profile $AWS_PROFILE

# build
sam build --profile $AWS_PROFILE

# deploy guided
sam deploy --guided  --capabilities CAPABILITY_IAM --capabilities CAPABILITY_NAMED_IAM --profile $AWS_PROFILE
```

Are you using [AWS Cli Profiles?](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-profiles.html), read [AWS-SAM](docs/AWS-SAM.md)

### In Local

You will have two ways to execute the binaries of this project in local, building these or using the pre-built stored in the github repository project.

#### Building the project

To build the project in local, you will need to have installed and configured at least the following:

1. [git](https://git-scm.com/)
2. [Go](https://go.dev/learn/)
3. [make](https://www.gnu.org/software/make/)

Then you will need to clone the repository in your local machine, and execute the following commands:

* Compile for your Operating System:

```bash
make
```

then the binaries are in `build/` folder.

* `Cross-compiling` the project for `Windows`, `MacOS` and `Linux` (default)

```bash
make clean
make test # optional
make build-dist
```

then the binaries are in `dist/` folder.

* Others Operating Systems, see the list of supported platforms in the [syslist.go](https://github.com/golang/go/blob/master/src/go/build/syslist.go)

```bash
make clean
GO_OS=<something from goosList in syslist.go> GO_ARCH=<something from goarchList in syslist.go> make test # optional
GO_OS=<something from goosList in syslist.go> GO_ARCH=<something from goarchList in syslist.go> make build-dist
```

then the binaries are in `dist/` folder.

* Execute

```bash
./idpscim --help
#or
./idpscimcli --help
```

#### Using the pre-built binaries

This is the easy way, just download the binaries you need from the [github repository releases](https://github.com/slashdevops/idp-scim-sync/releases)

and see [Execute](#executing) the binaries.

#### Using the pre-built binaries in local

Example [docs/Demo.md](https://github.com/slashdevops/idp-scim-sync/blob/main/docs/Demo.md)

## License

This module is released under the Apache License Version 2.0:

* [http://www.apache.org/licenses/LICENSE-2.0.html](http://www.apache.org/licenses/LICENSE-2.0.html)
