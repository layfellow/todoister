## todoister list

```sh
todoister list [flags] [NAME]...
```

Listar proyectos y subproyectos.

`NAME` es el nombre de uno o más proyectos de los cuales listar tareas.
Si no se proporciona `NAME`, se listan todos los proyectos.
Puede especificar el nombre de un proyecto por su ruta completa, por ejemplo, `Work/Project`.
Los nombres no distinguen entre mayúsculas y minúsculas.


### Flags globales:

<dl>
  <dt><code>-t</code>, <code>--token</code> <code>&lt;string&gt;</code></dt>
  <dd>usa <code>&lt;string&gt;</code> como token de la API de Todoist</dd>
</dl>

### Ejemplos

```sh
# Listar todos los proyectos y subproyectos:
todoister ls

# Listar los proyectos Work y Life y sus subproyectos:
todoister ls Work Life

# Listar todos los subproyectos de Project, que es un subproyecto de Work:
todoister ls Work/Project
```
