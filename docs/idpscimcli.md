# idpscimcli

This is is a `command-line tool` to check and validate some functionalities implemented in idpscim.

The Idea of this tools is check the configuration you will implement in the [idpscim](idpscim.md) program, like the `filter` used for `Google Workspace Groups` and `Users` and the inventory of users and goups in both sides.

## idpscimcli --help

```bash
./idpscimcli --help

This is a Command-Line Interfaced (CLI) to help you validate and check your source and target Single Sing-On endpoints.
Check your AWS Single Sign-On (SSO) / Google Workspace Groups users and groups and validate your filters over Google Workspace users and groups.

Usage:
  idpscimcli [command]

Available Commands:
  aws         AWS SSO SCIM commands
  completion  Generate the autocompletion script for the specified shell
  gws         Google Workspace commands
  help        Help about any command

Flags:
  -c, --config-file string     configuration file (default ".idpscim.yaml")
  -d, --debug                  enable log debug level
  -h, --help                   help for idpscimcli
  -f, --log-format string      set the log format (default "text")
  -l, --log-level string       set the log level (default "info")
      --output-format string   output format (json|yaml) (default "json")
      --timeout duration       requests timeout (default 10s)
  -v, --version                version for idpscimcli

Use "idpscimcli [command] --help" for more information about a command.
```

## Building the project

To build the project in local, you will need to have installed and configured at least the following:

1. [git](https://git-scm.com/)
2. [Go](https://go.dev/learn/)
3. [make](https://www.gnu.org/software/make/)

Then you will need to clone the repository in your local machine, and execute the following commands:

* Compile for your Operating System:

```bash
make
```

then the binaries are in `build/` folder.

* `Cross-compiling` the project for `Windows`, `MacOS` and `Linux` (default)

```bash
make clean
make test # optional
make build-dist
```

then the binaries are in `dist/` folder.

* Others Operating Systems, see the list of supported platforms in the [syslist.go](https://github.com/golang/go/blob/master/src/go/build/syslist.go)

```bash
make clean
GO_OS=<something from goosList in syslist.go> GO_ARCH=<something from goarchList in syslist.go> make test # optional
GO_OS=<something from goosList in syslist.go> GO_ARCH=<something from goarchList in syslist.go> make build-dist
```

then the binaries are in `dist/` folder.

* Execute

```bash
./idpscimcli --help
```
