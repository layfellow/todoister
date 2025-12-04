package cmd

import (
	"fmt"
	"strings"

	"github.com/layfellow/todoister/util"
	"github.com/spf13/cobra"
)

var (
	projectColor string
	projectFlag  string
)

var addProjectCmd = &cobra.Command{
	Use:   "project [flags] [PARENT/.../]NAME",
	Short: "Add a new project",
	Long: "Add a new project to Todoist.\n\n" +
		"<code>NAME</code> is the name of the project to create.\n" +
		"Use <code>PARENT/NAME</code> to create a project within a parent project.\n" +
		"Use <code>PARENT/SUBPARENT/NAME</code> for nested parents.\n",
	Example: "# Add a root-level project:\n" +
		"todoister add project \"Shopping\"\n\n" +
		"# Add a project within a parent:\n" +
		"todoister add project \"Work/Reports\"\n\n" +
		"# Add a deeply nested project:\n" +
		"todoister add project \"Work/Projects/Q1\"\n\n" +
		"# Add a project with a color:\n" +
		"todoister add project -c blue \"Personal\"\n\n" +
		"# Add a colored project within a parent:\n" +
		"todoister add project --color=red \"Work/Urgent\"",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]

		// Validate color if provided
		if projectColor != "" && !util.ValidColors[projectColor] {
			util.Die(fmt.Sprintf("Invalid color '%s'. Valid colors are: berry_red, red, orange, yellow, olive_green, lime_green, green, mint_green, teal, sky_blue, light_blue, blue, grape, violet, lavender, magenta, salmon, charcoal, grey, taupe", projectColor), nil)
		}

		// Parse the path to extract parent and project name
		parts := strings.Split(path, "/")
		projectName := parts[len(parts)-1]
		var parentPath string
		var parentID string

		// If there are parent parts, we need to find the parent project ID
		if len(parts) > 1 {
			parentPath = strings.Join(parts[:len(parts)-1], "/")

			// Fetch only projects and find the parent ID (lightweight operation)
			projects := util.GetProjects(ConfigValue.Token)
			parentID = util.GetProjectIDByPathFromProjects(parentPath, projects)

			if parentID == "" {
				util.Die(fmt.Sprintf("Parent project '%s' not found", parentPath), nil)
			}
		}

		// Create the project
		project, err := util.CreateProject(ConfigValue.Token, projectName, parentID, projectColor)
		if err != nil {
			util.Die("Failed to create project", err)
		}

		// Print success message
		if parentPath != "" {
			fmt.Printf("Created project '%s' in '%s' (ID: %s)\n", project.Name, parentPath, project.ID)
		} else {
			fmt.Printf("Created project '%s' (ID: %s)\n", project.Name, project.ID)
		}
	},
}

var addTaskCmd = &cobra.Command{
	Use:   "task [flags] [#][PARENT/.../PROJECT] TASK",
	Short: "Add a new task to a project",
	Long: "Add a new task to a Todoist project.\n\n" +
		"Use <code>#[PARENT/SUBPARENT.../]PROJECT</code> to specify the project name with optional\n" +
		"<code>PARENT</code> and <code>SUBPARENTS</code> (note the '<code>#</code>' character prefix and the single quotes).\n\n" +
		"Alternatively, you can use the <code>--project</code> flag to specify the project name\n" +
		"and omit the '<code>#</code>' prefix and the quotes.\n",
	Example: "# Add task to root-level project Work:\n" +
		"todoister add task '#Work' 'Complete report'\n\n" +
		"# Add task to project Reports of project Work:\n" +
		"todoister add task '#Work/Reports' 'Create quarterly report'\n\n" +
		"# Add tasks using project flag:\n" +
		"todoister add task -p Work/Reports 'Create monthly report'\n" +
		"todoister add task -p Personal 'Buy groceries'\n\n" +
		"# Add task to nested project using flag:\n" +
		"todoister add task --project=Personal/Shopping/List 'Buy milk'",
	Args: func(cmd *cobra.Command, args []string) error {
		// Handle both argument formats
		if projectFlag != "" {
			// Using -p/--project flag: expect exactly 1 argument (task title)
			if len(args) != 1 {
				return fmt.Errorf("when using --project flag, only TASK_TITLE is required")
			}
		} else {
			// Not using flag: expect exactly 2 arguments (project path and task title)
			if len(args) != 2 {
				return fmt.Errorf("expected PROJECT_PATH and TASK_TITLE, or use --project flag")
			}
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var projectPath, taskTitle string

		if projectFlag != "" {
			// Using project flag
			projectPath = projectFlag
			taskTitle = args[0]
		} else {
			// Using positional arguments
			projectPath = args[0]
			taskTitle = args[1]

			// Remove leading # if present
			projectPath = strings.TrimPrefix(projectPath, "#")
		}

		// Parse the project path to extract parent and project name
		parts := strings.Split(projectPath, "/")
		projectName := parts[len(parts)-1]
		var parentPath string
		var projectID string

		// If there are parent parts, we need to find the parent project ID
		if len(parts) > 1 {
			parentPath = strings.Join(parts[:len(parts)-1], "/")

			// Fetch only projects and find the project ID (lightweight operation)
			projects := util.GetProjects(ConfigValue.Token)
			projectID = util.GetProjectIDByPathFromProjects(projectPath, projects)

			if projectID == "" {
				util.Die(fmt.Sprintf("Project '%s' not found", projectPath), nil)
			}
		} else {
			// Single project name - find it by name
			projects := util.GetProjects(ConfigValue.Token)
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

		// Create the task
		task, err := util.CreateTask(ConfigValue.Token, taskTitle, projectID)
		if err != nil {
			util.Die("Failed to create task", err)
		}

		// Print success message
		if parentPath != "" {
			fmt.Printf("Created task '%s' in '%s/%s' (ID: %s)\n", task.Content, parentPath, projectName, task.ID)
		} else {
			fmt.Printf("Created task '%s' in '%s' (ID: %s)\n", task.Content, projectName, task.ID)
		}
	},
}

var addCmd = &cobra.Command{
	Use:   "add <resource> [arguments]",
	Short: "Add a new resource",
	Long:  "Add a new resource to Todoist (currently supports: project, task).\n",
}

func init() {
	addProjectCmd.Flags().StringVarP(&projectColor, "color", "c", "",
		"project color (berry_red, red, orange, yellow, olive_green, lime_green, green, mint_green, teal, sky_blue, light_blue, blue, grape, violet, lavender, magenta, salmon, charcoal, grey, taupe)")
	addProjectCmd.SetHelpFunc(util.CustomHelpFunc)

	addTaskCmd.Flags().StringVarP(&projectFlag, "project", "p", "",
		"project name or path (e.g., 'Work' or 'Work/Reports')")
	addTaskCmd.SetHelpFunc(util.CustomHelpFunc)

	addCmd.AddCommand(addProjectCmd)
	addCmd.AddCommand(addTaskCmd)
	addCmd.SetHelpFunc(util.CustomHelpFunc)

	RootCmd.AddCommand(addCmd)
}
