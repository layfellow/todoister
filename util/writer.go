package util

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

type ExportFormat int

const (
	JSON = iota
	YAML
)

const (
	DefaultExportPath = "index.json"
	YAMLExportPath    = "index.yaml"
)

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

func normalizePathnames(path string, format ExportFormat) (string, string, error) {
	expandedPath, err := ExpandPath(path)
	if err != nil {
		return "", "", err
	}
	dirname := filepath.Dir(expandedPath)
	basename := filepath.Base(expandedPath)

	if dirname == "." {
		dirname = ""
	}
	if basename == "." || !strings.Contains(basename, ".") {
		dirname = expandedPath
		basename = DefaultExportPath
		if format == YAML {
			basename = YAMLExportPath
		}
	}
	return dirname, basename, nil
}

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

func writeProject(dirname string, basename string, format ExportFormat, depth int, project *ExportedProject) error {
	if dirname != "" {
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
			err := writeProject(filepath.Join(dirname, subproject.Name), basename, format, depth-1, subproject)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func WriteHierarchicalData(roots []*ExportedProject, format ExportFormat, depth int, path string) error {
	dirname, basename, err := normalizePathnames(path, format)
	if err != nil {
		return err
	}
	project := ExportedProject{Subprojects: roots}
	project.Name = "Projects"
	if err := writeProject(dirname, basename, format, depth, &project); err != nil {
		return err
	}
	return nil
}
