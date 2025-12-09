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
	deleteLong = `Delete a resource from Todoist (currently supports: project, task).
`

	deleteProjectLong = `Delete a project from Todoist.

<code>NAME</code> is the name of the project to delete.
Use <code>PARENT/NAME</code> to locate a project within a parent project.
Use <code>PARENT/SUBPARENT/NAME</code> for nested parents.
Note that <code>NAMES</code>, <code>PARENTS</code> and <code>SUBPARENTS</code> are case-insensitive.

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

	deleteTaskLong = `Delete a task from Todoist.

Use <code>#[PARENT/SUBPARENT.../]PROJECT</code> to specify the project name with optional
<code>PARENT</code> and <code>SUBPARENTS</code> (note the '<code>#</code>' character prefix and the single quotes).

Alternatively, you can use the <code>--project</code> flag to specify the project name
and omit the '<code>#</code>' prefix and the quotes.
Note that <code>PROJECT</code>, <code>PARENTS</code> and <code>SUBPARENTS</code> are case-insensitive.

You can identify a <code>TASK</code> by its partial name. If multiple tasks match,
an error is shown.

This command deletes the task and all its sub-tasks.
`

	deleteTaskExample = `# Delete task from root-level project Work:
todoister delete task '#Work' 'Complete report'

# Delete task from nested project:
todoister delete task '#Work/Reports' 'Create quarterly report'

# Delete task using partial name:
todoister delete task '#Work/Reports' 'Create q'

# Delete task using project flag:
todoister delete task -p Work/Reports 'Create monthly report'

# Delete task without confirmation:
todoister delete task -f -p Personal 'Buy groceries'
todoister rm task --force '#Work' 'Old task'`
)

var (
	forceDelete       bool
	deleteProjectFlag string
)

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

var deleteTaskCmd = &cobra.Command{
	Use:     "task [flags] [#][PARENT/.../PROJECT] TASK",
	Short:   "Delete a task from a project",
	Long:    deleteTaskLong,
	Example: deleteTaskExample,
	Args: func(cmd *cobra.Command, args []string) error {
		if deleteProjectFlag != "" {
			if len(args) != 1 {
				return fmt.Errorf("when using --project flag, only TASK is required")
			}
		} else {
			if len(args) != 2 {
				return fmt.Errorf("expected PROJECT_PATH and TASK, or use --project flag")
			}
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var projectPath, taskContent string

		if deleteProjectFlag != "" {
			projectPath = deleteProjectFlag
			taskContent = args[0]
		} else {
			projectPath = args[0]
			taskContent = args[1]
			projectPath = strings.TrimPrefix(projectPath, "#")
		}

		// Parse the project path to extract parent and project name
		parts := strings.Split(projectPath, "/")
		projectName := parts[len(parts)-1]
		var projectID string

		// Fetch Todoist data and find the project ID
		todoistData := util.GetTodoistData(ConfigValue.Token)
		projects := todoistData.Projects

		if len(parts) > 1 {
			projectID = util.GetProjectIDByPathFromProjects(projectPath, projects)
			if projectID == "" {
				util.Die(fmt.Sprintf("Project '%s' not found", projectPath), nil)
			}
		} else {
			for _, proj := range projects {
				if strings.EqualFold(proj.Name, projectName) && proj.ParentID == "" {
					projectID = proj.ID
					break
				}
			}
			if projectID == "" {
				util.Die(fmt.Sprintf("Root-level project '%s' not found", projectName), nil)
			}
		}

		// Find matching tasks using prefix match
		matches := util.FindTasksByPrefix(projectID, taskContent, todoistData)

		if len(matches) == 0 {
			util.Die(fmt.Sprintf("Task '%s' not found in project '%s'", taskContent, projectPath), nil)
		}

		if len(matches) > 1 {
			msg := fmt.Sprintf("Multiple tasks match '%s':\n", taskContent)
			for _, task := range matches {
				msg += fmt.Sprintf("  - %s\n", task.Content)
			}
			msg += "Please provide a more specific task name."
			util.Die(msg, nil)
		}

		task := matches[0]

		// Unless --force is set, prompt for confirmation
		if !forceDelete {
			fmt.Printf("Delete task '%s'? [y/N]: ", task.Content)
			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')
			if err != nil {
				util.Die("Failed to read input", err)
			}
			if strings.ToLower(strings.TrimSpace(response)) != "y" {
				return
			}
		}

		// Delete the task
		err := util.DeleteTask(ConfigValue.Token, task.ID)
		if err != nil {
			util.Die("Failed to delete task", err)
		}

		fmt.Printf("Deleted task '%s'\n", task.Content)
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

	deleteTaskCmd.Flags().StringVarP(&deleteProjectFlag, "project", "p", "",
		"project name or path (e.g., 'Work' or 'Work/Reports')")
	deleteTaskCmd.Flags().BoolVarP(&forceDelete, "force", "f", false,
		"skip confirmation prompt")
	deleteTaskCmd.SetHelpFunc(util.CustomHelpFunc)

	deleteCmd.AddCommand(deleteProjectCmd)
	deleteCmd.AddCommand(deleteTaskCmd)
	deleteCmd.SetHelpFunc(util.CustomHelpFunc)

	RootCmd.AddCommand(deleteCmd)
}
