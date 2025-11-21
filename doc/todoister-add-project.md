## todoister add project

```sh
todoister add project [flags] [PARENT/.../]NAME
```

Add a new project to Todoist.

`NAME` is the name of the project to create.
Use `PARENT/NAME` to create a project within a parent project.
Use `PARENT/SUBPARENT/NAME` for nested parents.


### Flags:

<dl>
  <dt><code>-c</code>, <code>--color</code> <code>&lt;string&gt;</code></dt>
  <dd>project color (berry_red, red, orange, yellow, olive_green, lime_green, green, mint_green, teal, sky_blue, light_blue, blue, grape, violet, lavender, magenta, salmon, charcoal, grey, taupe)</dd>
</dl>

### Global Flags:

<dl>
  <dt><code>-t</code>, <code>--token</code> <code>&lt;string&gt;</code></dt>
  <dd>use <code>&lt;string&gt;</code> as Todoist API token</dd>
</dl>

### Examples

```sh
# Add a root-level project:
todoister add project "Shopping"

# Add a project within a parent:
todoister add project "Work/Reports"

# Add a deeply nested project:
todoister add project "Work/Projects/Q1"

# Add a project with a color:
todoister add project -c blue "Personal"

# Add a colored project within a parent:
todoister add project --color=red "Work/Urgent"
```

