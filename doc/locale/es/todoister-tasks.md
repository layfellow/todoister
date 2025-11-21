## todoister tasks

```sh
todoister tasks [flags] NAME...
```

Lista las tareas de un proyecto.

`NAME` es el nombre de uno o más proyectos de los cuales listar tareas.
Puede especificar el nombre de un proyecto por su ruta completa, ej., `Work/Project`.
Los nombres no distinguen entre mayúsculas y minúsculas.


### Flags globales:

<dl>
  <dt><code>-t</code>, <code>--token</code> <code>&lt;string&gt;</code></dt>
  <dd>usa <code>&lt;string&gt;</code> como token de la API de Todoist</dd>
</dl>

### Ejemplos

```sh
# Listar tareas para el proyecto Life:
todoister tasks Life

# Listar tareas para el subproyecto Project del proyecto Work:
todoister tasks Work/Project

# Listar tareas para ambos proyectos:
todoister tasks Life Work/Project
```
