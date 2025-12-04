package cmd

import (
	"fmt"

	"github.com/layfellow/todoister/util"
	"github.com/spf13/cobra"
)

const (
	tasksLong = `List project tasks.

<code>NAME</code> is the name of one or more projects to list tasks from.
You can specify a project name by its full path, e.g., <code>Work/Project</code>.
Names are case-insensitive.
`

	tasksExample = `# List tasks for project Life:
todoister tasks Life

# List tasks for subproject Project of project Work:
todoister tasks Work/Project

# List tasks for both projects:
todoister tasks Life Work/Project`
)

func printTasks(tasks []*util.ExportedTask) {
	for _, task := range tasks {
		fmt.Printf("  - %s\n", task.Content)
		if task.Description != "" {
			fmt.Printf("%s\n\n", util.IndentMultilineString(task.Description, 4))
		}
	}
}

var tasksCmd = &cobra.Command{
	Use:     "tasks [flags] NAME...",
	Aliases: []string{"items"},
	Short:   "List project tasks",
	Long:    tasksLong,
	Example: tasksExample,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectData := util.HierarchicalData(util.GetTodoistData(ConfigValue.Token))
		project := util.ExportedProject{Subprojects: projectData}
		project.Name = "Projects"

		for _, arg := range args {
			if actualPathname, p := util.GetProjectByPathName(arg, &project); p != nil && p.Tasks != nil {
				fmt.Printf("\n# %s\n\n", actualPathname)
				printTasks(p.Tasks)
				if p.Sections != nil {
					for _, s := range p.Sections {
						fmt.Printf("\n  /%s\n\n", s.Name)
						printTasks(s.Tasks)
					}
				}
			}
		}
	},
}

func init() {
	tasksCmd.SetHelpFunc(util.CustomHelpFunc)
	RootCmd.AddCommand(tasksCmd)
}
