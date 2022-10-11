# Minimal Todoist CLI client

## Installation

TODO

## Configuration

Write a `.todo.yaml` (*notice the dot `.`*) file in your `$HOME` directory with the following:

```yaml
TOKEN: <Your Todoist API token>
```

or, alternatively, set an environment variable:

```sh
$ export TODOIST_TOKEN=<Your Todoist API token>
```

The exported environment variable overrides the YAML file.

## Commands

**todo backup [DEST]**

Download the most recent Todoist backup to optional directory **`DEST`**. By default, use the current directory.

*Examples*

```sh
$ todo backup
$ todo backup ~/Downloads
```
