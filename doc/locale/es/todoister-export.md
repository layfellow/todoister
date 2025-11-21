## todoister export

```sh
todoister export [flags] [PATH]
```

Exporta todos los proyectos de Todoist como un árbol de archivos JSON o YAML.

- `PATH` es un archivo o directorio donde exportar los proyectos, por defecto `index.json`.


### Flags:

<dl>
  <dt><code>-d</code>, <code>--depth</code> <code>&lt;int&gt;</code></dt>
  <dd>profundidad del árbol de subdirectorios a crear en el sistema de archivos al exportar
(por defecto es 0, es decir, sin subdirectorios)</dd>
  <dt><code>--json</code></dt>
  <dd>exportar en formato JSON (por defecto)</dd>
  <dt><code>--yaml</code></dt>
  <dd>exportar en formato YAML</dd>
</dl>

### Flags globales:

<dl>
  <dt><code>-t</code>, <code>--token</code> <code>&lt;string&gt;</code></dt>
  <dd>usar <code>&lt;string&gt;</code> como token de API de Todoist</dd>
</dl>

### Ejemplos

```sh
# Export to a single index.json file in the current directory:
todoister export

# Export to todoist.json file in the home directory:
todoister export ~/todoist.json

# Export to todoist.yaml file in the home directory:
todoister export --yaml ~/todoist.yaml

# Export to a projects directory in the home, with subdirectories down to 3 levels deep:
todoister export --json -d 3 ~/projects
```
