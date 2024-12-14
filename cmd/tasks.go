package cmd

import (
	"fmt"
	"github.com/layfellow/todoister/util"
	"github.com/spf13/cobra"
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
	Use:     "tasks project...",
	Aliases: []string{"items"},
	Short:   "List project tasks",
	Long: "List project tasks.\n\n" +
		"`project` is the name of one or more projects to list tasks from.\n" +
		"You can specify a project name by its full path, e.g., `Work/Project`.\n" +
		"Names are case-insensitive.\n",
	Example: "# List tasks for project Life:\n" +
		"todoister tasks Life\n\n" +
		"# List tasks for subproject Project of project Work:\n" +
		"todoister tasks Work/Project\n\n",
	Args: cobra.MinimumNArgs(1),
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
