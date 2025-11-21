package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/layfellow/todoister/util"
	"github.com/spf13/cobra"
)

var (
	projectColor string
	projectFlag  string
)

// validColors are the allowed color values for projects
var validColors = map[string]bool{
	"berry_red":   true,
	"red":         true,
	"orange":      true,
	"yellow":      true,
	"olive_green": true,
	"lime_green":  true,
	"green":       true,
	"mint_green":  true,
	"teal":        true,
	"sky_blue":    true,
	"light_blue":  true,
	"blue":        true,
	"grape":       true,
	"violet":      true,
	"lavender":    true,
	"magenta":     true,
	"salmon":      true,
	"charcoal":    true,
	"grey":        true,
	"taupe":       true,
}

// ProjectCreateRequest represents the request body for creating a project
type ProjectCreateRequest struct {
	Name     string `json:"name"`
	ParentID string `json:"parent_id,omitempty"`
	Color    string `json:"color,omitempty"`
}

// ProjectResponse represents a project response from the API
type ProjectResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ParentID string `json:"parent_id"`
	Color    string `json:"color"`
}

// createProject makes a POST request to create a new project
func createProject(token, name, parentID, color string) (*ProjectResponse, error) {
	client := &http.Client{}

	reqBody := ProjectCreateRequest{
		Name:     name,
		ParentID: parentID,
		Color:    color,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/projects", util.TodoistBaseURL)
	req, err := http.NewRequest("POST", url, strings.NewReader(string(bodyBytes)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			util.Warn("Failed to close response body", cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var project ProjectResponse
	if err := json.Unmarshal(body, &project); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &project, nil
}

var addProjectCmd = &cobra.Command{
	Use:   "project [flags] [PARENT/.../]NAME",
	Short: "Add a new project",
	Long: "Add a new project to Todoist.\n\n" +
		"`NAME` is the name of the project to create.\n" +
		"Use `PARENT/NAME` to create a project within a parent project.\n" +
		"Use `PARENT/SUBPARENT/NAME` for nested parents.\n",
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
		if projectColor != "" && !validColors[projectColor] {
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
		project, err := createProject(ConfigValue.Token, projectName, parentID, projectColor)
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
		"Use `#[PARENT/SUBPARENT.../]PROJECT` to specify the project name\n" +
		"(with optional PARENT and SUBPARENT names); note the `#` character.\n\n" +
		"Alternatively, use the --project flag to specify the project name,\n" +
		"you can omit the `#` character\n",
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
			if strings.HasPrefix(projectPath, "#") {
				projectPath = projectPath[1:]
			}
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
