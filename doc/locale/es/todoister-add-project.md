## todoister add project

```sh
todoister add project [flags] [PADRE/.../]NOMBRE
```

Añade un nuevo proyecto a Todoist.

`NOMBRE` es el nombre del proyecto a crear.
Usa `PADRE/NOMBRE` para crear un proyecto dentro de un proyecto padre.
Usa `PADRE/SUBPADRE/NOMBRE` para padres anidados.


### Flags:

<dl>
  <dt><code>-c</code>, <code>--color</code> <code>&lt;string&gt;</code></dt>
  <dd>color del proyecto (berry_red, red, orange, yellow, olive_green, lime_green, green, mint_green, teal, sky_blue, light_blue, blue, grape, violet, lavender, magenta, salmon, charcoal, grey, taupe)</dd>
</dl>

### Flags globales:

<dl>
  <dt><code>-t</code>, <code>--token</code> <code>&lt;string&gt;</code></dt>
  <dd>usa <code>&lt;string&gt;</code> como token del API de Todoist</dd>
</dl>

### Ejemplos

```sh
# Añadir un proyecto de nivel raíz:
todoister add project "Shopping"

# Añadir un proyecto dentro de un padre:
todoister add project "Work/Reports"

# Añadir un proyecto profundamente anidado:
todoister add project "Work/Projects/Q1"

# Añadir un proyecto con color:
todoister add project -c blue "Personal"

# Añadir un proyecto coloreado dentro de un padre:
todoister add project --color=red "Work/Urgent"
```

