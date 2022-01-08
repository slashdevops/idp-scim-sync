# idpscim

This is in charge of keeping Google Workspace Groups and Users sync with AWS Single Sing-On service using SCIM protocol, to do that in an efficient way it stores the last sync data into an `AWS S3 Bucket` and periodically it will check if there are changes.

This program could work in three ways:

1. As an [AWS Lambda function](https://aws.amazon.com/lambda/?nc1=h_ls) deployed via [AWS SAM](https://aws.amazon.com/serverless/sam/) or consumed directly from the [AWS Serverless Application Repository](https://aws.amazon.com/serverless/serverlessrepo/?nc1=h_ls)
2. As a command line tool
3. As a Docker container

To understand how to configure the program, please read the [Configuration](Configuration.md) file.

## idpscim --help

```bash
./idpscim --help

Sync your Google Workspace Groups and Users to AWS Single Sing-On using
AWS SSO SCIM API (https://docs.aws.amazon.com/singlesignon/latest/developerguide/what-is-scim.html).

Usage:
  idpscim [flags]

Flags:
  -k, --aws-s3-bucket-key string          AWS S3 Bucket key to store the state (default "state.json")
  -b, --aws-s3-bucket-name string         AWS S3 Bucket name to store the state
  -t, --aws-scim-access-token string      AWS SSO SCIM API Access Token
  -e, --aws-scim-endpoint string          AWS SSO SCIM API Endpoint
  -c, --config-file string                configuration file (default ".idpscim.yaml")
  -d, --debug                             fast way to set the log-level to debug
  -n, --disable-state                     state status [true|false]
  -q, --gws-groups-filter strings         GWS Groups query parameter, example: --gws-groups-filter 'name:Admin* email:admin*' --gws-groups-filter 'name:Power* email:power*'
  -s, --gws-service-account-file string   path to Google Workspace service account file (default "credentials.json")
  -u, --gws-user-email string             GWS user email with allowed access to the Google Workspace Service Account
  -r, --gws-users-filter strings          GWS Users query parameter, example: --gws-users-filter 'name:Admin* email:admin*' --gws-users-filter 'name:Power* email:power*'
  -h, --help                              help for idpscim
  -f, --log-format string                 set the log format (default "text")
  -l, --log-level string                  set the log level [panic|fatal|error|warn|info|debug|trace] (default "info")
  -m, --sync-method string                Sync method to use [groups] (default "groups")
  -v, --version                           version for idpscim
```

## Using the AWS Lambda function

This could be deployed using the [official AWS Serverless public repository]() or using the method explained in the [AWS SAM](docs/AWS-SAM.md) section.

## Using the command line tool

This could be used following the instructions in the main [README.md](docs/README.md) file.

## Using the Docker image

this is a __WIP__

Test and build the Docker image

```bash
make test
make container-build
```

Execute

```bash
docker run -it -v $HOME/tmp/idpscim.yaml:/app/.idpscim.yaml ghcr.io/slashdevops/idp-scim-sync-linux-arm64v8 idpscim --debug
```
