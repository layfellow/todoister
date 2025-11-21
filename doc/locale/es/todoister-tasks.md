## todoister tasks

```sh
todoister tasks [flags] NAME...
```

Lista las tareas de un proyecto.

`NAME` es el nombre de uno o más proyectos de los cuales listar tareas.
Puede especificar el nombre de un proyecto mediante su ruta completa, por ejemplo, `Trabajo/Proyecto`.
Los nombres no distinguen entre mayúsculas y minúsculas.


### Flags globales:

<dl>
  <dt><code>-t</code>, <code>--token</code> <code>&lt;string&gt;</code></dt>
  <dd>usa <code>&lt;string&gt;</code> como token de la API de Todoist</dd>
</dl>

### Ejemplos

```sh
# Lista las tareas del proyecto Vida:
todoister tasks Vida

# Lista las tareas del subproyecto Proyecto del proyecto Trabajo:
todoister tasks Trabajo/Proyecto

# Lista las tareas de ambos proyectos:
todoister tasks Vida Trabajo/Proyecto
```
