## todoister add task

```sh
todoister add task [flags] [#][PADRE/.../PROYECTO] TAREA
```

Añade una nueva tarea a un proyecto de Todoist.

Usa `#[PADRE/SUBPADRE.../]PROYECTO` para especificar el nombre del proyecto con
`PADRE` y `SUBPADRES` opcionales (observe el prefijo '`#`' y las comillas simples).

Alternativamente, puede usar el flag `--project` para especificar el nombre del proyecto
y omitir el prefijo '`#`' y las comillas.


### Flags:

<dl>
  <dt><code>-p</code>, <code>--project</code> <code>&lt;string&gt;</code></dt>
  <dd>nombre o ruta del proyecto (por ejemplo, 'Trabajo' o 'Trabajo/Informes')</dd>
</dl>

### Flags globales:

<dl>
  <dt><code>-t</code>, <code>--token</code> <code>&lt;string&gt;</code></dt>
  <dd>usa <code>&lt;string&gt;</code> como token de la API de Todoist</dd>
</dl>

### Ejemplos

```sh
# Añadir tarea al proyecto de nivel raíz Trabajo:
todoister add task '#Trabajo' 'Completar informe'

# Añadir tarea al proyecto Informes del proyecto Trabajo:
todoister add task '#Trabajo/Informes' 'Crear informe trimestral'

# Añadir tareas usando el flag --project:
todoister add task -p Trabajo/Informes 'Crear informe mensual'
todoister add task -p Personal 'Comprar comestibles'

# Añadir tarea a proyecto anidado usando el flag --project:
todoister add task --project=Personal/Compras/Lista 'Comprar leche'
```
