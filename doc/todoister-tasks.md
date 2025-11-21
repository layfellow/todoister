## todoister tasks

```sh
todoister tasks [flags] NAME...
```

List project tasks.

`NAME` is the name of one or more projects to list tasks from.
You can specify a project name by its full path, e.g., `Work/Project`.
Names are case-insensitive.


### Global Flags:

<dl>
  <dt><code>-t</code>, <code>--token</code> <code>&lt;string&gt;</code></dt>
  <dd>use <code>&lt;string&gt;</code> as Todoist API token</dd>
</dl>

### Examples

```sh
# List tasks for project Life:
todoister tasks Life

# List tasks for subproject Project of project Work:
todoister tasks Work/Project

# List tasks for both projects:
todoister tasks Life Work/Project
```

