## todoister add task

```sh
todoister add task [flags] [#][PARENT/.../PROJECT] TASK
```

Add a new task to a Todoist project.

Use `#[PARENT/SUBPARENT.../]PROJECT` to specify the project name with optional
`PARENT` and `SUBPARENTS` (note the '`#`' character prefix and the single quotes).

Alternatively, you can use the --project flag to specify the project name
and omit the '`#`' character prefix and the quotes.


### Flags:

<dl>
  <dt><code>-p</code>, <code>--project</code> <code>&lt;string&gt;</code></dt>
  <dd>project name or path (e.g., 'Work' or 'Work/Reports')</dd>
</dl>

### Global Flags:

<dl>
  <dt><code>-t</code>, <code>--token</code> <code>&lt;string&gt;</code></dt>
  <dd>use <code>&lt;string&gt;</code> as Todoist API token</dd>
</dl>

### Examples

```sh
# Add task to root-level project Work:
todoister add task '#Work' 'Complete report'

# Add task to project Reports of project Work:
todoister add task '#Work/Reports' 'Create quarterly report'

# Add tasks using project flag:
todoister add task -p Work/Reports 'Create monthly report'
todoister add task -p Personal 'Buy groceries'

# Add task to nested project using flag:
todoister add task --project=Personal/Shopping/List 'Buy milk'
```

