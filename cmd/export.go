package cmd

import (
	"github.com/layfellow/todoister/util"
	"github.com/spf13/cobra"
	"strings"
)

var useJSON bool
var useYAML bool
var depth int

var exportCmd = &cobra.Command{
	Use:   "export [path]",
	Short: "Export projects in JSON or YAML format",
	Long: `Export all Todoist projects.

  path is a file or directory where to export the projects, by default index.json
       or index.yaml (if --yaml is specified) in the current directory.

  -d, --depth N when provided, todoister will create directories up to N levels deep,
       writing each subproject to a separate file.`,
	Example: `  todoister export
  todoister export ~/todoist.json
  todoister export ~/todoist.yaml --yaml
  todoister export ~/projects --json -d 3`,

	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		useJSON = useJSON || strings.ToLower(ConfigValue.Export.Format) == "json"
		useYAML = !useJSON && (useYAML || strings.ToLower(ConfigValue.Export.Format) == "yaml")
		if depth == 0 {
			depth = ConfigValue.Export.Depth
		}

		var exportFormat util.ExportFormat = util.JSON
		exportPath := util.DefaultExportPath
		if useYAML {
			exportFormat = util.YAML
			exportPath = util.YAMLExportPath
		}
		if ConfigValue.Export.Path != "" {
			exportPath = ConfigValue.Export.Path
		}
		if len(args) > 0 {
			exportPath = args[0]
		}

		hierarchicalData := util.HierarchicalData(util.GetTodoistData(ConfigValue.Token))
		err := util.WriteHierarchicalData(hierarchicalData, exportFormat, depth, exportPath)
		if err != nil {
			util.Die("Failed to export", err)
		}
	},
}

func init() {
	exportCmd.Flags().BoolVar(&useJSON, "json", false,
		"Export in JSON format (default)")
	exportCmd.Flags().BoolVar(&useYAML, "yaml", false,
		"Export in YAML format")
	exportCmd.Flags().IntVarP(&depth, "depth", "d", 0,
		"Depth of subdirectory tree when exporting (default 0, i.e., no subdirectories)")

	rootCmd.AddCommand(exportCmd)
}
