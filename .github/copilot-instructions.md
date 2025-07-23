---
applyTo: "**"
---

# Development Guidelines

This document contains the critical information about working with the project codebase.
Follows these guidelines precisely to ensure consistency and maintainability of the code.

## Stack

- Language: Go (Go 1.24+)
- Framework: Go standard library
- Testing: Go's built-in testing package
- Build Tool: `make` using the Makefile with all the targets defined to build, test, and run the application
- Dependency Management: Go modules
- Version Control: Git
- Documentation: GoDoc
- Code Review: Pull requests on GitHub
- CI/CD: GitHub Actions
- Database: PostgreSQL
- Logging: `slog` package from the standard library

## Project Structure

This project implements a kind of Model-View-Controller (MVC) architecture, the main components are:

- Core: located in the `internal/core` package, it contains the business logic of the application.
- Config: located in the `internal/config` package, it contains the configuration of the application
- Model: located in the `internal/model` package, it contains the data structures or entities used in the application.
- Repository: located in the `internal/repository` package, it contains the data access layer for interacting with the database.
- Identity Provider conversion: located in the `internal/idp/` package, This is the glue between the logic into `internal/core` and the identity provider located at `pkg/google/` package.
- SCIM conversion: located in the `internal/scim/` package, This is the glue between the logic into `internal/core` and the SCIM API located at `pkg/aws/` package.

The project structure is organized as follows:

```plaintext
.
├── build
│   ├── coverage.txt
│   ├── idpscim
│   └── idpscimcli
├── cmd
│   ├── idpscim
│   │   ├── cmd
│   │   │   └── root.go
│   │   └── main.go
│   └── idpscimcli
│       ├── cmd
│       │   ├── aws.go
│       │   ├── common.go
│       │   ├── gws.go
│       │   └── root.go
│       └── main.go
├── CODE_OF_CONDUCT.md
├── CONTRIBUTING.md
├── coverage.out
├── DCO
├── Dockerfile
├── docs
│   ├── AWS-SAM-Template.md
│   ├── AWS-SAM.md
│   ├── Configuration.md
│   ├── Demo.md
│   ├── Development.md
│   ├── idpscim.md
│   ├── idpscimcli.md
│   ├── images
│   │   ├── demo
│   │   │   ├── aws-groups-developers.png
│   │   │   ├── aws-groups.png
│   │   │   ├── aws-s3-state-file.png
│   │   │   ├── aws-users.png
│   │   │   ├── gws-groups-developers.png
│   │   │   ├── gws-groups.png
│   │   │   ├── gws-users.png
│   │   │   ├── idpscim-help.png
│   │   │   ├── state-file.png
│   │   │   ├── sync-1.png
│   │   │   └── sync-2.png
│   │   └── diagrams
│   │       ├── idpscim--workflow.drawio.html
│   │       ├── idpscim--workflow.drawio.png
│   │       └── ipd-scim-sync.drawio.png
│   ├── Release.md
│   ├── State-File-example.md
│   └── Using-SSO.md
├── extras
│   └── infra
│       ├── 1_cfn-slashdevops-iam-oidc-github-provider.template
│       ├── 2_cfn-slashdevops-iam-role-idp-scim-sync.template
│       ├── 3_cfn-slashdevops-ecr-repo-idp-scim-sync.template
│       ├── 4_cfn-slashdevops-iam-policy-idp-scim-sync-ecr.template
│       ├── 4_cfn-slashdevops-s3-idp-scim-sync-sam.template
│       ├── 5_cfn-slashdevops-iam-policy-idp-scim-sync-sam.template
│       └── README.md
├── go.mod
├── go.sum
├── internal
│   ├── config
│   ├── core
│   ├── deepcopy
│   ├── idp
│   ├── model
│   ├── repository
│   ├── scim
│   └── version
├── LICENSE
├── Makefile
├── mocks
├── pkg
│   ├── aws
│   └── google
├── README.md
├── SECURITY.md
└── template.yaml
```

## Code Style

- Follow Go's idiomatic style defined in
  - [Go Style Guide](https://google.github.io/styleguide/go/guide)
  - [Go Style Decisions](https://google.github.io/styleguide/go/decisions)
  - [Go Style Best Practices](https://google.github.io/styleguide/go/best-practices)
  - [Effective Go](https://golang.org/doc/effective_go.html)
- Use meaningful names for variables, functions, and packages.
- Keep functions small and focused on a single task.
- Use comments to explain complex logic or decisions.
- Use dependency injection for services and repositories to facilitate testing and maintainability.

## Testing

- Write unit tests implementing test tables for all functions when applicable.
- Use the `testing` package for writing tests.
- Use `go test` to run tests.
- Use `go test -cover` to check code coverage.
- Use `go test -race` to check for race conditions.
- Use `go test -bench` to run benchmarks.
- Use `go test -v` for verbose output during testing.
- Use `go test -run` to run specific tests.
- Use `go test -short` to run short tests.
