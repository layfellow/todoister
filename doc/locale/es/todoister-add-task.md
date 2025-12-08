## todoister add task

```sh
todoister add task [flags] [#][PARENT/.../PROJECT] TASK
```

Añade una nueva tarea a un proyecto de Todoist.

Utilice <code>#[PARENT/SUBPARENT.../]PROJECT</code> para especificar el nombre del proyecto con
<code>PARENT</code> y <code>SUBPARENTS</code> opcionales (observe el carácter '<code>#</code>' como prefijo y las comillas simples).

Alternativamente, puede utilizar la opción <code>--project</code> para especificar el nombre del proyecto
y omitir el prefijo '<code>#</code>' y las comillas.


### Opciones:

<dl>
  <dt><code>-d</code>, <code>--date</code> <code>&lt;string&gt;</code></dt>
  <dd>fecha de vencimiento (YYYY-MM-DD, YYYY-MM-DD HH:MM, o una cadena como 'tomorrow',
consulte https://www.todoist.com/help/articles/introduction-to-dates-and-time
para obtener ayuda sobre cómo escribir fechas en lenguaje natural)</dd>
  <dt><code>-p</code>, <code>--project</code> <code>&lt;string&gt;</code></dt>
  <dd>nombre o ruta del proyecto (por ejemplo, 'Work' o 'Work/Reports')</dd>
</dl>

### Opciones globales:

<dl>
  <dt><code>-t</code>, <code>--token</code> <code>&lt;string&gt;</code></dt>
  <dd>utiliza <code>&lt;string&gt;</code> como token de API de Todoist</dd>
</dl>

### Ejemplos

```sh
# Añadir tarea al proyecto Work en el nivel raíz:
todoister add task '#Work' 'Complete report'

# Añadir tarea al proyecto Reports dentro del proyecto Work:
todoister add task '#Work/Reports' 'Create quarterly report'

# Añadir tareas utilizando la opción project:
todoister add task -p Work/Reports 'Create monthly report'
todoister add task -p Personal 'Buy groceries'

# Añadir tarea a un proyecto anidado utilizando la opción:
todoister add task --project=Personal/Shopping/List 'Buy milk'

# Añadir tarea con fecha de vencimiento:
todoister add task -p Work -d '2026-01-15' 'Submit report'
todoister add task -p Work -d 'Jan. 16' 'Submit another report'
todoister add task -p Work --date='2026-01-16' 'Submit yet another report'
todoister add task -p Work -d 'next tuesday 14:00' 'Team meeting'
todoister add task -p Personal -d 'tomorrow' 'Call dentist'
todoister add task -p Personal --date='every friday' 'Weekly review'
```

