## todoister check

```sh
todoister check [flags] [#][PARENT/.../PROJECT] TASK
```

Mark a <code>TASK</code> in a <code>PROJECT</code> as completed.

Use <code>#[PARENT/SUBPARENT.../]PROJECT</code> to specify the project name with optional
<code>PARENT</code> and <code>SUBPARENTS</code> (note the <code>'#'</code> character prefix and the single quotes).

Alternatively, you can use the <code>--project</code> flag to specify the project name
and omit the <code>'#'</code> prefix and the quotes.

The command matches tasks by prefix (case-insensitive). If multiple tasks
match, an error is shown with a list of matching tasks.

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
  # Check a task in a root project
  todoister check '#Work' 'Write report'
  todoister check -p Work 'Write report'

  # Check a task in a nested project
  todoister check '#Work/Reports' 'Q4 summary'
  todoister check -p Work/Reports 'Q4 summary'
```

