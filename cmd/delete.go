package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/layfellow/todoister/util"
	"github.com/spf13/cobra"
)

const (
	deleteLong = `Delete a resource from Todoist (currently supports: project).
`

	deleteProjectLong = `Delete a project from Todoist.

<code>NAME</code> is the name of the project to delete.
Use <code>PARENT/NAME</code> to locate a project within a parent project.
Use <code>PARENT/SUBPARENT/NAME</code> for nested parents.

This command deletes the project and all its descendants (subprojects and tasks).
`

	deleteProjectExample = `# Delete a root-level project:
todoister delete project Shopping

# Delete a project within a parent:
todoister delete project Work/Reports

# Delete a deeply nested project:
todoister delete project Work/Projects/Q1

# Delete without confirmation:
todoister delete project -f Shopping
todoister rm project --force Work/Old`
)

var forceDelete bool

var deleteProjectCmd = &cobra.Command{
	Use:     "project [flags] [PARENT/.../]NAME",
	Short:   "Delete a project",
	Long:    deleteProjectLong,
	Example: deleteProjectExample,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]

		// Fetch Todoist data and find the project ID
		todoistData := util.GetTodoistData(ConfigValue.Token)
		projectID := util.GetProjectIDByPathFromProjects(path, todoistData.Projects)

		if projectID == "" {
			util.Die(fmt.Sprintf("Project '%s' not found", path), nil)
		}

		// Unless --force is set, prompt for confirmation
		if !forceDelete {
			fmt.Printf("Delete project '%s' and all its descendants? [y/N]: ", path)
			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')
			if err != nil {
				util.Die("Failed to read input", err)
			}
			if strings.ToLower(strings.TrimSpace(response)) != "y" {
				return
			}
		}

		// Delete the project
		err := util.DeleteProject(ConfigValue.Token, projectID)
		if err != nil {
			util.Die("Failed to delete project", err)
		}

		fmt.Printf("Deleted project '%s'\n", path)
	},
}

var deleteCmd = &cobra.Command{
	Use:     "delete <resource> [arguments]",
	Aliases: []string{"del", "rm"},
	Short:   "Delete a resource",
	Long:    deleteLong,
}

func init() {
	deleteProjectCmd.Flags().BoolVarP(&forceDelete, "force", "f", false,
		"skip confirmation prompt")
	deleteProjectCmd.SetHelpFunc(util.CustomHelpFunc)

	deleteCmd.AddCommand(deleteProjectCmd)
	deleteCmd.SetHelpFunc(util.CustomHelpFunc)

	RootCmd.AddCommand(deleteCmd)
}
