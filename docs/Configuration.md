# Configuration

The `idpscim` program could be configured using three ways: `configuration file`, `command line arguments` and `environment variables`.

NOTES:

* The configuration could result in a merge of the three sources.
* The precedence order is `configuration file`, `command line arguments` and `environment variables`.

## Prerequisites

To have access to the `AWS Resources` like `AWS S3` you will need `AWS Credentials`, so, before execute the `idpscim` or `idpscimcli` program use one of the following configuration methods:

Using `AWS_ACCESS_KEY`

```bash
export AWS_ACCESS_KEY_ID="<your access key>"
export AWS_SECRET_ACCESS_KEY="<your secret key>"
export AWS_REGION="<your region>"
```

Using `profiles`

```bash
export AWS_PROFILE="slashdevops"
export AWS_REGION="<your region>"
export AWS_ROLE_SESSION_NAME="idp-scim-sync"
```

__NOTES:__

* This support profiles with and without `role` and `mfa`

## Configuration file

create a `.idpscim.yaml` file in the `$HOME/` directory or in the same path where the `idpscim` program is.

```yaml
log_level: trace
log_format: text

gws_service_account_file: /path/to/gws_service_account.json
gws_user_email: my.user@gws-email.com
gws_groups_filter:
  - 'name:AWS* email:aws*'
  - 'email:administrators*'

aws_scim_endpoint: https://scim.eu-west-1.amazonaws.com/<tenant id>/scim/v2/
aws_scim_access_token: <access token>

aws_s3_bucket_name: my-bucket
aws_s3_bucket_key: data/state.json

sync_method: groups
use_secrets_manager: false
```

then run the `idpscim` program

```bash
./idpscim
```

or create a `<any filename on whenever place>.yaml` then run the `idpscim` program with the `--config-file` option.

```bash
./idpscim --config-file <any filename on whenever place>.yaml
```

## Command line arguments

```bash
# execute the program with the following arguments
./idpscim \
  --aws-s3-bucket-name "my-bucket" \
  --aws-s3-bucket-key "data/state.json" \
  --aws-scim-access-token "<access token>" \
  --aws-scim-endpoint "https://scim.eu-west-1.amazonaws.com/<tenant id>/scim/v2/" \
  --gws-service-account-file "/path/to/gws_service_account.json" \
  --gws-user-email "my.user@gws-email.com" \
  --gws-groups-filter 'name:AWS* email:aws*' \
  --gws-groups-filter 'email:administrators*' \
  --sync-method 'groups'
  --log-level trace
```

## Environment variables

```bash
# first export the environment variables
export IDPSCIM_AWS_S3_BUCKET_NAME="my-bucket"
export IDPSCIM_AWS_S3_BUCKET_KEY="data/state.json"
export IDPSCIM_AWS_SCIM_ACCESS_TOKEN="<access token>"
export IDPSCIM_AWS_SCIM_ENDPOINT="https://scim.eu-west-1.amazonaws.com/<tenant id>/scim/v2/"
export IDPSCIM_GWS_SERVICE_ACCOUNT_FILE="/path/to/gws_service_account.json"
export IDPSCIM_GWS_USER_EMAIL="my.user@gws-email.com"
export IDPSCIM_GWS_GROUPS_FILTER='name:AWS* email:aws*','email:administrators*'
export IDPSCIM_SYNC_METHOD="groups"
export IDPSCIM_LOG_LEVEL="trace"

# then execute the program
./idpscim
```
