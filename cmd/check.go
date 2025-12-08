// Copyright 2025 Marco Bravo Mejías. All rights reserved.
// Use of this source code is governed by a GPL v3 license
// that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"strings"

	"github.com/layfellow/todoister/util"
	"github.com/spf13/cobra"
)

const checkCmdShortHelp = "Mark a task as completed"

const checkCmdLongHelp = `Mark a TASK in a PROJECT as completed.

Use #[PARENT/SUBPARENT.../]PROJECT to specify the project name with optional
PARENT and SUBPARENTS (note the '#' character prefix and the single quotes).

Alternatively, you can use the --project flag to specify the project name
and omit the '#' prefix and the quotes.

The command matches tasks by prefix (case-insensitive). If multiple tasks
match, an error is shown with a list of matching tasks.`

const checkCmdExample = `  # Check a task in a root project
  todoister check '#Work' 'Write report'
  todoister check -p Work 'Write report'

  # Check a task in a nested project
  todoister check '#Work/Reports' 'Q4 summary'
  todoister check -p Work/Reports 'Q4 summary'`

var checkProjectFlag string

var checkCmd = &cobra.Command{
	Use:     "check [flags] [#][PARENT/.../PROJECT] TASK",
	Short:   checkCmdShortHelp,
	Long:    checkCmdLongHelp,
	Example: checkCmdExample,
	Args:    cobra.RangeArgs(1, 2),
	Run:     runCheckCmd,
}

func init() {
	checkCmd.SetHelpFunc(util.CustomHelpFunc)
	RootCmd.AddCommand(checkCmd)
	checkCmd.Flags().StringVarP(&checkProjectFlag, "project", "p", "", "project name or path (e.g., 'Work' or 'Work/Reports')")
}

func runCheckCmd(cmd *cobra.Command, args []string) {
	var projectPath string
	var taskPrefix string

	// Parse arguments based on whether --project flag is used
	if checkProjectFlag != "" {
		projectPath = checkProjectFlag
		taskPrefix = args[0]
	} else {
		if len(args) != 2 {
			util.Die("Expected 2 arguments when not using --project flag", nil)
		}
		projectPath = strings.TrimPrefix(args[0], "#")
		taskPrefix = args[1]
	}

	// Get all Todoist data
	todoistData := util.GetTodoistData(ConfigValue.Token)
	projects := todoistData.Projects

	// Resolve project path to project ID
	var projectID string
	parts := strings.Split(projectPath, "/")

	if len(parts) > 1 {
		// Nested project path
		projectID = util.GetProjectIDByPathFromProjects(projectPath, projects)
	} else {
		// Single project name - find root project
		projectName := parts[0]
		for _, proj := range projects {
			if strings.EqualFold(proj.Name, projectName) && proj.ParentID == "" {
				projectID = proj.ID
				break
			}
		}
	}

	if projectID == "" {
		util.Die(fmt.Sprintf("Project not found: %s", projectPath), nil)
	}

	// Find tasks matching the prefix
	matches := util.FindTasksByPrefix(projectID, taskPrefix, todoistData)

	if len(matches) == 0 {
		util.Die(fmt.Sprintf("No incomplete tasks found matching '%s' in project '%s'", taskPrefix, projectPath), nil)
	}

	if len(matches) > 1 {
		// Multiple matches - show error with list
		fmt.Printf("Error: Multiple tasks match '%s' in project '%s':\n\n", taskPrefix, projectPath)
		for i, task := range matches {
			status := "incomplete"
			if task.CompletedAt != "" {
				status = "completed"
			}
			fmt.Printf("  %d. [%s] %s (ID: %s)\n", i+1, status, task.Content, task.ID)
		}
		fmt.Println("\nPlease provide a more specific task prefix to match exactly one task.")
		util.Die("", nil)
	}

	// Single match - complete the task
	task := matches[0]

	// Check if already completed (silently succeed)
	if task.CompletedAt != "" {
		// Idempotent behavior - task already done
		fmt.Printf("Task already completed: %s\n", task.Content)
		return
	}

	// Complete the task
	err := util.CompleteTask(ConfigValue.Token, task.ID)
	if err != nil {
		util.Die(fmt.Sprintf("Failed to complete task '%s'", task.Content), err)
	}

	fmt.Printf("✓ Completed: %s\n", task.Content)
}
