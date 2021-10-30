# idp-scim-sync

[![CII Best Practices](https://bestpractices.coreinfrastructure.org/en/projects/5348/badge)](https://bestpractices.coreinfrastructure.org/en/projects/5348)
[![e2e](https://github.com/slashdevops/idp-scim-sync/workflows/e2e/badge.svg)](https://github.com/slashdevops/idp-scim-sync/actions)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/slashdevops/idp-scim-sync?style=plastic)
[![Go Report Card](https://goreportcard.com/badge/github.com/slashdevops/idp-scim-sync)](https://goreportcard.com/report/github.com/slashdevops/idp-scim-sync)
[![license](https://img.shields.io/github/license/slashdevops/idp-scim-sync.svg)](https://github.com/slashdevops/idp-scim-sync/blob/main/LICENSE)
[![release](https://img.shields.io/github/release/slashdevops/idp-scim-sync/all.svg)](https://github.com/slashdevops/idp-scim-sync/releases)
[![Maintainability](https://api.codeclimate.com/v1/badges/8f88180aebaca6fc4923/maintainability)](https://codeclimate.com/github/slashdevops/idp-scim-sync/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/8f88180aebaca6fc4923/test_coverage)](https://codeclimate.com/github/slashdevops/idp-scim-sync/test_coverage)

`idpscim` is a for keeping Google Workspace Groups and Users with AWS Single Sing-On service using SCIM protocol.

`idpscimcli` is a command line tool to check and validate some functionalities implemented in `idpscim`

## Available Commands

### idpscimcli

```cmd
idpscimcli gws groups list -u "user.email@google.com" -s "./credentials.json" -q "name:Admin*" -q "name:SuperAdmin*"
```
