## todoister check

```sh
todoister check [opciones] [#][PADRE/.../PROYECTO] TAREA
```

Marca una <code>TAREA</code> en un <code>PROYECTO</code> como completada.

Use <code>#[PADRE/SUBPADRE.../]PROYECTO</code> para especificar el nombre del proyecto con
<code>PADRE</code> y <code>SUBPADRES</code> opcionales (note el carácter <code>'#'</code> como prefijo y las comillas simples).

Alternativamente, puede usar la opción <code>--project</code> para especificar el nombre del
proyecto y omitir el prefijo <code>'#'</code> y las comillas.

El comando busca coincidencias de tareas por prefijo (sin distinguir mayúsculas
de minúsculas). Si varias tareas coinciden, se muestra un error con una lista de
las tareas coincidentes.

### Opciones:

<dl>
  <dt><code>-p</code>, <code>--project</code> <code>&lt;string&gt;</code></dt>
  <dd>nombre o ruta del proyecto (ej., 'Trabajo' o 'Trabajo/Informes')</dd>
</dl>

### Opciones globales:

<dl>
  <dt><code>-t</code>, <code>--token</code> <code>&lt;string&gt;</code></dt>
  <dd>usar <code>&lt;string&gt;</code> como token de la API de Todoist</dd>
</dl>

### Ejemplos

```sh
  # Completar una tarea en un proyecto raíz
  todoister check '#Work' 'Write report'
  todoister check -p Work 'Write report'

  # Completar una tarea en un proyecto anidado
  todoister check '#Work/Reports' 'Q4 summary'
  todoister check -p Work/Reports 'Q4 summary'
```

