package cmd

import (
	"fmt"
	"github.com/layfellow/todoister/util"
	"github.com/spf13/cobra"
	"strings"
)

func walkProject(project *util.ExportedProject, depth int) {

	// Indent depth * 2 spaces
	fmt.Printf("%s# %s\n", strings.Repeat(" ", depth*2), project.Name)

	if project.Subprojects != nil {
		for _, subproject := range project.Subprojects {
			walkProject(subproject, depth+1)
		}
	}
}

var listCmd = &cobra.Command{
	Use:     "list [project]...",
	Aliases: []string{"ls", "projects"},
	Short:   "List projects",
	Long: `List projects and subprojects.

  project... are the names of one or more project or subproject names to list.
      If no project name is given, all projects are listed.
      You can specify a project name by its full path, e.g., "Work/Project".
      Names are case-insensitive.`,
	Example: `  todoister ls
  todoister ls Work Life
  todoister ls Work/Project`,

	Run: func(cmd *cobra.Command, args []string) {
		projectData := util.HierarchicalData(util.GetTodoistData(ConfigValue.Token))
		project := util.ExportedProject{Subprojects: projectData}
		project.Name = "Projects"

		if len(args) == 0 {
			walkProject(&project, 0)
		} else {
			for _, arg := range args {
				if _, p := util.GetProjectByPathName(arg, &project); p != nil {
					walkProject(p, 0)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
