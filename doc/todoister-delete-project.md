## todoister delete project

```sh
todoister delete project [flags] [PARENT/.../]NAME
```

Delete a project from Todoist.

<code>NAME</code> is the name of the project to delete.
Use <code>PARENT/NAME</code> to locate a project within a parent project.
Use <code>PARENT/SUBPARENT/NAME</code> for nested parents.
Note that <code>NAMES</code>, <code>PARENTS</code> and <code>SUBPARENTS</code> are case-insensitive.

This command deletes the project and all its descendants (subprojects and tasks).


### Flags:

<dl>
  <dt><code>-f</code>, <code>--force</code></dt>
  <dd>skip confirmation prompt</dd>
</dl>

### Global Flags:

<dl>
  <dt><code>-t</code>, <code>--token</code> <code>&lt;string&gt;</code></dt>
  <dd>use <code>&lt;string&gt;</code> as Todoist API token</dd>
</dl>

### Examples

```sh
# Delete a root-level project:
todoister delete project Shopping

# Delete a project within a parent:
todoister delete project Work/Reports

# Delete a deeply nested project:
todoister delete project Work/Projects/Q1

# Delete without confirmation:
todoister delete project -f Shopping
todoister rm project --force Work/Old
```

