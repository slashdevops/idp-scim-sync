---
applyTo: "**"
---

# Development Guidelines

Critical information for working in this codebase. Follow it precisely.

This one file is the source of truth for every assistant: `CLAUDE.md`, `AGENTS.md`,
`.cursorrules`, `.clinerules` and `DEVELOPMENT_GUIDELINES.md` are all symlinks to it. Edit this file,
never the symlinks.

## What this project is

`idp-scim-sync` keeps AWS IAM Identity Center (formerly AWS SSO) in sync with a Google Workspace
directory. It ships as an AWS Lambda function (`provided.al2023`, `arm64`, invoked on a schedule) plus
two CLIs built from `cmd/`: `idpscim` (the sync itself) and `idpscimcli` (an inspection tool for
validating configuration).

**It creates, updates and deletes real users and groups in a production identity provider.** A bug
here can remove someone's access, or grant access that should have been revoked. Correctness outranks
elegance, and a change that "probably works" is not acceptable.

## Stack

- **Language**: Go 1.27+
- **Runtime**: AWS Lambda `provided.al2023` on `arm64`
- **State**: a single JSON object in S3 — **there is no database**
- **Secrets**: AWS Secrets Manager
- **Libraries**: standard library first; `cobra`/`viper` (CLI + config), `aws-sdk-go-v2`,
  `google.golang.org/api/admin/directory/v1`, `slashdevops/httpx` (retrying HTTP client),
  `golang.org/x/sync/errgroup`
- **Logging**: `log/slog`, structured
- **Testing**: standard `testing` + `testify`, `go.uber.org/mock` (via the `tool` directive), fuzz
  targets, `testing/synctest` for concurrency
- **Build**: `make` — see the Makefile for every target
- **Lint**: `golangci-lint` against the committed `.golangci.yml`
- **CI/CD**: GitHub Actions; deployment via AWS SAM
- **Docs**: GoDoc plus `docs/` (start at [docs/Architecture.md](../docs/Architecture.md))

## Architecture invariants

Read these before changing anything in `internal/model`, `internal/core`, or the state file. Breaking
one is how you cause a silent production incident.

1. **Never construct a `model.*` value directly.** Always use the `XxxBuilder()` chain. `Build()` is
   what sets `Items` and computes `HashCode`; a hand-built struct carries a zero hash and silently
   compares unequal to everything.
2. **The hand-written `MarshalBinary`/`UnmarshalBinary` methods in `internal/model` define the hash
   input.** Hashing goes through `gob`, not JSON. Changing field order, adding a field, or removing
   one changes every hash, which invalidates every deployed state file and forces a full re-sync.
   `internal/model/golden_hash_test.go` pins the current values — if it fails, you have changed the
   wire format, and that must be deliberate and documented in `docs/Whats-New.md`.
3. **`SCIMID` is deliberately excluded from hashes.** It is assigned by AWS, not the identity
   provider, so including it would make every first-sync comparison differ.
4. **Sync decisions compare container hashes**, not the `State` hash — see
   `internal/core/actions.go`. Those container hashes must stay independent of the order resources
   arrive in, because the fan-outs build their slices from maps.
5. **Changing the state-file schema requires bumping `model.StateSchemaVersion`** and documenting a
   migration path.
6. **Interfaces are declared by the consumer**, not next to the implementation — `core.SCIMService`,
   `core.IdentityProviderService`, `core.StateRepository`, `idp.GoogleProviderService`,
   `scim.AWSSCIMProvider`, `repository.S3ClientAPI`, `aws.HTTPClient`. Keep it that way; it is what
   makes this codebase testable.
7. **`internal/core` must not import concrete AWS or Google types.** It depends only on the
   interfaces it declares.
8. **A `buildUser` that returns `nil` means "reject this record".** Callers must skip and warn, never
   store the `nil`.

## Code style

Follow Go's idiomatic style:

- [Go Style Guide](https://google.github.io/styleguide/go/guide) ·
  [Style Decisions](https://google.github.io/styleguide/go/decisions) ·
  [Best Practices](https://google.github.io/styleguide/go/best-practices) ·
  [Effective Go](https://golang.org/doc/effective_go.html)

Specifically, in this repo:

- `errors.New` for constant sentinel messages; `fmt.Errorf` with **`%w`** to wrap. Never `%v` on an
  error — `internal/core` relies on `errors.AsType` working through the chain.
- No naked `return` in anything longer than a few lines. Don't name results you don't use.
- Prefer `errgroup.WithContext` + `SetLimit` over hand-rolled `WaitGroup` + semaphore + error channel.
  Fan-outs must cancel siblings on first error.
- `context.Context` flows down from `cmd/`. `context.Background()` belongs only at the entrypoint.
- Structured `slog` with key/value pairs, never preformatted message strings.
- **Never log secrets, bearer tokens, or the service-account JSON.**
- Prefer `slices`/`maps` over `sort.Slice` and manual loops. Use `SortStableFunc` when the sort key
  isn't unique and the result feeds a hash.
- Exported identifiers need doc comments. Every package keeps a `docs.go` with a package comment.
- Meaningful names; small, single-purpose functions; comments that explain *why*, not *what*.
- Dependency injection through interfaces, for testability.

## Testing

- **Bug fixes and refactors start with a failing test.** Write the red test, confirm it fails, then
  make it green. Put the failing output in the PR description. For a pure refactor with no bug to
  reproduce, write a characterization test that pins existing behaviour instead, and confirm it
  passes both before and after — don't manufacture a fake red.
- Table-driven tests; fixtures in `testdata/`.
- `-race` is mandatory (`make test` already does it).
- Use `t.Setenv` and `t.TempDir()`, never `os.Setenv`/`os.TempDir` — subtests must not leak state
  into each other. (Three tests here once passed only because they did.)
- Mocks are generated, never hand-edited. Run `make go-generate` after any interface change.
- New parsing or decoding paths get a fuzz target. Run `make fuzz` (`FUZZ_TIME=30s` for a quick pass).
- Anything touching hashing or the state file must keep `internal/model/golden_hash_test.go` green.

## Post-change checklist

```bash
go fix ./...            # Go 1.27+ modernizers
make go-fmt             # Format
make go-betteralign     # Align struct fields for memory layout
golangci-lint run ./... # Must report 0 issues
make build              # Verify build
make test               # Full -race suite
```

## Pre-PR checklist

- Never develop on `main`. Branch from `main` for every change.
- Bring the branch up to date with `main` before opening the PR.
- Clear, descriptive commit messages explaining *why*.
- All tests pass locally, and `golangci-lint` reports no new issues versus the baseline.
- Ask whether the release is **patch, minor or major**, and make the changelog reflect it.
- Update `docs/Whats-New.md` with the change, its motivation, its impact, and the target version.
- Update `README.md` if usage, configuration, or any user-facing behaviour changed.
- Update `docs/` when behaviour, configuration or architecture changed; author diagrams as
  **mermaid** in markdown so they are reviewable in a diff, not as binary exports.
- PR description gives reviewers context and links relevant issues and docs.

## Planning and tracking

For non-trivial work, keep a plan in the gitignored `.plan/` folder with a `PROGRESS.md` tracker, and
keep it current as you go.

**Decisions that change behaviour get recorded and approved before implementation, not after.** That
includes anything affecting sync results, the state-file schema or bytes, hash values, or the AWS and
Google API calls made. When in doubt, ask rather than assume — and state the risk plainly.
