# Configuration

This document describes how to configure the project binaries, especially `idpscim`, which is the program executed by Lambda and by the local sync CLI.

Use this document together with:

* [idpscim.md](idpscim.md) for the executable and flag overview
* [idpscimcli.md](idpscimcli.md) for the validation CLI commands
* [AWS-SAM.md](AWS-SAM.md) for Lambda deployment examples

## Configuration Sources

The project uses a combination of Cobra flags and Viper configuration loading.

In practice, you have these input sources:

* Built-in defaults from the Go code
* A config file such as `.idpscim.yaml`
* Environment variables prefixed with `IDPSCIM_`
* Command-line flags for the active command

Because some settings can be supplied from more than one place, the safest operating model is to choose one primary source for each deployment style:

* Local development: use a config file plus a small number of CLI flags when needed
* AWS SAM and Lambda: use environment variables and secret names from the template
* One-off testing: use CLI flags directly

## AWS Credentials

Both `idpscim` and `idpscimcli` rely on the AWS SDK default credential chain when they need AWS access.

Common approaches:

Using direct credentials:

```bash
export AWS_ACCESS_KEY_ID="<your access key>"
export AWS_SECRET_ACCESS_KEY="<your secret key>"
export AWS_REGION="<your region>"
```

Using a profile:

```bash
export AWS_PROFILE="my-profile"
export AWS_REGION="<your region>"
export AWS_ROLE_SESSION_NAME="idp-scim-sync"
```

Profiles with role assumption and MFA are supported by the AWS SDK as long as your local AWS CLI configuration is valid.

## Main `idpscim` Settings

The main sync program uses these configuration groups:

| Group | Settings |
| --- | --- |
| Logging | `log_level`, `log_format`, `debug` |
| Google Workspace | `gws_service_account_file`, `gws_user_email`, `gws_groups_filter` |
| Google Workspace secret names | `gws_service_account_file_secret_name`, `gws_user_email_secret_name` |
| AWS SCIM | `aws_scim_endpoint`, `aws_scim_access_token` |
| AWS SCIM secret names | `aws_scim_endpoint_secret_name`, `aws_scim_access_token_secret_name` |
| State repository | `aws_s3_bucket_name`, `aws_s3_bucket_key` |
| Sync behavior | `sync_method`, `sync_user_fields`, `use_secrets_manager` |

Important notes:

* `sync_method` currently supports `groups`
* `sync_user_fields` is optional; when empty, all supported optional user attributes are synced
* `use_secrets_manager=true` tells the program to resolve credential values from AWS Secrets Manager using the configured secret names
* The code default for `aws_s3_bucket_key` is `state.json`, while the AWS SAM template overrides it to `data/state.json` unless you change the template parameter

## Config File Example

For local usage, place `.idpscim.yaml` in your home directory, the current working directory, or point to a file explicitly with `--config-file`.

```yaml
log_level: debug
log_format: text

gws_service_account_file: /path/to/gws_service_account.json
gws_user_email: admin@example.com
gws_groups_filter:
  - 'name:AWS*'
  - 'email:aws-admins@example.com'

aws_scim_endpoint: https://example.awsapps.com/scim/v2/
aws_scim_access_token: <access token>

aws_s3_bucket_name: idp-scim-sync-state-123456789012-us-east-1
aws_s3_bucket_key: data/state.json

sync_method: groups
sync_user_fields:
  - phoneNumbers
  - addresses
  - enterpriseData
use_secrets_manager: false
```

Run with the default config file name:

```bash
./build/idpscim
```

Or point to a specific file:

```bash
./build/idpscim --config-file /path/to/custom.idpscim.yaml
```

## CLI Example

Use flags directly for one-off tests or automation:

```bash
./build/idpscim \
  --aws-s3-bucket-name "idp-scim-sync-state-123456789012-us-east-1" \
  --aws-s3-bucket-key "data/state.json" \
  --aws-scim-access-token "<access token>" \
  --aws-scim-endpoint "https://example.awsapps.com/scim/v2/" \
  --gws-service-account-file "/path/to/gws_service_account.json" \
  --gws-user-email "admin@example.com" \
  --gws-groups-filter 'name:AWS*' \
  --sync-method groups \
  --sync-user-fields phoneNumbers,addresses,enterpriseData \
  --log-level debug
```

## Environment Variable Example

Environment variables are especially useful for AWS SAM and Lambda deployments.

```bash
export IDPSCIM_AWS_S3_BUCKET_NAME="idp-scim-sync-state-123456789012-us-east-1"
export IDPSCIM_AWS_S3_BUCKET_KEY="data/state.json"
export IDPSCIM_AWS_SCIM_ACCESS_TOKEN="<access token>"
export IDPSCIM_AWS_SCIM_ENDPOINT="https://example.awsapps.com/scim/v2/"
export IDPSCIM_GWS_SERVICE_ACCOUNT_FILE="/path/to/gws_service_account.json"
export IDPSCIM_GWS_USER_EMAIL="admin@example.com"
export IDPSCIM_GWS_GROUPS_FILTER='name:AWS*'
export IDPSCIM_SYNC_METHOD="groups"
export IDPSCIM_SYNC_USER_FIELDS="phoneNumbers,addresses,enterpriseData"
export IDPSCIM_LOG_LEVEL="debug"

./build/idpscim
```

If you want runtime secret resolution instead of passing raw credentials:

```bash
export IDPSCIM_USE_SECRETS_MANAGER=true
export IDPSCIM_GWS_SERVICE_ACCOUNT_FILE_SECRET_NAME="IDPSCIM_GWSServiceAccountFile"
export IDPSCIM_GWS_USER_EMAIL_SECRET_NAME="IDPSCIM_GWSUserEmail"
export IDPSCIM_AWS_SCIM_ENDPOINT_SECRET_NAME="IDPSCIM_SCIMEndpoint"
export IDPSCIM_AWS_SCIM_ACCESS_TOKEN_SECRET_NAME="IDPSCIM_SCIMAccessToken"
```

## `idpscimcli` Notes

`idpscimcli` uses the same config file name and some of the same fields, but it is command-oriented:

* AWS SCIM inspection commands usually rely on `--aws-scim-endpoint` and `--aws-scim-access-token`
* Google Workspace inspection commands usually rely on `--gws-service-account-file` and `--gws-user-email`
* Output formatting is controlled with `--output-format`
* Request timeout is controlled with `--timeout`

For the exact command tree and examples, see [idpscimcli.md](idpscimcli.md).

## Related Documentation

* [idpscim.md](idpscim.md)
* [idpscimcli.md](idpscimcli.md)
* [AWS-SAM.md](AWS-SAM.md)
* [Development.md](Development.md)
