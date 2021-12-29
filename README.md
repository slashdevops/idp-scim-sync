# idp-scim-sync

[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/5348/badge)](https://bestpractices.coreinfrastructure.org/projects/5348)
[![CI Checks](https://github.com/slashdevops/idp-scim-sync/actions/workflows/main.yaml/badge.svg)](https://github.com/slashdevops/idp-scim-sync/actions)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/slashdevops/idp-scim-sync?style=plastic)
[![Go Report Card](https://goreportcard.com/badge/github.com/slashdevops/idp-scim-sync)](https://goreportcard.com/report/github.com/slashdevops/idp-scim-sync)
[![license](https://img.shields.io/github/license/slashdevops/idp-scim-sync.svg)](https://github.com/slashdevops/idp-scim-sync/blob/main/LICENSE)
[![release](https://img.shields.io/github/release/slashdevops/idp-scim-sync/all.svg)](https://github.com/slashdevops/idp-scim-sync/releases)
[![Maintainability](https://api.codeclimate.com/v1/badges/8f88180aebaca6fc4923/maintainability)](https://codeclimate.com/github/slashdevops/idp-scim-sync/maintainability)
[![codecov](https://codecov.io/gh/slashdevops/idp-scim-sync/branch/main/graph/badge.svg?token=H72NWJGHZ0)](https://codecov.io/gh/slashdevops/idp-scim-sync)

This project is composed of two main components:

1. [idpscim](cmd/idpscim/cmd/root.go) is a program for keeping Google Workspace Groups and Users sync with AWS Single Sing-On service using SCIM protocol.

2. [idpscimcli](cmd/idpscimcli/cmd/root.go) is is a command-line tool to check and validate some functionalities implemented in `idpscim`

## idpscim

This program could work in three ways:

1. As an [AWS Lambda function](https://aws.amazon.com/lambda/?nc1=h_ls) deployed via [AWS SAM](https://aws.amazon.com/serverless/sam/) or consumed directly from the [AWS Serverless Application Repository](https://aws.amazon.com/serverless/serverlessrepo/?nc1=h_ls)
2. As a command line tool
3. As a Docker container

## idpscimcli

```cmd
idpscimcli --help
```
