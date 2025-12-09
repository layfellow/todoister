## todoister delete task

```sh
todoister delete task [flags] [#][PARENT/.../PROJECT] TASK
```

Delete a task from Todoist.

Use <code>#[PARENT/SUBPARENT.../]PROJECT</code> to specify the project name with optional
<code>PARENT</code> and <code>SUBPARENTS</code> (note the '<code>#</code>' character prefix and the single quotes).

Alternatively, you can use the <code>--project</code> flag to specify the project name
and omit the '<code>#</code>' prefix and the quotes.
Note that <code>PROJECT</code>, <code>PARENTS</code> and <code>SUBPARENTS</code> are case-insensitive.

You can identify a <code>TASK</code> by its partial name. If multiple tasks match,
an error is shown.

This command deletes the task and all its sub-tasks.


### Flags:

<dl>
  <dt><code>-f</code>, <code>--force</code></dt>
  <dd>skip confirmation prompt</dd>
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
# Delete task from root-level project Work:
todoister delete task '#Work' 'Complete report'

# Delete task from nested project:
todoister delete task '#Work/Reports' 'Create quarterly report'

# Delete task using partial name:
todoister delete task '#Work/Reports' 'Create q'

# Delete task using project flag:
todoister delete task -p Work/Reports 'Create monthly report'

# Delete task without confirmation:
todoister delete task -f -p Personal 'Buy groceries'
todoister rm task --force '#Work' 'Old task'
```

