# Para desarrolladores

[CONTRIBUTING in English](CONTRIBUTING.md)

Los *pull requests* son bienvenidos. Se acepta el código generado por LLM,
siempre y cuando esté claramente etiquetado como tal y haya sido revisado
por un humano.

Todoister está escrito en Go (versión mínima 1.24). Utiliza el
[framework Cobra](https://cobra.dev/)
para la CLI. Los comandos residen en `cmd`, las utilidades en `util`.

Uso un Makefile para simplificar algunas tareas rutinarias.

Para actualizar dependencias y actualizar `go.mod` y `go.sum`:

    $ make dependencies

Para ejecutar `golangci-lint` (requiere [golangci-lint](https://golangci-lint.run/)):

    $ make lint

Para construir el binario para su plataforma:

    $ make build

Para instalar el binario en su ruta predeterminada:

    $ make install

Para crear un nuevo GitHub Release usando la última etiqueta (requiere [GitHub CLI](https://cli.github.com/)):

    $ make releases
