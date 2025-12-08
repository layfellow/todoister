package cmd

import (
	"fmt"
	"time"

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
		if task.Due != nil && task.Due.Datetime != "" {
			// Task has a specific datetime in the Datetime field
			dueStr := task.Due.Datetime
			// Try RFC3339 first, then without timezone
			if t, err := time.Parse(time.RFC3339, task.Due.Datetime); err == nil {
				dueStr = t.Format("Jan 2, 2006, 3:04 PM")
			} else if t, err := time.Parse("2006-01-02T15:04:05", task.Due.Datetime); err == nil {
				dueStr = t.Format("Jan 2, 2006, 3:04 PM")
			}
			fmt.Printf("  - %s (%s)\n", task.Content, dueStr)
		} else if task.Due != nil && task.Due.Date != "" {
			// Task has a Date field - may contain date-only or datetime
			dueStr := task.Due.Date
			// Check if Date field contains a datetime (API sometimes puts datetime here)
			if t, err := time.Parse("2006-01-02T15:04:05", task.Due.Date); err == nil {
				dueStr = t.Format("Jan 2, 2006, 3:04 PM")
			} else if t, err := time.Parse(time.RFC3339, task.Due.Date); err == nil {
				dueStr = t.Format("Jan 2, 2006, 3:04 PM")
			} else if t, err := time.Parse("2006-01-02", task.Due.Date); err == nil {
				dueStr = t.Format("Jan 2, 2006")
			}
			fmt.Printf("  - %s (%s)\n", task.Content, dueStr)
		} else {
			fmt.Printf("  - %s\n", task.Content)
		}
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
