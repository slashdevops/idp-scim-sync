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

If you are intentionally modernizing syntax or APIs for Go 1.26+, run `go fix ./...` manually as a separate step.

Always keep `README.md` and `docs/` updated for any architectural changes, new commands/flags, removed flags, new packages, or changes to the development workflow. This includes command usage examples, flag tables, and CLI reference sections.
