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
  <dd>usa <code>&lt;string&gt;</code> como token de la API de Todoist</dd>
</dl>

### Ejemplos

```sh
# Exportar a un único archivo index.json en el directorio actual:
todoister export

# Exportar a un archivo todoist.json en el directorio home:
todoister export ~/todoist.json

# Exportar a un archivo todoist.yaml en el directorio home:
todoister export --yaml ~/todoist.yaml

# Exportar a un directorio projects en el home, con subdirectorios hasta 3 niveles de profundidad:
todoister export --json -d 3 ~/projects
```
