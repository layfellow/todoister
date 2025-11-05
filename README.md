<h1 align="center">CLI Client for Todoist</h1>
<p align="center"><img src="icon.svg" width="80" height="78" alt="Todoister"/></p>

[LÉAME en español](README.es.md)

Todoister is a simple [Todoist](https://todoist.com/) CLI client written in Go.

**This is an early release, with very reduced functionality.** Currenty implemented:

- `list`: list projects.
- `tasks`: list project tasks.
- `export`: export projects and tasks to JSON or YAML files.

I wrote this because I wanted a simple and quick way to check my Todoist tasks and projects
when working in the terminal.

Also, I was dissatisfied with the only export option of Todoist being unstructured
comma-separated values
([the horror, the horror](https://www.oxfordreference.com/display/10.1093/acref/9780199567454.001.0001/acref-9780199567454-e-931)),
which lack the detail I need. I wanted something cron-job-friendly for unattended
backups, in a more manageable format, like JSON or YAML.

*More features like task management, project creation, tag management, etc. coming soon.*

## Installation

For Linux and macOS, use:

```sh
$ curl -sfL https://layfellow.net/todoister/installer.sh | sh
```

This script fetches the latest binary for your platform and installs it in `~/.local/bin` or
`~/bin`.

For Windows ... huh,
[I don’t use Windows](https://www.fsf.org/news/lifes-better-together-when-you-avoid-windows-11),
so there are no releases for it, but the Linux binary should work under
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


## Cron job

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


## Commands

See the User’s Guide at [https://layfellow.net/todoister/](https://layfellow.net/todoister/) for a complete reference of the implemented commands.
