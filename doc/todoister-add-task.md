## todoister add task

```sh
todoister add task [flags] [#][PARENT/.../PROJECT] TASK
```

Add a new task to a Todoist project.

Use <code>#[PARENT/SUBPARENT.../]PROJECT</code> to specify the project name with optional
<code>PARENT</code> and <code>SUBPARENTS</code> (note the '<code>#</code>' character prefix and the single quotes).

Alternatively, you can use the <code>--project</code> flag to specify the project name
and omit the '<code>#</code>' prefix and the quotes.


### Flags:

<dl>
  <dt><code>-d</code>, <code>--date</code> <code>&lt;string&gt;</code></dt>
  <dd>due date (YYYY-MM-DD, YYYY-MM-DD HH:MM, or a string like 'tomorrow',
see https://www.todoist.com/help/articles/introduction-to-dates-and-time
for help on how to write natural language dates )</dd>
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

# Add task with due date:
todoister add task -p Work -d '2026-01-15' 'Submit report'
todoister add task -p Work -d 'Jan. 16' 'Submit another report'
todoister add task -p Work --date='2026-01-16' 'Submit yet another report'
todoister add task -p Work -d 'next tuesday 14:00' 'Team meeting'
todoister add task -p Personal -d 'tomorrow' 'Call dentist'
todoister add task -p Personal --date='every friday' 'Weekly review'
```

