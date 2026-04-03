---
applyTo: "**"
---

# Development Guidelines

This document contains the critical information about working with the project codebase.
Follows these guidelines precisely to ensure consistency and maintainability of the code.

## Stack

- Language: Go (Go 1.26+)
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

## Post-Change Checklist

Prefer the Make targets that the repo already defines after making changes:

```bash
go fix ./...            # Optional manual step for Go 1.26+ syntax
make go-fmt             # Format code
make go-betteralign     # Align struct fields for optimal memory layout
golangci-lint run ./... # Run linter (also checks formatting, vet, and other issues)
make build               # Verify build
make test                # Run tests
```

## Pre-PR Checklist

- never use the `main` branch for development. Always create a new feature branch from `main` for your changes.
- ensure your branch is up to date with `main` before creating a pull request
- write clear and descriptive commit messages that explain the purpose of each change
- ensure all tests pass locally before pushing your changes
- ask the version is going to be released as a patch, minor, or major update and ensure the change log reflects this
- update the `docs/Whats-New.md` file with a summary of the change, including the motivation and impact and the version it will be released in
- ensure the PR description is clear and concise, providing context for reviewers and linking to any relevant issues or documentation
- update the `README.md` file if the change affects usage instructions, configuration, or any other user-facing documentation
