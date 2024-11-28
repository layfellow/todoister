package cmd

import (
	"fmt"
	"github.com/layfellow/todoister/util"
	"github.com/spf13/cobra"
	"strings"
)

func indentMultilineString(s string, indent int) string {
	lines := strings.Split(s, "\n")
	indentStr := strings.Repeat(" ", indent)
	for i, line := range lines {
		lines[i] = indentStr + line
	}
	return strings.Join(lines, "\n")
}

func printTasks(tasks []*util.ExportedTask) {
	for _, task := range tasks {
		fmt.Printf("  - %s\n", task.Content)
		if task.Description != "" {
			fmt.Printf("%s\n\n", indentMultilineString(task.Description, 4))
		}
	}
}

var tasksCmd = &cobra.Command{
	Use:     "tasks project...",
	Aliases: []string{"items"},
	Short:   "List project tasks",
	Long: `List project tasks.

  project... are the names of one or more projects whose tasks to list.
      You can specify a project name by its full path, e.g., "Work/Project".
      Names are case-insensitive.`,
	Example: `  todoister tasks Life
  todoister tasks Work/Project`,

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
	RootCmd.AddCommand(tasksCmd)
}
