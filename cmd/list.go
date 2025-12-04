package cmd

import (
	"fmt"
	"strings"

	"github.com/layfellow/todoister/util"
	"github.com/spf13/cobra"
)

const (
	listLong = `List projects and subprojects.

<code>NAME</code> is the name of one or more projects to list tasks from.
If no <code>NAME</code> is given, all projects are listed.
You can specify a project name by its full path, e.g., <code>Work/Project</code>.
Names are case-insensitive.
`

	listExample = `# List all projects and subprojects:
todoister ls

# List projects Work and Life and their subprojects:
todoister ls Work Life

# List all subprojects of Project, which is a subproject of Work:
todoister ls Work/Project`
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
	Use:     "list [flags] [NAME]...",
	Aliases: []string{"ls", "projects"},
	Short:   "List projects",
	Long:    listLong,
	Example: listExample,
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
	listCmd.SetHelpFunc(util.CustomHelpFunc)
	RootCmd.AddCommand(listCmd)
}
