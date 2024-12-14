## todoister export

```sh
todoister export [path] [flags]
```

Exporte todos los proyectos de Todoist como un árbol de archivos JSON o YAML.

- `path` es un archivo o directorio donde exportar los proyectos, por defecto `index.json`.

### Opciones

<dl>
  <dt><code>-d</code>, <code>--depth</code> <code>&lt;int&gt;</code></dt>
  <dd>profundidad del árbol de subdirectorios a crear en el sistema de archivos al exportar
(por defecto es 0, es decir, sin subdirectorios)</dd>
  <dt><code>--json</code></dt>
  <dd>exportar en formato JSON (por defecto)</dd>
  <dt><code>--yaml</code></dt>
  <dd>exportar en formato YAML</dd>
</dl>

### Opciones Globales

<dl>
  <dt><code>-t</code>, <code>--token</code> <code>&lt;string&gt;</code></dt>
  <dd>utilizar <code>&lt;string&gt;</code> como token de la API de Todoist</dd>
</dl>

### Ejemplos

```sh
# Exportar a un único archivo index.json en el directorio actual:
todoister export

# Exportar a un archivo todoist.json en el directorio de inicio:
todoister export ~/todoist.json

# Exportar a un archivo todoist.yaml en el directorio de inicio:
todoister export ~/todoist.yaml --yaml

# Exportar a un directorio de proyectos en el inicio, con subdirectorios hasta 3 niveles de profundidad:
todoister export ~/projects --json -d 3
```
