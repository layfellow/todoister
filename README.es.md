<h1 align="center">Cliente CLI para Todoist</h1>
<p align="center"><img src="icon.svg" width="80" height="78" alt="Todoister"/></p>

[README in English](README.md)

Todoister es un cliente CLI simple para [Todoist](https://todoist.com/) escrito en Go.
Úselo para gestionar rápidamente sus tareas y proyectos de Todoist en un terminal.

También ofrece una función de exportación mucho mejor que la copia de seguridad CSV estándar:
la exportación de Todoister admite JSON o YAML estructurados con profundidad configurable
para directorios anidados.

> [!NOTE]
> Consulte la Guía del Usuario en [https://todoister.parroquiano.net/](https://todoister.parroquiano.net/) para una referencia completa de los comandos implementados.

*Todoister aún está en desarrollo. Próximamente se añadirán más funciones, como la edición de tareas, la gestión de etiquetas, etc.*

## Instalación

Para Linux y macOS, use:

```sh
curl -sfL https://todoister.layfellow.net/installer.sh | sh
```

Este script descarga el binario más reciente para su plataforma y lo instala en `~/.local/bin` o `~/bin`.

Para Windows ... eh,
[no uso Windows](https://www.fsf.org/es/news/la-vida-es-mejor-juntos-cuando-evitas-windows-11),
así que no hay versiones para éste, pero el binario de Linux debería funcionar bajo
[WSL 2](https://learn.microsoft.com/en-us/windows/wsl/).

Alternativamente, si tiene Go (versión 1.24 o posterior), puede descargar, compilar e instalar
Todoister con:

```sh
go install github.com/layfellow/todoister@latest
```

## Configuración

Necesita un token de API de Todoist; inicie sesión en su cuenta de Todoist y cree uno
[aquí](https://app.todoist.com/app/settings/integrations/developer).

Luego escriba un archivo `~/.config/todoister/config.toml` o `~/.todoister.toml` y establezca el token:

```toml
token = "su-token-de-API-de-todoist"
```

El formato completo del archivo de configuración es:

```toml
# Establezca aquí el token.
token = ""

# Archivo de log para ejecución no interactiva.
# El valor predeterminado es $HOME/.cache/todoister/out.log
[log]
name = ""

# Argumento y opciones para `todoister export`.
# Los valores predeterminados son path: directorio actual, format: json, depth: 0.
[export]
path = ""
format = ""
depth = 0
```

**Alternativas de configuración**

En lugar de un archivo de configuración, puede utilizar una variable de entorno para el token:

```sh
export TODOIST_TOKEN='su-token-de-API-de-todoist'
```

O pase el token directamente a través de la línea de comandos:

```sh
todoister --token='su-token-de-API-de-todoist' comando ...
```
La opción `--token` tiene prioridad sobre la variable de entorno, que a su vez tiene precedencia sobre el archivo de configuración.

## Exportación

Es posible ejecutar `todoister export` en un cron job como una forma de crear respaldos automáticos de Todoist en un formato legible.
Puede establecer las opciones de exportación directamente en el archivo de configuración `config.toml`, para que no tenga que editar el cron tab.

```toml
[export]
path = ""
format = ""
depth = 0
```

Por ejemplo, en lugar de `todoister export --yaml -d 3 ~/projects` simplemente ejecute `todoist export` con:

```toml
[export]
path = "$HOME/projects"
format = "yaml"
depth = 3
```
Cuando se ejecuta como un cron job, `todoister export` registra su actividad en un archivo de log como se establece en:

```toml
[log]
name = "/ruta/al/archivo.log"
```

## Formato de log

Los logs siguen el formato de
[log estructurado](https://pkg.go.dev/log/slog) y se rotan automáticamente.
No se escriben logs en modo interactivo.
