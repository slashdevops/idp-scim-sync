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

The tables below mirror `idpscim --help`, including defaults.

### General And Logging

| Flag | Purpose | Default |
| --- | --- | --- |
| `--config-file`, `-c` | Path to the configuration file | `.idpscim.yaml` |
| `--debug`, `-d` | Shortcut to force debug logging | `false` |
| `--log-format`, `-f` | Log output format: `text` or `json` | `text` |
| `--log-level`, `-l` | Log verbosity: `debug`, `info`, `warn`, `error`. `fatal` and `panic` are accepted as aliases for `error`. | `info` |
| `--version`, `-v` | Show version information | — |
| `--help`, `-h` | Show help | — |

Use `--log-format json` in Lambda so CloudWatch Logs Insights can query the structured fields. See
[Configuration.md](Configuration.md#log-levels) for the full log-level semantics.

### Google Workspace Input

| Flag | Purpose | Default |
| --- | --- | --- |
| `--gws-service-account-file`, `-s` | Path to the Google Workspace service account JSON. In Lambda this holds the JSON *content*, resolved from Secrets Manager, not a path. | `credentials.json` |
| `--gws-user-email`, `-u` | Delegated Google Workspace user email | — |
| `--gws-groups-filter`, `-q` | One or more filters that restrict which groups are synchronized. Repeatable. | `[""]` (all groups) |
| `--gws-service-account-file-secret-name`, `-o` | Secret name used when resolving the service account JSON from AWS Secrets Manager | `IDPSCIM_GWSServiceAccountFile` |
| `--gws-user-email-secret-name`, `-p` | Secret name used when resolving the delegated user email from AWS Secrets Manager | `IDPSCIM_GWSUserEmail` |

### AWS SCIM And State Storage

| Flag | Purpose | Default |
| --- | --- | --- |
| `--aws-scim-endpoint`, `-e` | AWS IAM Identity Center SCIM endpoint | — |
| `--aws-scim-access-token`, `-t` | AWS IAM Identity Center SCIM access token | — |
| `--aws-scim-endpoint-secret-name`, `-n` | Secret name used when resolving the SCIM endpoint from AWS Secrets Manager | `IDPSCIM_SCIMEndpoint` |
| `--aws-scim-access-token-secret-name`, `-j` | Secret name used when resolving the SCIM token from AWS Secrets Manager | `IDPSCIM_SCIMAccessToken` |
| `--aws-s3-bucket-name`, `-b` | S3 bucket used to store the sync state | — |
| `--aws-s3-bucket-key`, `-k` | S3 object key used for the sync state. The AWS SAM template overrides this to `data/state.json`. | `state.json` |
| `--use-secrets-manager`, `-g` | Load credential values from AWS Secrets Manager | `false` |

### Sync Behavior

| Flag | Purpose | Default |
| --- | --- | --- |
| `--sync-method`, `-m` | Sync strategy. The only implemented value is `groups`. | `groups` |
| `--sync-user-fields` | Optional user attributes to synchronize. Empty means **all** of them. No shorthand. | *(empty — all fields)* |

Valid `--sync-user-fields` values: `phoneNumbers`, `addresses`, `title`, `preferredLanguage`,
`locale`, `timezone`, `nickName`, `profileURL`, `userType`, `enterpriseData`. Narrowing the set also
narrows what is requested from the Google API. Changing it changes user hashes, so the next run
updates every user once.

## Example Local Run

Build the binary:

```bash
make build
```

Run the program with direct credentials and a state bucket:

```bash
./build/idpscim \
  --gws-service-account-file "$HOME/.config/idpscim/credentials.json" \
  --gws-user-email admin@example.com \
  --gws-groups-filter 'name:AWS*' \
  --aws-scim-endpoint https://example.awsapps.com/scim/v2/ \
  --aws-scim-access-token "$SCIM_ACCESS_TOKEN" \
  --aws-s3-bucket-name idp-scim-sync-state-123456789012-us-east-1 \
  --aws-s3-bucket-key data/state.json \
  --sync-method groups \
  --sync-user-fields phoneNumbers,addresses,enterpriseData
```

> [!CAUTION]
> Keep the service-account JSON **outside the repository**. The default value is `credentials.json`,
> which resolves to the current working directory — convenient, but it puts a private key in a
> directory that gets archived and shared. In production the file is not used at all: the JSON is
> resolved from AWS Secrets Manager.

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
