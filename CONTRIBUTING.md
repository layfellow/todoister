# For developers

[CONTRIBUCIONES en español](CONTRIBUTING.es.md)

Pull requests are welcome. LLM-generated code is acceptable, as long as it is
clearly labeled as such and was reviewed by a human.

Todoister is written in Go (minimum version 1.24). It uses the
[Cobra framework](https://cobra.dev/)
for the CLI. Commands reside in `cmd`, utilities in `util`.

I wrote a Makefile to assist with routine tasks.

To update dependencies and update `go.mod` and `go.sum`:

    $ make dependencies

To run `golangci-lint` (requires [golangci-lint](https://golangci-lint.run/)):

    $ make lint

To build the binary for your platform:

    $ make build

To install the binary in your default path:

    $ make install

To create a new GitHub Release using the latest tag (requires [GitHub CLI](https://cli.github.com/)):

    $ make releases
