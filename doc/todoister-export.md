## todoister export

```sh
todoister export [path] [flags]
```

Export all Todoist projects as a tree of JSON or YAML files.

- `path` is a file or directory where to export the projects, by default `index.json`.


### Options

<dl>
  <dt><code>-d</code>, <code>--depth</code> <code>&lt;int&gt;</code></dt>
  <dd>depth of subdirectory tree to create on the filesystem when exporting
(default is 0, i.e., no subdirectories)</dd>
  <dt><code>--json</code></dt>
  <dd>export in JSON format (default)</dd>
  <dt><code>--yaml</code></dt>
  <dd>export in YAML format</dd>
</dl>

### Global Options

<dl>
  <dt><code>-t</code>, <code>--token</code> <code>&lt;string&gt;</code></dt>
  <dd>use <code>&lt;string&gt;</code> as Todoist API token</dd>
</dl>

### Examples

```sh
# Export to a single index.json file in the current directory:
todoister export

# Export to todoist.json file in the home directory:
todoister export ~/todoist.json

# Export to todoist.yaml file in the home directory:
todoister export ~/todoist.yaml --yaml

# Export to a projects directory in the home, with subdirectories down to 3 levels deep:
todoister export ~/projects --json -d 3
```

