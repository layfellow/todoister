## todoister list

```sh
todoister list [flags] [NAME]...
```

Lista proyectos y subproyectos.

<code>NAME</code> es el nombre de uno o más proyectos de los cuales listar tareas.
Si no se especifica <code>NAME</code>, se listan todos los proyectos.
Puede especificar el nombre de un proyecto mediante su ruta completa, por ejemplo, <code>Work/Project</code>.
Los nombres no distinguen entre mayúsculas y minúsculas.


### Opciones globales:

<dl>
  <dt><code>-t</code>, <code>--token</code> <code>&lt;string&gt;</code></dt>
  <dd>utilizar <code>&lt;string&gt;</code> como token de API de Todoist</dd>
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
