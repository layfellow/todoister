## todoister add task

```sh
todoister add task [flags] [PROJECT_PATH] TASK_TITLE
```

Add a new task to a Todoist project.

Usage formats:
  todoister add task '#[PARENT PROJECT/]PROJECT NAME' 'TASK TITLE'
  todoister add task -p '[PARENT PROJECT/]PROJECT NAME' 'TASK TITLE'
  todoister add task --project '[PARENT PROJECT/]PROJECT NAME' 'TASK TITLE'

Examples:
  # Add task to root-level project:
  todoister add task '#Work' 'Complete report'

  # Add task to nested project:
  todoister add task '#Work/Reports' 'Create quarterly report'

  # Add task using project flag:
  todoister add task -p 'Personal' 'Buy groceries'

  # Add task to nested project using flag:
  todoister add task --project 'Personal/Shopping' 'Buy milk'

### Options

<dl>
  <dt><code>-p</code>, <code>--project</code> <code>&lt;string&gt;</code></dt>
  <dd>project name or path (e.g., 'Work' or 'Work/Reports')</dd>
</dl>

### Global Options

<dl>
  <dt><code>-t</code>, <code>--token</code> <code>&lt;string&gt;</code></dt>
  <dd>use <code>&lt;string&gt;</code> as Todoist API token</dd>
</dl>

