# Cliente CLI Minimalista para Todoist

![Todoister](icon.png)

[README in English](README.md)

Todoister es un cliente CLI simple para [Todoist](https://todoist.com/) escrito en Go.

**Esta es una versión temprana, con funcionalidad muy reducida.** Por ahora, solo
puede exportar proyectos y tareas a un árbol jerárquico de archivos JSON o YAML.

Lo desarrollé porque la única opción de exportación de Todoist es hacia un archivo de valores
separados por comas no estructurado
([el horror, el horror](https://www.oxfordreference.com/display/10.1093/acref/9780199567454.001.0001/acref-9780199567454-e-931)), que carece del detalle necesario. Además, buscaba algo compatible con tareas de cron
para respaldar automáticamente mis proyectos en un formato más manejable, como JSON o YAML.

Pronto espero agregar más características, como la gestión de tareas y la creación de proyectos.

## Instalación

Para Linux y macOS, use:

```sh
$ curl -sfL https://parroquiano.net/todoister/installer.sh | sh
```

Este script descarga el binario más reciente para su plataforma y lo instala en `~/.local/bin`
o `~/bin`.

Para Windows... eh, no tengo acceso a un sistema Windows para desarrollar y probar, así que no
hay versiones para esta plataforma, pero el binario de Linux debería funcionar en
[WSL 2](https://learn.microsoft.com/en-us/windows/wsl/).

Alternativamente, si tiene Go (versión 1.22 o posterior), se puede descargar, compilar e instalar
Todoister con:

```sh
$ go install github.com/layfellow/todoister@latest
```

## Configuración

Necesita un token de API de Todoist; inicie sesión en su cuenta de Todoist y cree uno
[aquí](https://app.todoist.com/app/settings/integrations/developer).

Luego, cree un archivo `~/.config/todoister/config.toml` o `~/.todoister.toml` y copie el token:

```toml
token = "su-token-de-API-de-todoist"
```

Alternativamente, use una variable de entorno:

```sh
$ export TODOIST_TOKEN='su-token-de-API-de-todoist'
```
O pase el token directamente vía la línea de comandos:

```sh
$ todoister --token='su-token-de-API-de-todoist' comando ...
```
La opción `--token` tiene prioridad sobre la variable de entorno, que a su vez tiene prioridad
sobre el archivo de configuración.

## Uso

```sh
todoister [OPCIONES] COMANDO
```

`OPCIONES` para todos los comandos:

- `-h`, `--help` Muestra el mensaje de ayuda y sale
- `-v`, `--version` Muestra la versión y sale
- `-t`, `--token` Usa este token de API de Todoist en lugar de usar el del archivo de configuración
   o variable de entorno

## COMANDOS

### `help`

```sh
todoister help [COMANDO]
```

Muestra el mensaje de ayuda para `COMANDO` o ayuda general si no se proporciona `COMANDO`.

### `export`

```sh
todoister export [RUTA] [OPCIONES]
```

Exporta todos los proyectos de Todoist al archivo o directorio `RUTA` (por defecto es `index.json`
en el directorio actual).

**OPCIONES**

- `--json` Usa JSON (por defecto)
- `--yaml` Usa YAML 
- `-d N`, `--depth=N`  Crea directorios hasta `N` niveles de profundidad, escribiendo cada
  subproyecto en un archivo separado (por defecto es 0, es decir, sin subdirectorios)
 
**Ejemplos**

Exporta a un único archivo `index.json` en el directorio actual:

```sh
$ todoister export
```

Exporta a un archivo `todoist.json` en el directorio $HOME:

```sh
$ todoister export ~/todoist.json
```

Exporta a un archivo `todoist.yaml` en el directorio $HOME:

```sh
$ todoister export ~/todoist.yaml --yaml
```

Exporta a un directorio `projects` en el $HOME, con subdirectorios hasta 3
niveles de profundidad:

```sh
$ todoister export ~/projects --json -d 3
```

## Ejecución de `export` como una tarea de cron

Se puede ejecutar `todoister export` en una tarea de cron como una forma de crear respaldos
automáticos de Todoist en un formato manejable. Se pueden indicar las opciones de exportación
en el archivo de configuración, de manera que no sea necesario cambiar el cron tab.

```toml
[export]
path = ""
format = ""
depth = 0
```

Por ejemplo, en lugar de `todoister export ~/projects --yaml d 3` puede usar simplemente
`todoist export` con:

```toml
[export]
path = "$HOME/projects"
format = "yaml"
depth = 3
```
Al ejecutarse como una tarea de cron, `todoister export` registra su actividad en un archivo de
log especificado con:

```toml
[log]
name = "/ruta/al/archivo/de/log.log"
```

Consulte `config.toml.example` para un ejemplo completo.

Los logs siguen el formato de [log estructurado](https://pkg.go.dev/log/slog) y se rotan automáticamente.
No se escriben logs en modo interactivo.

## Para desarrolladores

Todoister está escrito en Go (versión mínima 1.22). Utilizo el
[framework Cobra](https://cobra.dev/) para la CLI.
Los comandos residen en `cmd`, el código de utilitarios en `util`.

Uso un Makefile para abreviar tareas rutinarias.

Para actualizar dependencias y actualizar `go.mod` y `go.sum`:

    $ make dependencies

Para ejecutar `golangci-lint` (requiere [golangci-lint](https://golangci-lint.run/)):

    $ make lint

Para construir el binario para su plataforma:

    $ make build

Para instalar el binario en su ruta por defecto:

    $ make install

Para crear un nuevo GitHub release usando la última etiqueta
(requiere [GitHub CLI](https://cli.github.com/)):

    $ make releases

Los *pull requests* son bienvenidos.
