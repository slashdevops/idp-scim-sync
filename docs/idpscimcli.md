# idpscimcli

`idpscimcli` is the helper command-line program used to validate and inspect the same systems that `idpscim` synchronizes.

Use it when you want to answer questions such as:

- Can I reach the AWS IAM Identity Center SCIM endpoint?
- Which groups or users currently exist in AWS SCIM?
- Which Google Workspace groups match my filters?
- Which members are inside those groups?

## Source And Build Output

| Item | Location |
| --- | --- |
| Source entry point | [cmd/idpscimcli/main.go](../cmd/idpscimcli/main.go) |
| Local binary | `build/idpscimcli` |

## Command Tree

The current command surface is:

```text
idpscimcli
|- aws
|  |- service config
|  |- groups list
|  `- users list
`- gws
   |- groups list
   |- groups members list
   `- users list
```

## Root Flags

These flags apply to the whole CLI:

| Flag | Purpose |
| --- | --- |
| `--config-file`, `-c` | Path to the configuration file |
| `--debug`, `-d` | Enable debug logging |
| `--log-format`, `-f` | Log output format |
| `--log-level`, `-l` | Log verbosity |
| `--output-format` | Output format: `json` or `yaml` |
| `--timeout` | Request timeout for API calls |
| `--version`, `-v` | Show version information |

## AWS Commands

Use the `aws` command group to inspect the AWS IAM Identity Center SCIM API.

### Shared AWS Flags

| Flag | Purpose |
| --- | --- |
| `--aws-scim-endpoint`, `-e` | AWS IAM Identity Center SCIM endpoint |
| `--aws-scim-access-token`, `-t` | AWS IAM Identity Center SCIM access token |

### Available AWS Subcommands

| Command | Purpose |
| --- | --- |
| `idpscimcli aws service config` | Show SCIM service provider configuration |
| `idpscimcli aws groups list` | List groups from the AWS SCIM API |
| `idpscimcli aws users list` | List users from the AWS SCIM API |

### AWS Examples

Show SCIM service configuration:

```bash
./build/idpscimcli aws service config \
  --aws-scim-endpoint https://example.awsapps.com/scim/v2/ \
  --aws-scim-access-token "$SCIM_ACCESS_TOKEN"
```

List groups with a SCIM filter:

```bash
./build/idpscimcli aws groups list \
  --aws-scim-endpoint https://example.awsapps.com/scim/v2/ \
  --aws-scim-access-token "$SCIM_ACCESS_TOKEN" \
  --filter 'displayName eq "Engineering"'
```

List users:

```bash
./build/idpscimcli aws users list \
  --aws-scim-endpoint https://example.awsapps.com/scim/v2/ \
  --aws-scim-access-token "$SCIM_ACCESS_TOKEN"
```

## Google Workspace Commands

Use the `gws` command group to inspect Google Workspace objects with the same credentials model used by the main sync program.

### Shared Google Workspace Flags

| Flag | Purpose |
| --- | --- |
| `--gws-service-account-file`, `-s` | Path to the Google Workspace service account JSON |
| `--gws-user-email`, `-u` | Delegated Google Workspace user email |

### Available Google Workspace Subcommands

| Command | Purpose |
| --- | --- |
| `idpscimcli gws groups list` | List groups that match the provided group filters |
| `idpscimcli gws groups members list` | List the members of the groups that match the filters |
| `idpscimcli gws users list` | List users that match the provided user filters |

### Google Workspace Filter Flags

| Command | Flag |
| --- | --- |
| `gws groups list` | `--gws-groups-filter`, `-q` |
| `gws groups members list` | `--gws-groups-filter`, `-q` |
| `gws users list` | `--gws-users-filter`, `-r` |

### Google Workspace Examples

List groups:

```bash
./build/idpscimcli gws groups list \
  --gws-service-account-file credentials.json \
  --gws-user-email admin@example.com \
  --gws-groups-filter 'name=AWS*'
```

List members of matching groups:

```bash
./build/idpscimcli gws groups members list \
  --gws-service-account-file credentials.json \
  --gws-user-email admin@example.com \
  --gws-groups-filter 'email=aws-admins@example.com'
```

List users:

```bash
./build/idpscimcli gws users list \
  --gws-service-account-file credentials.json \
  --gws-user-email admin@example.com \
  --gws-users-filter 'email=alice@example.com'
```

Return YAML instead of JSON:

```bash
./build/idpscimcli gws groups list \
  --gws-service-account-file credentials.json \
  --gws-user-email admin@example.com \
  --gws-groups-filter 'name=AWS*' \
  --output-format yaml
```

## Build And Run

Build the binary locally:

```bash
make build
./build/idpscimcli --help
```

Cross-compile for distribution:

```bash
make build-dist
```

## Run From The Container Image

Build the image locally:

```bash
make build-dist
GIT_VERSION=test make container-build
```

Run the CLI from the image:

```bash
podman run --rm -it \
  -v "$PWD/.idpscim.yaml:/app/.idpscim.yaml:ro" \
  ghcr.io/slashdevops/idp-scim-sync:latest \
  idpscimcli --config-file .idpscim.yaml
```

## Related Documentation

- [idpscim.md](idpscim.md)
- [Configuration.md](Configuration.md)
- [Development.md](Development.md)
