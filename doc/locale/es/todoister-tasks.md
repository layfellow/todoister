## todoister tasks

```sh
todoister tasks project... [flags]
```

Listar tareas del proyecto.

`project` es el nombre de uno o más proyectos de los cuales listar tareas.
Puede especificar un nombre de proyecto por su ruta completa, por ejemplo, `Work/Project`.
Los nombres no distinguen entre mayúsculas y minúsculas.


### Opciones Globales

<dl>
  <dt><code>-t</code>, <code>--token</code> <code>&lt;string&gt;</code></dt>
  <dd>usar <code>&lt;string&gt;</code> como token de la API de Todoist</dd>
</dl>

### Ejemplos

```sh
# Listar tareas para el proyecto Life:
todoister tasks Life

# Listar tareas para el subproyecto Project del proyecto Work:
todoister tasks Work/Project
```
