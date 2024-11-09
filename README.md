# Minimal Todoist CLI client

## Installation

Minimum Go release required: **1.22.2**


**TODO**

## Configuration

Write a `~/.config/todoister/config.json` or  `~/.todoister.json` file and set:

```json
{
  "token": "<Your Todoist API token>"
}
```

or set an environment variable:

```sh
$ export TODOIST_TOKEN='<Your Todoist API token>'
```

or provide the token as a command line argument:

```sh
$ todoister --token='<Your Todoist API token>' COMMAND
```

The command line `--token` option overrides the environment variable, which in turn overrides the configuration file.

## Commands

**todoister download [DEST]**

Download all Todoist projects to directory **`DEST`**. By default, use the current directory.

*Examples*

```sh
$ todoister download 
$ todoister download ~/projects
```
