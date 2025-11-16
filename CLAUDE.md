# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Todoister is a CLI client for Todoist written in Go. It allows users to list projects, view tasks, and export Todoist data to JSON or YAML format. The project emphasizes simplicity and uses structured, hierarchical data representation for exports.

## Development Commands

### Building
```sh
# Build for current platform
make build

# Build for all platforms (Linux amd64, macOS amd64/arm64)
# Outputs to build/ directory
make build
```

The build process uses ldflags to inject the version number from `cmd.Version` variable.

### Testing and Quality
```sh
# Run tests
make test

# Run golangci-lint (requires golangci-lint to be installed)
make lint

# Update dependencies
make dependencies
```

### Installation
```sh
# Install to ~/bin
make install

# Install via go install (requires Go 1.22+)
go install github.com/layfellow/todoister@latest
```

### Documentation
```sh
# Generate documentation (runs doc/doc.go)
make doc
```

### Releases
```sh
# Create GitHub release with binaries using gh CLI
make releases
```

## Architecture

### CLI Framework
Uses **Cobra** for command-line interface structure. All commands are in `cmd/` directory:
- `cmd/root.go` - Root command and initialization
- `cmd/list.go` - List projects
- `cmd/tasks.go` - List tasks
- `cmd/export.go` - Export to JSON/YAML
- `cmd/version.go` - Version information

### Configuration System
Configuration uses **Viper** with multi-source priority (highest to lowest):
1. `--token` CLI flag
2. `TODOIST_TOKEN` environment variable
3. `~/.config/todoister/config.toml` (XDG-compliant)
4. `~/.todoister.toml` (fallback)

Configuration structure in `util/config.go`:
```go
type ConfigType struct {
    Token  string
    Log    Log     // For cron job logging
    Export Export  // Default export settings
}
```

### API Integration
The Todoist unified API v1 is used via `util/api.go`:
- Base URL: `https://api.todoist.com/api/v1`
- Multiple endpoints:
  - `GET /api/v1/projects` - all active projects
  - `GET /api/v1/sections` - all sections
  - `GET /api/v1/tasks` - all active tasks
  - `GET /api/v1/labels` - all labels
  - `GET /api/v1/comments?project_id={id}` - comments per project
- Returns flat `TodoistData` structure combining all responses
- **Note**: This is the current unified API that replaced the deprecated Sync API v9 and REST API v2

### Data Transformation Pipeline
The key architectural pattern is converting flat API data to hierarchical structures:

1. **API Layer** (`util/api.go`): Makes multiple unified API v1 calls and combines into flat `TodoistData`
   - Fetches projects, sections, tasks, labels in separate GET requests
   - Fetches comments for each project (one request per project)
   - Combines all responses into single `TodoistData` structure
2. **Parser Layer** (`util/parser.go`): `HierarchicalData()` transforms flat data into nested `ExportedProject` tree
   - Creates maps for O(1) lookups by ID
   - Builds parent-child relationships for projects, sections, tasks
   - Links associated data (labels, comments) to their parents
3. **Writer Layer** (`util/writer.go`): Serializes hierarchical data to JSON/YAML files

### Data Structures
Two parallel type hierarchies exist in `util/parser.go`:
- **Todoist types**: Unified API v1 response structures (e.g., `TodoistProject`, `TodoistItem`)
  - Use simple string IDs (no `v2_` prefix like old Sync API v9)
  - Duration and Due are pointers (can be null)
  - Comments unified (single type for both task and project comments)
- **Exported types**: Hierarchical structures for export (e.g., `ExportedProject`, `ExportedTask`)

Key relationships:
- Projects can contain subprojects, sections, tasks, and comments
- Sections contain tasks
- Tasks can have labels, comments, duration, and due dates

**API Migration**: The code uses the unified Todoist API v1, which replaced the deprecated Sync API v9 and REST API v2. The unified API provides standard REST endpoints for all resources.

### Logging
Structured logging via Go's `log/slog` in `util/logger.go`:
- Only writes logs in non-interactive mode (cron jobs)
- Auto-rotates log files
- Configured via `log.name` in config file

### Utility Functions
`util/` directory contains shared functionality:
- `util/help.go` - Custom help formatting
- `util/writer.go` - JSON/YAML file writing with depth control
- Path expansion, error handling (`Die`, `Warn`)

## Key Implementation Details

### Version Injection
Version is set via ldflags during build:
```makefile
LDFLAGS= -ldflags="-X 'github.com/layfellow/todoister/cmd.Version=$(VERSION)'"
```

### Export Depth Control
The `-d` flag controls directory depth when exporting:
- `depth=0`: Single file (default)
- `depth>0`: Creates nested directory structure mirroring project hierarchy

### Configuration Precedence
Commands check for config values in this order:
1. Command-line flags
2. Config file values
3. Hardcoded defaults

This pattern appears in `cmd/export.go` for handling format and path settings.

### Error Handling
Uses `util.Die()` for fatal errors and `util.Warn()` for non-fatal warnings. All errors are user-facing with descriptive messages.

## Testing Locally

To test with actual Todoist data:
1. Create API token at https://app.todoist.com/app/settings/integrations/developer
2. Set up config file or environment variable
3. Run commands: `go run . list`, `go run . export`
