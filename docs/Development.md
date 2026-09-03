# Development

This document covers the day-to-day development workflow for this repository, including local Go development, the two CLI programs built from `cmd/`, and the AWS SAM workflow used to deploy the Lambda function.

## Prerequisites

Install the following tools before contributing:

- [Git](https://git-scm.com/)
- [Go](https://go.dev/doc/install) **1.27 or newer** — the `go` directive in `go.mod` is `1.27.0`, and
  `GOTOOLCHAIN=auto` (the default) will fetch a matching toolchain if yours is older
- `make`
- [golangci-lint](https://golangci-lint.run/) **v2.x** — configured by the committed
  [`.golangci.yml`](../.golangci.yml)
- [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html) if you want to deploy or test in AWS
- [AWS SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/install-sam-cli.html) if you want to work with the serverless deployment
- [podman](https://podman.io/) if you want to build or publish the container image

## Main Entry Points

This repository builds two user-facing programs:

| Program | Source directory | Purpose |
| --- | --- | --- |
| `idpscim` | [cmd/idpscim](../cmd/idpscim) | Main synchronization program used by Lambda, local CLI execution, and the container image |
| `idpscimcli` | [cmd/idpscimcli](../cmd/idpscimcli) | Helper CLI used to inspect Google Workspace and AWS SCIM data while validating configuration |

Local builds are written to `build/`:

- `build/idpscim`
- `build/idpscimcli`

## Recommended Local Workflow

Run the repository quality checks in the same order the project expects:

```bash
go fix ./...            # Go 1.27+ modernizers
make go-fmt             # Format
make go-betteralign     # Struct field alignment for memory layout
golangci-lint run ./... # Must report 0 issues
make build              # Verify build
make test               # Full -race suite
```

`golangci-lint` reads the committed [`.golangci.yml`](../.golangci.yml), so your results match CI.
**The baseline is 0 issues** — anything reported is from your change.

### Additional checks

Run these when the relevant area changes:

```bash
# After changing any consumed interface — mocks are generated, never hand-edited
make go-generate

# After changing parsing or decoding (state file, SCIM responses)
FUZZ_TIME=30s make fuzz

# After changing dependencies
go mod verify
go run golang.org/x/vuln/cmd/govulncheck@latest ./...

# Coverage report in your browser
make go-test-coverage
```

### Updating dependencies

```bash
make go-mod-update   # go mod tidy, then `go get -u` for each direct dependency, then tidy again
make go-generate     # regenerate mocks in case an interface changed shape
make test
```

> [!IMPORTANT]
> Any change under `internal/model` must keep
> [`golden_hash_test.go`](../internal/model/golden_hash_test.go) green. It pins the hash values that
> the S3 state file is compared by; if it fails, you have changed a wire format and every deployed
> installation will do a full re-sync. See the
> [Implementation Guide](Implementation-Guide.md#-state-file-compatibility).

### A note for macOS developers on Go 1.27

Go 1.27 made `crypto/x509.SystemCertPool` honour `SSL_CERT_FILE` and `SSL_CERT_DIR` on macOS and
Windows. If either is set in your shell, local TLS verification will use those roots instead of the
system keychain, which can cause confusing certificate failures against Google or AWS endpoints.
Unset them, or set `GODEBUG=x509sslcertoverrideplatform=0`. Lambda runs on Linux, where this was
already the behaviour.

## Build And Run The Programs Locally

Build both binaries:

```bash
make build
```

Show the help for each program:

```bash
./build/idpscim --help
./build/idpscimcli --help
```

Useful `idpscimcli` validation commands during development:

```bash
./build/idpscimcli aws service config --help
./build/idpscimcli aws groups list --help
./build/idpscimcli gws groups list --help
./build/idpscimcli gws groups members list --help
./build/idpscimcli gws users list --help
```

## Build Distribution Artifacts

To cross-compile the binaries for the supported operating systems and architectures:

```bash
make build-dist
```

This writes the artifacts to `dist/`.

## Use AWS SAM During Development

AWS SAM is the deployment path for the Lambda version of `idpscim`.

### Validate And Build

```bash
export AWS_PROFILE=<profile-name>
export AWS_REGION=us-east-1

aws cloudformation validate-template \
  --template-body file://template.yaml \
  --profile "$AWS_PROFILE"

sam validate \
  --profile "$AWS_PROFILE" \
  --region "$AWS_REGION"

GIT_VERSION=dev sam build \
  --profile "$AWS_PROFILE" \
  --region "$AWS_REGION"
```

### First Deploy

Use the guided workflow the first time so you can enter the required parameters interactively:

```bash
sam deploy --guided \
  --stack-name idp-scim-sync \
  --capabilities CAPABILITY_IAM CAPABILITY_NAMED_IAM \
  --profile "$AWS_PROFILE" \
  --region "$AWS_REGION"
```

### Update An Existing Development Stack

After the first guided deployment, iterate with:

```bash
GIT_VERSION=dev sam build \
  --profile "$AWS_PROFILE" \
  --region "$AWS_REGION"

sam deploy \
  --stack-name idp-scim-sync \
  --capabilities CAPABILITY_IAM CAPABILITY_NAMED_IAM \
  --profile "$AWS_PROFILE" \
  --region "$AWS_REGION"
```

Use this when you want to validate application changes against real AWS resources.

## Container Workflow

Container images are published to [GitHub Container Registry](https://github.com/slashdevops/idp-scim-sync/pkgs/container/idp-scim-sync) using `podman`.

Build locally:

```bash
make build-dist
GIT_VERSION=test make container-build
podman images | grep idp-scim-sync
```

Publish:

```bash
REPOSITORY_REGISTRY_TOKEN=<your-token> \
REPOSITORY_REGISTRY_USERNAME=<your-username> \
make container-login

GIT_VERSION=<version> make container-publish
```

## Related Documentation

- [Implementation-Guide.md](Implementation-Guide.md) — invariants, how to add a synced attribute,
  testing conventions, state-file compatibility rules
- [Architecture.md](Architecture.md) — package topology, sync algorithm, hashing
- [AWS-SAM.md](AWS-SAM.md) for private deploys, SAR updates, and public publishing
- [idpscim.md](idpscim.md) for the main sync program
- [idpscimcli.md](idpscimcli.md) for the validation CLI
- [Release.md](Release.md) for the release process
