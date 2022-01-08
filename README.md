# idp-scim-sync

[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/5348/badge)](https://bestpractices.coreinfrastructure.org/projects/5348)
[![Main branch workflow](https://github.com/slashdevops/idp-scim-sync/actions/workflows/main.yaml/badge.svg?branch=main)](https://github.com/slashdevops/idp-scim-sync/actions/workflows/main.yaml)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/slashdevops/idp-scim-sync?style=plastic)
[![Go Report Card](https://goreportcard.com/badge/github.com/slashdevops/idp-scim-sync)](https://goreportcard.com/report/github.com/slashdevops/idp-scim-sync)
[![license](https://img.shields.io/github/license/slashdevops/idp-scim-sync.svg)](https://github.com/slashdevops/idp-scim-sync/blob/main/LICENSE)
[![release](https://img.shields.io/github/release/slashdevops/idp-scim-sync/all.svg)](https://github.com/slashdevops/idp-scim-sync/releases)
[![Maintainability](https://api.codeclimate.com/v1/badges/8f88180aebaca6fc4923/maintainability)](https://codeclimate.com/github/slashdevops/idp-scim-sync/maintainability)
[![codecov](https://codecov.io/gh/slashdevops/idp-scim-sync/branch/main/graph/badge.svg?token=H72NWJGHZ0)](https://codecov.io/gh/slashdevops/idp-scim-sync)

Keep your [AWS Single Sign-On (SSO) groups and users](https://aws.amazon.com/single-sign-on/) in sync with your [Google Workspace directory](https://workspace.google.com/) using and [AWS Lambda function](https://aws.amazon.com/lambda/).

![On AWS](https://raw.githubusercontent.com/slashdevops/idp-scim-sync/main/docs/images/diagrams/ipd-scim-sync.drawio.png)

This project is [100% Golang](https://go.dev/) and has two main components:

1. [idpscim](docs/idpscim.md) is a program for keeping AWS Single Sing-On Groups and Users sync with Google Workspace service using SCIM protocol. More details [here](docs/idpscim.md).

2. [idpscimcli](docs/idpscimcli.md) is is a command-line tool to check and validate some functionalities implemented in `idpscim`. More details [here](docs/idpscimcli.md).

## Examples

look here [docs/Demo.md](docs/Demo.md)

## License

This module is released under the Apache License Version 2.0:

* [http://www.apache.org/licenses/LICENSE-2.0.html](http://www.apache.org/licenses/LICENSE-2.0.html)