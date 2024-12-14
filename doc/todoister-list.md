## todoister list

```sh
todoister list [project]... [flags]
```

List projects and subprojects.

`project` is the name of one or more projects to list tasks from.
If no `project` is given, all projects are listed.
You can specify a project name by its full path, e.g., `Work/Project`.
Names are case-insensitive.


### Global Options

<dl>
  <dt><code>-t</code>, <code>--token</code> <code>&lt;string&gt;</code></dt>
  <dd>use <code>&lt;string&gt;</code> as Todoist API token</dd>
</dl>

### Examples

```sh
# List all projects and subprojects:
todoister ls

# List projects Work and Life and their subprojects:
todoister ls Work Life

# List all subprojects of Project, which is a subproject of Work:
todoister ls Work/Project
```

