## todoister add task

```sh
todoister add task [flags] [#][PARENT/.../PROJECT] TASK
```

Agregar una nueva tarea a un proyecto de Todoist.

Use `#[PARENT/SUBPARENT.../]PROJECT` para especificar el nombre del proyecto
(con nombres opcionales de PARENT y SUBPARENT); note el carácter `#`.

Alternativamente, use el flag --project para especificar el nombre del proyecto,
puede omitir el carácter `#`


### Flags:

<dl>
  <dt><code>-p</code>, <code>--project</code> <code>&lt;string&gt;</code></dt>
  <dd>nombre o ruta del proyecto (p. ej., 'Work' o 'Work/Reports')</dd>
</dl>

### Flags globales:

<dl>
  <dt><code>-t</code>, <code>--token</code> <code>&lt;string&gt;</code></dt>
  <dd>usa <code>&lt;string&gt;</code> como token de la API de Todoist</dd>
</dl>

### Ejemplos

```sh
# Agregar tarea a proyecto de nivel raíz Work:
todoister add task '#Work' 'Complete report'

# Agregar tarea a proyecto Reports del proyecto Work:
todoister add task '#Work/Reports' 'Create quarterly report'

# Agregar tareas usando el flag project:
todoister add task -p Work/Reports 'Create monthly report'
todoister add task -p Personal 'Buy groceries'

# Agregar tarea a proyecto anidado usando flag:
todoister add task --project=Personal/Shopping/List 'Buy milk'
```

