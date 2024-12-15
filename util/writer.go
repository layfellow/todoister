package util

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

// ExportFormat is either JSON or YAML.
type ExportFormat int

const (
	JSON = iota
	YAML
)

// Deafult export paths.
const (
	JSONExportPath = "index.json"
	YAMLExportPath = "index.yaml"
)

// ExpandPath expands a path that starts with a tilde (~) or contains environment variables.
// - path: the path to expand
//
// Returns the expanded path and an error, if any.
func ExpandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(homeDir, path[1:])
	}
	path = os.ExpandEnv(path)
	return path, nil
}

// normalizePathnames extracts the directory and basename from a path, setting the basename
// to the default if it isnʼt a filename.
// - path: the path to normalize
// - format: the export format (JSON or YAML)
//
// Returns the directory and basename and an error, if any.
func normalizePathnames(path string, format ExportFormat) (string, string, error) {
	expandedPath, err := ExpandPath(path)
	if err != nil {
		return "", "", err
	}
	dirname := filepath.Dir(expandedPath)
	basename := filepath.Base(expandedPath)

	// filepath.Dir returns "." for empty paths
	if dirname == "." {
		dirname = ""
	}
	// If the basename is empty or isnʼt a filename, use the default export path.
	if basename == "." || !strings.Contains(basename, ".") {
		dirname = expandedPath
		basename = JSONExportPath
		if format == YAML {
			basename = YAMLExportPath
		}
	}
	return dirname, basename, nil
}

// clearDirectory removes the directory and its contents, then creates a new directory.
// - dirname: the directory to clear
func clearDirectory(dirname string) error {
	pathInfo, err := os.Stat(dirname)
	if err == nil && pathInfo.IsDir() {
		err = os.RemoveAll(dirname)
		if err != nil {
			return err
		}
	}
	if err = os.MkdirAll(dirname, 0750); err != nil {
		return err
	}
	return nil
}

// writeProjectFile writes a project to a file in the specified format.
// - path: the path to the file
// - format: the export format (JSON or YAML)
// - project: a pointer to the project to write, already prepared for export
// - includeSubprojects: whether to include subprojects in the export
func writeProjectFile(path string, format ExportFormat, project *ExportedProject, includeSubprojects bool) error {
	outputFile, err := os.Create(path)
	if err != nil {
		return err
	}

	defer func() {
		if cerr := outputFile.Close(); cerr != nil {
			Warn("Failed to close output file", cerr)
		}
	}()

	outputProject := project
	if !includeSubprojects {
		// Make a temporary copy of project without subprojects.
		outputProject = &ExportedProject{
			Project:     project.Project,
			Subprojects: nil,
			Sections:    project.Sections,
			Tasks:       project.Tasks,
			Comments:    project.Comments,
		}
	}

	var encoder interface{ Encode(v interface{}) error }
	if format == JSON {
		encoder = json.NewEncoder(outputFile)
		encoder.(*json.Encoder).SetIndent("", "  ")
	} else {
		encoder = yaml.NewEncoder(outputFile)
		encoder.(*yaml.Encoder).SetIndent(2)
	}

	if err := encoder.Encode(outputProject); err != nil {
		return err
	}
	return nil
}

// writeProject writes a project to a file or directory in the specified format.
// - dirname: the directory to write to, already normalized
// - basename: the base filename, already normalized
// - format: the export format (JSON or YAML)
// - depth: the depth of the subdirectory tree
// - project: a pointer to the project to write
func writeProject(dirname string, basename string, format ExportFormat, depth int, clear bool, project *ExportedProject) error {
	if clear && dirname != "" {
		if err := clearDirectory(dirname); err != nil {
			Warn("Failed to clear directory", err)
		}
	}
	if project.Subprojects == nil || len(project.Subprojects) == 0 || depth == 0 {
		err := writeProjectFile(filepath.Join(dirname, basename), format, project, true)
		if err != nil {
			return err
		}
	} else {
		err := writeProjectFile(filepath.Join(dirname, basename), format, project, false)
		if err != nil {
			return err
		}
		for _, subproject := range project.Subprojects {
			// Recursively walk the subprojects.
			err := writeProject(filepath.Join(dirname, subproject.Name), basename, format, depth-1, true, subproject)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// WriteHierarchicalData writes a hierarchical project structure to a file or directory in the specified format.
// - roots: a slice of root ExportedProject references
// - format: the export format (JSON or YAML)
// - depth: the depth of the subdirectory tree
// - path: the path to write to
func WriteHierarchicalData(roots []*ExportedProject, format ExportFormat, depth int, path string) error {
	dirname, basename, err := normalizePathnames(path, format)

	// DEBUG
	println("format: ", format)
	println("dirname: ", dirname)
	println("basename: ", basename)

	if err != nil {
		return err
	}
	project := ExportedProject{Subprojects: roots}
	project.Name = "Projects"
	if err := writeProject(dirname, basename, format, depth, false, &project); err != nil {
		return err
	}
	return nil
}
