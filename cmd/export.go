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
	Use:   "export [flags] [PATH]",
	Short: "Export projects in JSON or YAML format",
	Long: "Export all Todoist projects as a tree of JSON or YAML files.\n\n" +
		"- `PATH` is a file or directory where to export the projects, by default `index.json`.\n",
	Example: "# Export to a single index.json file in the current directory:\n" +
		"todoister export\n\n" +
		"# Export to todoist.json file in the home directory:\n" +
		"todoister export ~/todoist.json\n\n" +
		"# Export to todoist.yaml file in the home directory:\n" +
		"todoister export --yaml ~/todoist.yaml\n\n" +
		"# Export to a projects directory in the home, with subdirectories down to 3 levels deep:\n" +
		"todoister export --json -d 3 ~/projects",

	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		if !useJSON && !useYAML {
			useJSON = strings.ToLower(ConfigValue.Format) == "json"
			useYAML = strings.ToLower(ConfigValue.Format) == "yaml"
		}

		if depth < 0 {
			depth = ConfigValue.Depth
		}

		var exportFormat util.ExportFormat
		var exportPath string
		if useYAML {
			exportFormat = util.YAML
			exportPath = util.YAMLExportPath
		} else {
			exportFormat = util.JSON
			exportPath = util.JSONExportPath
		}

		if ConfigValue.Path != "" {
			exportPath = ConfigValue.Path
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
		"export in JSON format (default)")
	exportCmd.Flags().BoolVar(&useYAML, "yaml", false,
		"export in YAML format")
	exportCmd.Flags().IntVarP(&depth, "depth", "d", -1,
		"depth of subdirectory tree to create on the filesystem when exporting\n(default is 0, i.e., no subdirectories)")
	exportCmd.SetHelpFunc(util.CustomHelpFunc)

	RootCmd.AddCommand(exportCmd)
}
