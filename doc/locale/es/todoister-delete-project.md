## todoister delete project

```sh
todoister delete project [flags] [PADRE/.../]NOMBRE
```

Eliminar un proyecto de Todoist.

<code>NOMBRE</code> es el nombre del proyecto a eliminar.
Utilice <code>PADRE/NOMBRE</code> para localizar un proyecto dentro de un proyecto padre.
Utilice <code>PADRE/SUBPADRE/NOMBRE</code> para padres anidados.

Este comando elimina el proyecto y todos sus descendientes (subproyectos y tareas).


### Opciones:

<dl>
  <dt><code>-f</code>, <code>--force</code></dt>
  <dd>omitir solicitud de confirmación</dd>
</dl>

### Opciones globales:

<dl>
  <dt><code>-t</code>, <code>--token</code> <code>&lt;string&gt;</code></dt>
  <dd>utilizar <code>&lt;string&gt;</code> como token de API de Todoist</dd>
</dl>

### Ejemplos

```sh
# Eliminar un proyecto de nivel raíz:
todoister delete project Shopping

# Eliminar un proyecto dentro de un padre:
todoister delete project Work/Reports

# Eliminar un proyecto profundamente anidado:
todoister delete project Work/Projects/Q1

# Eliminar sin confirmación:
todoister delete project -f Shopping
todoister rm project --force Work/Old
```

