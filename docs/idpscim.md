# idpscim

`idpscim` is the main synchronization program in this repository. It reads Google Workspace groups and members, compares them with AWS IAM Identity Center through the SCIM API, and stores synchronization state in S3 so later runs can avoid unnecessary updates.

This is the program executed by the deployed Lambda function.

## Source And Build Output

| Item | Location |
| --- | --- |
| Source entry point | [cmd/idpscim/main.go](../cmd/idpscim/main.go) |
| Local binary | `build/idpscim` |
| Lambda build artifact | `.aws-sam/build/LambdaFunction/bootstrap` |

## Supported Run Modes

`idpscim` can run in three ways:

1. As an AWS Lambda function deployed with AWS SAM or consumed from the AWS Serverless Application Repository.
2. As a local command-line program.
3. As a container image.

For the full configuration model, see [Configuration.md](Configuration.md).

## Key Flags

The current `idpscim --help` surface is summarized below.

### General And Logging

| Flag | Purpose |
| --- | --- |
| `--config-file`, `-c` | Path to the configuration file |
| `--debug`, `-d` | Shortcut to force debug logging |
| `--log-format`, `-f` | Log output format |
| `--log-level`, `-l` | Log verbosity |
| `--version`, `-v` | Show version information |

### Google Workspace Input

| Flag | Purpose |
| --- | --- |
| `--gws-service-account-file`, `-s` | Path to the Google Workspace service account JSON |
| `--gws-user-email`, `-u` | Delegated Google Workspace user email |
| `--gws-groups-filter`, `-q` | One or more filters that restrict which groups are synchronized |
| `--gws-service-account-file-secret-name`, `-o` | Secret name used when resolving the service account JSON from AWS Secrets Manager |
| `--gws-user-email-secret-name`, `-p` | Secret name used when resolving the delegated user email from AWS Secrets Manager |

### AWS SCIM And State Storage

| Flag | Purpose |
| --- | --- |
| `--aws-scim-endpoint`, `-e` | AWS IAM Identity Center SCIM endpoint |
| `--aws-scim-access-token`, `-t` | AWS IAM Identity Center SCIM access token |
| `--aws-scim-endpoint-secret-name`, `-n` | Secret name used when resolving the SCIM endpoint from AWS Secrets Manager |
| `--aws-scim-access-token-secret-name`, `-j` | Secret name used when resolving the SCIM token from AWS Secrets Manager |
| `--aws-s3-bucket-name`, `-b` | S3 bucket used to store the sync state |
| `--aws-s3-bucket-key`, `-k` | S3 object key used for the sync state |
| `--use-secrets-manager`, `-g` | Tell the program to load values from AWS Secrets Manager |

### Sync Behavior

| Flag | Purpose |
| --- | --- |
| `--sync-method`, `-m` | Sync strategy. The implemented value is `groups` |
| `--sync-user-fields` | Optional user fields to synchronize |

## Example Local Run

Build the binary:

```bash
make build
```

Run the program with direct credentials and a state bucket:

```bash
./build/idpscim \
  --gws-service-account-file credentials.json \
  --gws-user-email admin@example.com \
  --gws-groups-filter 'name:AWS*' \
  --aws-scim-endpoint https://example.awsapps.com/scim/v2/ \
  --aws-scim-access-token "$SCIM_ACCESS_TOKEN" \
  --aws-s3-bucket-name idp-scim-sync-state-123456789012-us-east-1 \
  --aws-s3-bucket-key data/state.json \
  --sync-method groups \
  --sync-user-fields phoneNumbers,addresses,enterpriseData
```

If you prefer to resolve secrets at runtime, provide the secret names and add `--use-secrets-manager`.

## Deploy As Lambda

You can deploy the Lambda version in either of these ways:

* Public application page: [AWS Serverless Application Repository](https://serverlessrepo.aws.amazon.com/applications/us-east-1/889836709304/idp-scim-sync)
* From source using the workflow described in [AWS-SAM.md](AWS-SAM.md)

For most users, the public AWS Serverless Application Repository page is the simplest option. For contributors and private deployments, use the source-based SAM workflow.

## Run As A Container

Build the image locally:

```bash
make build-dist
GIT_VERSION=test make container-build
```

Run the container:

```bash
podman run --rm -it \
  -v "$PWD/.idpscim.yaml:/app/.idpscim.yaml:ro" \
  ghcr.io/slashdevops/idp-scim-sync:latest \
  idpscim --config-file .idpscim.yaml
```

## Related Documentation

* [Configuration.md](Configuration.md)
* [AWS-SAM.md](AWS-SAM.md)
* [Development.md](Development.md)
* [README.md](../README.md)
