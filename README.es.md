<h1 align="center">Cliente CLI para Todoist</h1>
<p align="center"><img src="icon.svg" width="80" height="78" alt="Todoister"/></p>

[README in English](README.md)

Todoister es un cliente CLI simple para [Todoist](https://todoist.com/) escrito en Go.

**Funciones actuales:**

- `list`: listar proyectos en vista jerárquica.
- `tasks`: listar tareas dentro de proyectos.
- `add project`: crear nuevos proyectos (con color y jerarquía opcionales).
- `add task`: crear nuevas tareas en proyectos existentes con sintaxis flexible.
- `export`: exportar proyectos y tareas a JSON o YAML con estructura jerárquica.

Escribí esto porque quería una forma simple y rápida de gestionar mis tareas y proyectos de Todoist en un terminal.

Además, estaba insatisfecho con la única opción de exportación de Todoist, que es un archivo de valores separados por comas no estructurado
([el horror, el horror](https://www.oxfordreference.com/display/10.1093/acref/9780199567454.001.0001/acref-9780199567454-e-931)),
que carece del detalle que necesito. Quería algo compatible con cron jobs para respaldos no supervisados, en un formato más manejable, como JSON o YAML.

*Más características como edición de tareas, gestión de etiquetas, etc. próximamente.*

## Instalación

Para Linux y macOS, use:

```sh
$ curl -sfL https://layfellow.net/todoister/installer.sh | sh
```

Este script descarga el binario más reciente para su plataforma y lo instala en `~/.local/bin` o `~/bin`.

Para Windows ... eh,
[no uso Windows](https://www.fsf.org/es/news/la-vida-es-mejor-juntos-cuando-evitas-windows-11),
así que no hay versiones para éste, pero el binario de Linux debería funcionar bajo
[WSL 2](https://learn.microsoft.com/en-us/windows/wsl/).

Alternativamente, si tiene Go (versión 1.22 o posterior), puede descargar, compilar e instalar
Todoister con:

```sh
$ go install github.com/layfellow/todoister@latest
```

## Configuración

Necesita un token de API de Todoist; inicie sesión en su cuenta de Todoist y cree uno
[aquí](https://app.todoist.com/app/settings/integrations/developer).

Luego escriba un archivo `~/.config/todoister/config.toml` o `~/.todoister.toml` y establezca el token:

```toml
token = "su-token-de-API-de-todoist"
```

Alternativamente, establezca una variable de entorno:

```sh
$ export TODOIST_TOKEN='su-token-de-API-de-todoist'
```
O pase el token directamente a través de la línea de comandos:

```sh
$ todoister --token='su-token-de-API-de-todoist' comando ...
```
La opción `--token` tiene prioridad sobre la variable de entorno, que a su vez tiene precedencia sobre el archivo de configuración.

## Cron job

Es posible ejecutar `todoister export` en un cron job como una forma de crear respaldos automáticos de Todoist en un formato legible.
Puede establecer las opciones de exportación directamente en el archivo de configuración `config.toml`, para que no tenga que editar el cron tab.

```toml
[export]
path = ""
format = ""
depth = 0
```

Por ejemplo, en lugar de `todoister export ~/projects --yaml -d 3` simplemente ejecute `todoist export` con:

```toml
[export]
path = "$HOME/projects"
format = "yaml"
depth = 3
```
Cuando se ejecuta como un cron job, `todoister export` registra su actividad en un archivo de log como se establece en:

```toml
[log]
name = "/ruta/al/archivo/de/registro.log"
```

Consulte el archivo de ejemplo `config.toml.example` para más detalles.

Los registros siguen el formato de
[registro estructurado](https://pkg.go.dev/log/slog) y se rotan automáticamente.
No se escriben registros en modo interactivo.

## Comandos

Consulte la Guía del Usuario en [https://parroquiano.net/todoister/](https://parroquiano.net/todoister/) para una referencia completa de los comandos implementados.

### Ejemplos de Uso Básico

```sh
# Listar todos los proyectos
todoister list

# Listar tareas en un proyecto
todoister tasks NombreProyecto

# Crear un nuevo proyecto
todoister add project "Nuevo Proyecto"

# Crear un proyecto con color
todoister add project --color=blue "Cosas Importantes"

# Crear un proyecto anidado
todoister add project "Trabajo/Informes"

# Añadir una tarea a un proyecto
todoister add task -p "Trabajo" "Completar el informe"

# Añadir una tarea a un proyecto anidado
todoister add task -p "Trabajo/Informes" "Crear resumen trimestral"

# Sintaxis alternativa para añadir tareas
todoister add task "#Personal" "Comprar víveres"

# Exportar a JSON
todoister export ~/respaldo.json

# Exportar a YAML con estructura de directorios anidada
todoister export ~/respaldo --yaml -d 3
```

