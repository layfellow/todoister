package cmd

import (
	"strings"

	"github.com/layfellow/todoister/util"
	"github.com/spf13/cobra"
)

const (
	exportLong = `Export all Todoist projects as a tree of JSON or YAML files.

- <code>PATH</code> is a file or directory where to export the projects, by default <code>index.json</code>.
`

	exportExample = `# Export to a single index.json file in the current directory:
todoister export

# Export to todoist.json file in the home directory:
todoister export ~/todoist.json

# Export to todoist.yaml file in the home directory:
todoister export --yaml ~/todoist.yaml

# Export to a projects directory in the home, with subdirectories down to 3 levels deep:
todoister export --json -d 3 ~/projects`
)

var useJSON bool
var useYAML bool
var depth int

var exportCmd = &cobra.Command{
	Use:     "export [flags] [PATH]",
	Short:   "Export projects in JSON or YAML format",
	Long:    exportLong,
	Example: exportExample,
	Args:    cobra.MaximumNArgs(1),
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
