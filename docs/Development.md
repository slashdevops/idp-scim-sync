# Development

This document covers the day-to-day development workflow for this repository, including local Go development, the two CLI programs built from `cmd/`, and the AWS SAM workflow used to deploy the Lambda function.

## Prerequisites

Install the following tools before contributing:

- [Git](https://git-scm.com/)
- [Go](https://go.dev/doc/install)
- `make`
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
make go-fmt
make go-betteralign
golangci-lint run ./...
make build
make test
```

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

- [AWS-SAM.md](AWS-SAM.md) for private deploys, SAR updates, and public publishing
- [idpscim.md](idpscim.md) for the main sync program
- [idpscimcli.md](idpscimcli.md) for the validation CLI
