## todoister tasks

```sh
todoister tasks [flags] NOMBRE...
```

Listar las tareas de un proyecto.

<code>NOMBRE</code> es el nombre de uno o más proyectos de los cuales listar las tareas.
Puede especificar el nombre de un proyecto mediante su ruta completa, por ejemplo, <code>Trabajo/Proyecto</code>.
Los nombres no distinguen entre mayúsculas y minúsculas.


### Opciones globales:

<dl>
  <dt><code>-t</code>, <code>--token</code> <code>&lt;string&gt;</code></dt>
  <dd>usar <code>&lt;string&gt;</code> como token de la API de Todoist</dd>
</dl>

### Ejemplos

```sh
# Listar las tareas del proyecto Vida:
todoister tasks Vida

# Listar las tareas del subproyecto Proyecto del proyecto Trabajo:
todoister tasks Trabajo/Proyecto

# Listar las tareas de ambos proyectos:
todoister tasks Vida Trabajo/Proyecto
```

