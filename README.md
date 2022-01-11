# idp-scim-sync

[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/5348/badge)](https://bestpractices.coreinfrastructure.org/projects/5348) [![CodeQL Analysis](https://github.com/slashdevops/idp-scim-sync/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/slashdevops/idp-scim-sync/actions/workflows/codeql-analysis.yml) [![Main](https://github.com/slashdevops/idp-scim-sync/actions/workflows/main.yml/badge.svg)](https://github.com/slashdevops/idp-scim-sync/actions/workflows/main.yml) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/slashdevops/idp-scim-sync?style=plastic) [![Go Report Card](https://goreportcard.com/badge/github.com/slashdevops/idp-scim-sync)](https://goreportcard.com/report/github.com/slashdevops/idp-scim-sync) [![license](https://img.shields.io/github/license/slashdevops/idp-scim-sync.svg)](https://github.com/slashdevops/idp-scim-sync/blob/main/LICENSE) [![release](https://img.shields.io/github/release/slashdevops/idp-scim-sync/all.svg)](https://github.com/slashdevops/idp-scim-sync/releases) [![Maintainability](https://api.codeclimate.com/v1/badges/8f88180aebaca6fc4923/maintainability)](https://codeclimate.com/github/slashdevops/idp-scim-sync/maintainability) [![codecov](https://codecov.io/gh/slashdevops/idp-scim-sync/branch/main/graph/badge.svg?token=H72NWJGHZ0)](https://codecov.io/gh/slashdevops/idp-scim-sync)

Keep your [AWS Single Sign-On (SSO) groups and users](https://aws.amazon.com/single-sign-on/) in sync with your [Google Workspace directory](https://workspace.google.com/) using and [AWS Lambda function](https://aws.amazon.com/lambda/).

![On AWS](https://raw.githubusercontent.com/slashdevops/idp-scim-sync/main/docs/images/diagrams/ipd-scim-sync.drawio.png)

As the image above shows, the `AWS Lambda function` is triggered by a [CloudWatch event rule](https://docs.aws.amazon.com/eventbridge/latest/userguide/eb-create-rule-schedule.html).  The event rule is configured to run every [15 minutes (default in the cfn template)](template.yaml), and sync the `AWS Single Sign-On (SSO) groups and users` with `Google Workspace directory` using the respective `APIs`. During the `first sync`, the data of the `Groups and Users` are stored in the `AWS S3 bucket` as the [State file](docs/State-File-example.md)

[The State file](docs/State-File-example.md) is a custom implementation to save time and requests to the [AWS SSO SCIM API](https://docs.aws.amazon.com/singlesignon/latest/developerguide/what-is-scim.html), also mitigate some limitations of this.

This project is developed using the [Go language](https://go.dev/) and also use [AWS SAM](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-specification.html), a tool for creating and deploying `AWS Serverless Applications` in an easy way.

If you want to know what creates the [CloudFormation Template](template.yaml), please check the [AWS SAM Template](docs/AWS-SAM-Template.md)

First time implementing [Single Sign-on on AWS](https://aws.amazon.com/single-sign-on)? please read [Using SSO](docs/Using-SSO.md)

## Important

The documentation is a __WIP__ and you can contribute!

## Repositories

* [AWS Serverless public repository - slashdevops/idp-scim-sync](https://serverlessrepo.aws.amazon.com/applications/us-east-1/889836709304/idp-scim-sync)
* [AWS ECR public repository - slashdevops/idp-scim-sync](https://gallery.ecr.aws/l2n7y5s7/slashdevops/idp-scim-sync)
* [GitHub public repository - slashdevops/idp-scim-sync](https://ghcr.io/slashdevops/idp-scim-sync)
* [Docker Hub public repository - slashdevops/idp-scim-sync](https://hub.docker.com/r/slashdevops/idp-scim-sync-linux-amd64)

## Limitations

Most of the limitations of this project are due to [AWS SSO SCIM API Limitations](https://docs.aws.amazon.com/singlesignon/latest/developerguide/what-is-scim.html).

* Use less than 50 Groups -->  [AWS SSO SCIM API (ListGroups)](https://docs.aws.amazon.com/singlesignon/latest/developerguide/listgroups.html#Constraints) has a limit of 50 Groups per request.
* Too much Users and Groups could generate a `ThrottlingException` of the some [AWS SSO SCIM API methods](https://docs.aws.amazon.com/singlesignon/latest/developerguide/what-is-scim.html)

NOTES:

1. The use of the [The State file](docs/State-File-example.md) could mitigate the number `1`, but I recommend you be cautious of these limitations as well.
2. The project implements a `well-known HTTP Retryable client` to mitigate the number `2`, but I recommend you be cautious of these limitations as well.

**Users that come from the project [SSO Sync](https://github.com/awslabs/ssosync)**

* This project only implement the `--sync-method` `groups`, so if you are using the `--sync-method` `users_groups` you can't use it, because this is going to delete and recreate your data in the AWS SSO side.
* This project only implement the `filter` for the `Google Workspace Groups`, so if you are using the `filter` for the `Google Workspace Users`, you can't use it. Please see [Using SSO](docs/Using-SSO.md) for more information.
* The flags names of this project are different from the ones of the [SSO Sync](https://github.com/awslabs/ssosync)
* Not "all the features" of the [SSO Sync](https://github.com/awslabs/ssosync) are not implemented here, and maybe will not.

## Components

1. [idpscim](docs/idpscim.md) is a program for keeping [AWS Single Sign-On (SSO) groups and users](https://aws.amazon.com/single-sign-on/) synced with [Google Workspace directory service](https://workspace.google.com/) using the [AWS SSO SCIM API](https://docs.aws.amazon.com/singlesignon/latest/developerguide/what-is-scim.html). Details [here](docs/idpscim.md).
2. [idpscimcli](docs/idpscimcli.md) is is a command-line tool to check and validate some functionalities implemented in `idpscim`. Details [here](docs/idpscimcli.md).

## Requirements

* [AWS Single Sign-On -> Connect to your external identity provider](https://docs.aws.amazon.com/singlesignon/latest/userguide/manage-your-identity-source-idp.html)
* [AWS Single Sign-On -> Automatic provisioning](https://docs.aws.amazon.com/singlesignon/latest/userguide/provision-automatically.html)

## How to use

To use this project you have different options, and depending on your needs you can use the following

### In AWS

There are two ways to use this project in AWS and described below.

#### Using AWS Serverless Repository

This is the easy way, this project is deployed as a [AWS Serverless Application](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-apps-overview.html) in [AWS Serverless Application Repository](https://aws.amazon.com/es/serverless/serverlessrepo/).

The public repository of the project is [slashdevops/idp-scim-sync](https://serverlessrepo.aws.amazon.com/applications/us-east-1/889836709304/idp-scim-sync)

NOTE: The repository depends on your `AWS Region`.

#### Using AWS SAM

This is the way if you want to build an deploy the lambda function from your local machine.

Requirements:

1. [git](https://git-scm.com/)
2. [Go](https://go.dev/learn/)
3. [AWS SAM Cli](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)

Validate, Build and Deploy:

```bash
aws cloudformation validate-template --template-body file://template.yaml 1>/dev/null
sam validate

sam deploy --guided
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