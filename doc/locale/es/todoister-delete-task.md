## todoister delete task

```sh
todoister delete task [flags] [#][PARENT/.../PROJECT] TASK
```

Elimina una tarea de Todoist.

Use <code>#[PARENT/SUBPARENT.../]PROJECT</code> para especificar el nombre del proyecto con
<code>PARENT</code> y <code>SUBPARENTS</code> opcionales (note el carácter '<code>#</code>' como prefijo y las comillas simples).

Alternativamente, puede usar la opción <code>--project</code> para especificar el nombre del proyecto
y omitir el prefijo '<code>#</code>' y las comillas.
Observe que <code>PROJECT</code>, <code>PARENTS</code> y <code>SUBPARENTS</code>
no distinguen entre mayúsculas y minúsculas.

Se puede identificar un <code>TASK</code> con parte del nombre. Si múltiples TASKS
coinciden con el nombre parcial, se muestra un error.

Este comando elimina la tarea y todas sus subtareas.


### Opciones:

<dl>
  <dt><code>-f</code>, <code>--force</code></dt>
  <dd>omitir solicitud de confirmación</dd>
  <dt><code>-p</code>, <code>--project</code> <code>&lt;string&gt;</code></dt>
  <dd>nombre o ruta del proyecto (por ejemplo, 'Work' o 'Work/Reports')</dd>
</dl>

### Opciones globales:

<dl>
  <dt><code>-t</code>, <code>--token</code> <code>&lt;string&gt;</code></dt>
  <dd>usar <code>&lt;string&gt;</code> como token de la API de Todoist</dd>
</dl>

### Ejemplos

```sh
# Eliminar tarea del proyecto raíz Work:
todoister delete task '#Work' 'Complete report'

# Eliminar tarea de un proyecto anidado:
todoister delete task '#Work/Reports' 'Create quarterly report'

# Eliminar una tarea con parte del nombre:
todoister delete task '#Work/Reports' 'Create q'

# Eliminar tarea usando la opción de proyecto:
todoister delete task -p Work/Reports 'Create monthly report'

# Eliminar tarea sin confirmación:
todoister delete task -f -p Personal 'Buy groceries'
todoister rm task --force '#Work' 'Old task'
```

