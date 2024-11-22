# Minimal Todoist CLI client

![Todoister](icon.png)

[LÉAME en español](README.es.md)

Todoister is a simple [Todoist](https://todoist.com/) CLI client written in Go.

**This is an early release, with very reduced functionality.** For now, it can only export
projects and tasks into a hierarchichal tree of JSON or YAML files.

I wrote it because I was dissatisfied with the only export option of Todoist being unstructured
comma-separated values
([the horror, the horror](https://www.oxfordreference.com/display/10.1093/acref/9780199567454.001.0001/acref-9780199567454-e-931)),
which lack the detail I need. Also, I wanted something cron-job-friendly for unattended
backups, in a more manageable format, like JSON or YAML.

More features like task management and project creation are on the horizon.

## Installation

For Linux and macOS, use:

```sh
$ curl -sfL https://layfellow.net/todoister/installer.sh | sh
```

This script fetches the latest binary for your platform and installs it in `~/.local/bin` or
`~/bin`.

For Windows ... huh, I don’t have access to a Windows system to develop and test, so there are
no releases for it, but the Linux binary should work under
[WSL 2](https://learn.microsoft.com/en-us/windows/wsl/).

Alternatively, if you have Go (version 1.22 or later), you can download, compile and install
Todoister with:

```sh
$ go install github.com/layfellow/todoister@latest
```

## Configuration

You need a Todoist API token; log in to your Todoist account and create one
[here](https://app.todoist.com/app/settings/integrations/developer).

Then write a `~/.config/todoister/config.toml` or  `~/.todoister.toml` file and set the token:

```toml
token = "your-todoist-API-token"
```

Alternatively, set an environment variable:

```sh
$ export TODOIST_TOKEN='your-todoist-API-token'
```
Or pass the token directly via the command line:

```sh
$ todoister --token='your-todoist-API-token' command ...
```
The `--token` option takes precedence over the environment variable, which in turn overrides the
configuration file.

## Usage

```sh
todoister [OPTIONS] COMMAND
```

`OPTIONS` for all commands:

- `-h`, `--help` Show help message and exit
- `-v`, `--version` Show version and exit
- `-t`, `--token` Override the Todoist API token in the configuration file or environment variable

## COMMANDS

### `help`

```sh
todoister help [COMMAND]
```

Show help message for `COMMAND` or general help if no `COMMAND` is provided.

### `export`

```sh
todoister export [PATH] [OPTIONS]
```

Export all Todoist projects to `PATH` file or directory (default is `index.json`
in the current directory).

**OPTIONS**

- `--json` Use JSON (default)
- `--yaml` Use YAML 
- `-d N`, `--depth=N`  Create directories up to `N` levels deep, writing each subproject to a
    separate file (default is 0, i.e.,no subdirectories)
 
**Examples**

Export to a single `index.json` file in the current directory:

```sh
$ todoister export
```

Export to `todoist.json` file in the home directory:

```sh
$ todoister export ~/todoist.json
```

Export to `todoist.yaml` file in the home directory:

```sh
$ todoister export ~/todoist.yaml --yaml
```

Export to a `projects` directory in the home, with subdirectories up to 3 levels
deep:

```sh
$ todoister export ~/projects --json -d 3
```

## Running `export` as a cron job

You can run `todoister export` in a cron job as a way create automatic Todoist backups in a
sane format. You can set the export options in the configuration file, so you don’t have
to edit the cron tab.

```toml
[export]
path = ""
format = ""
depth = 0
```

For instance, instead of `todoister export ~/projects --yaml -d 3` just run `todoist export`
with:

```toml
[export]
path = "$HOME/projects"
format = "yaml"
depth = 3
```
When running as a cron job, `todoister export` logs its activity to a log file as set in:

```toml
[log]
name = "/path/to/log/file.log"
```

Check the provided `config.toml.example`.

Logs follow the
[structured logging](https://pkg.go.dev/log/slog) format and are auto-rotated.
No logs are written in interactive mode.

## For developers

Todoister is written in Go (minimum version 1.22). It uses the
[Cobra framework](https://cobra.dev/)
for the CLI. Commands reside in `cmd`, utilities in `util`.

I wrote a Makefile to assist with routine tasks.

To update dependencies and update `go.mod` and `go.sum`:

    $ make dependencies

To run `golangci-lint` (requires [golangci-lint](https://golangci-lint.run/)):

    $ make lint

To build the binary for your platform:

    $ make build

To install the binary in your default path:

    $ make install

To create a new GitHub Release using the latest tag (requires [GitHub CLI](https://cli.github.com/)):

    $ make releases

Pull requests are welcome.
