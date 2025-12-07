# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Todoister is a CLI client for Todoist written in Go using the Cobra framework. It provides commands to list projects, view tasks, add tasks and projects, and export Todoist data to structured JSON/YAML formats.

## Development Commands

### Building
```bash
# Build for local platform
go build -ldflags="-X 'github.com/layfellow/todoister/cmd.Version=0.3.0'" -o build/todoister

# Build all platform binaries
make build

# Regenerate protobuf code (if modifying todoist.proto)
protoc --go_out=. --go_opt=paths=source_relative util/todoist.proto
```

### Testing
```bash
# Run all tests
make test

# Run specific test file
VERSION=0.3.0 go test -count=1 ./cmd -run TestAdd

# Run single test case
VERSION=0.3.0 go test -count=1 ./cmd -run TestVersion/version_flag
```

### Other Commands
```bash
make lint          # Run golangci-lint
make doc           # Generate documentation (runs doc/doc.go)
make dependencies  # Update dependencies and tidy go.mod
make install       # Install to ~/bin
make clean         # Remove build directory
```

## Architecture

### Package Structure

- **`main.go`**: Entry point that calls `cmd.Execute()`
- **`cmd/`**: Command implementations using Cobra
  - `root.go`: Root command setup, version, config initialization
  - `add.go`: Add projects and tasks (`add project`, `add task`)
  - `list.go`: List projects and subprojects
  - `tasks.go`: List tasks from specific projects
  - `export.go`: Export all Todoist data to JSON/YAML
  - `version.go`: Version command
  - `*_test.go`: Command tests
- **`util/`**: Core utility functions
  - `api.go`: Todoist Sync API v1 client functions
  - `cache.go`: Protobuf cache management (load/save)
  - `sync.go`: Data conversion between cache and application formats, incremental merge logic
  - `todoist.proto`: Protobuf schema for cache (generates `todoist.pb.go`)
  - `parser.go`: Data structure definitions and hierarchical parsing
  - `writer.go`: JSON/YAML export with directory tree generation
  - `config.go`: Configuration file and environment variable handling
  - `logger.go`: Structured logging (slog)
  - `help.go`: Custom help formatting
- **`build/`**: Build artifacts (binaries for different platforms)
- **`doc/`**: Documentation generation

### Data Flow

1. **API Data Retrieval with Caching** (`util/api.go`, `util/cache.go`, `util/sync.go`):
   - `GetTodoistData()` uses the Todoist Sync API v1 endpoint with local Protobuf caching
   - On first run: Performs full sync (`sync_token="*"`) and caches data to `~/.cache/todoister/todoist.pb`
   - On subsequent runs: Performs incremental sync using cached `sync_token`, fetching only changes
   - Cache format: Protocol Buffers (see `util/todoist.proto`)
   - Sync request fetches all data in one HTTP call: projects, sections, items, labels, notes, project_notes
   - `CreateTask()` and `CreateProject()` still use REST API POST endpoints for writes

2. **Data Transformation** (`util/parser.go`):
   - `HierarchicalData()` converts flat Todoist API data into nested `ExportedProject` structures
   - Builds parent-child relationships for projects, sections, and tasks
   - Links comments and labels to their respective tasks/projects
   - Returns slice of root-level `ExportedProject` pointers

3. **Path Resolution** (`util/parser.go`):
   - `GetProjectIDByPath()` resolves project paths like `"Work/Reports"` to Todoist project IDs
   - `GetProjectIDByPathFromProjects()` is a lightweight alternative using only project data
   - `GetProjectByPathName()` navigates hierarchical structures to find projects

4. **Data Export** (`util/writer.go`):
   - `WriteHierarchicalData()` writes hierarchical data to JSON or YAML
   - Supports configurable depth for filesystem directory trees
   - `writeProject()` recursively creates subdirectories for nested projects

### Configuration

Configuration is loaded via Viper from (in order of precedence):
1. `--token` command-line flag
2. `TODOIST_TOKEN` environment variable
3. `~/.config/todoister/config.toml` (or `$XDG_CONFIG_HOME/todoister/config.toml`)
4. `~/.todoister.toml` (fallback)

Example configuration:
```toml
token = "your-todoist-API-token"

[log]
name = "/path/to/log/file.log"

[export]
path = "$HOME/projects"
format = "yaml"
depth = 3
```

### Project Path Syntax

Commands that reference projects accept hierarchical paths:
- Root project: `"Work"`
- Nested project: `"Work/Reports"`
- Deeply nested: `"Work/Projects/Q1"`

The `add task` command uses `#` prefix for project paths in positional arguments:
- `todoister add task '#Work/Reports' 'Task title'`
- Or use the `--project` flag: `todoister add task -p Work/Reports 'Task title'`

### Cache Management

The application uses a local Protobuf cache for performance:

- **Cache location**: `~/.cache/todoister/todoist.pb` (or `$XDG_CACHE_HOME/todoister/todoist.pb`)
- **First run**: Performs full sync and creates cache
- **Subsequent runs**: Incremental sync fetches only changes since last sync
- **Cache corruption**: Automatically falls back to full sync
- **Clear cache**: Delete `~/.cache/todoister/todoist.pb` to force full resync

### Color Values for Projects

When creating projects with the `--color` flag, use one of these valid values:
`berry_red`, `red`, `orange`, `yellow`, `olive_green`, `lime_green`, `green`, `mint_green`, `teal`, `sky_blue`, `light_blue`, `blue`, `grape`, `violet`, `lavender`, `magenta`, `salmon`, `charcoal`, `grey`, `taupe`

### Version Management

The version is injected at build time via ldflags:
```bash
-ldflags="-X 'github.com/layfellow/todoister/cmd.Version=0.3.0'"
```

Tests require the `VERSION` environment variable to be set.

## Testing Notes

- Tests are in `cmd/*_test.go`
- All tests must set `VERSION` environment variable before running
- Use `-count=1` to disable test caching
- The test framework expects specific version output formats
