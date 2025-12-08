## todoister export

```sh
todoister export [flags] [PATH]
```

Exporta todos los proyectos de Todoist como un árbol de archivos JSON o YAML.

- <code>PATH</code> es un archivo o directorio donde exportar los proyectos, por defecto <code>index.json</code>.


### Opciones:

<dl>
  <dt><code>-d</code>, <code>--depth</code> <code>&lt;int&gt;</code></dt>
  <dd>profundidad del árbol de subdirectorios a crear en el sistema de archivos al exportar
(el valor predeterminado es 0, es decir, sin subdirectorios)</dd>
  <dt><code>--json</code></dt>
  <dd>exportar en formato JSON (predeterminado)</dd>
  <dt><code>--yaml</code></dt>
  <dd>exportar en formato YAML</dd>
</dl>

### Opciones globales:

<dl>
  <dt><code>-t</code>, <code>--token</code> <code>&lt;string&gt;</code></dt>
  <dd>utilizar <code>&lt;string&gt;</code> como token de API de Todoist</dd>
</dl>

### Ejemplos

```sh
# Exportar a un único archivo index.json en el directorio actual:
todoister export

# Exportar al archivo todoist.json en el directorio personal:
todoister export ~/todoist.json

# Exportar al archivo todoist.yaml en el directorio personal:
todoister export --yaml ~/todoist.yaml

# Exportar a un directorio projects en el directorio personal, con subdirectorios hasta 3 niveles de profundidad:
todoister export --json -d 3 ~/projects
```

