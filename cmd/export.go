package cmd

import (
	"github.com/spf13/cobra"
	"strings"
	"todoister/util"
)

var useJSON bool
var useYAML bool
var depth int

var exportCmd = &cobra.Command{
	Use:   "export [path]",
	Short: "Export projects in JSON or YAML format",
	Long: `Export all Todoist projects to path (file or directory)
The default is index.json or index.yaml (if --yaml is specified) in the
current directory.
if --depth is provided, todoister will create directories up to N levels deep,
writing each subproject to a separate file.`,
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
