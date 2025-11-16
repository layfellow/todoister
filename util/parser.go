package util

import "strings"

// TodoistData as returned by the unified API v1
type TodoistData struct {
	Projects []TodoistProject `json:"projects"`
	Sections []TodoistSection `json:"sections"`
	Items    []TodoistItem    `json:"items"`
	Labels   []TodoistLabel   `json:"labels"`
	Comments []TodoistComment `json:"comments"`
}

// Projects

type Project struct {
	Name      string `json:"name"`
	Color     string `json:"color"`
	ViewStyle string `json:"view_style"`
}

type TodoistProject struct {
	Project
	ID       string `json:"id"`
	ParentID string `json:"parent_id"`
}

type ExportedProject struct {
	Project
	Subprojects []*ExportedProject `json:"subprojects"`
	Sections    []*ExportedSection `json:"sections"`
	Tasks       []*ExportedTask    `json:"tasks"`
	Comments    []*ExportedComment `json:"comments"`
}

// Sections

type Section struct {
	Name      string `json:"name"`
	Collapsed bool   `json:"collapsed"`
}

type TodoistSection struct {
	Section
	ID        string `json:"id"`
	ProjectID string `json:"project_id"`
	Order     int    `json:"order"`
}

type ExportedSection struct {
	Section
	Tasks []*ExportedTask `json:"tasks"`
}

// Tasks (aka Items)

type Task struct {
	Content     string `json:"content"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
	ChildOrder  int    `json:"child_order"`
	Collapsed   bool   `json:"collapsed"`
	CompletedAt string `json:"completed_at"`
}

type TodoistItem struct {
	Task
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	SectionID string    `json:"section_id"`
	Labels    []string  `json:"labels"`
	Duration  *Duration `json:"duration"`
	Due       *Due      `json:"due"`
}

type ExportedTask struct {
	Task
	Labeled  []*ExportedLabel   `json:"labeled"`
	Comments []*ExportedComment `json:"comments"`
	Duration *Duration          `json:"duration"`
	Due      *Due               `json:"due"`
}

// Labels

type Label struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type TodoistLabel struct {
	Label
	ID string `json:"id"`
}

type ExportedLabel struct {
	Label
}

// Comments

type Comment struct {
	Content string `json:"content"`
}

type TodoistComment struct {
	Comment
	ID        string `json:"id"`
	TaskID    string `json:"task_id"`
	ProjectID string `json:"project_id"`
}

type ExportedComment struct {
	Comment
}

// Due dates

type Due struct {
	Date        string `json:"date"`
	IsRecurring bool   `json:"is_recurring"`
	String      string `json:"string"`
	Datetime    string `json:"datetime"`
	Timezone    string `json:"timezone"`
}

// Duration

type Duration struct {
	Amount int    `json:"amount"`
	Unit   string `json:"unit"`
}

// GetTasksByProjectID returns a slice of tasks for a given project ID.
//   - projectID: the project ID
//   - todoistData: pointer to TodoistData struct as returned by the API
//
// Returns a pointer to a slice of TodoistItem structs (aka Tasks).
func GetTasksByProjectID(projectID string, todoistData *TodoistData) *[]TodoistItem {
	tasks := make([]TodoistItem, 0)
	for _, item := range todoistData.Items {
		if item.ProjectID == projectID {
			tasks = append(tasks, item)
		}
	}
	return &tasks
}

// GetProjectIDByName returns the project ID for a given project name.
//   - name: the project name
//   - todoistData: pointer to TodoistData struct as returned by the API
//
// Returns the project ID as a string.
func GetProjectIDByName(name string, todoistData *TodoistData) string {
	for _, project := range todoistData.Projects {
		if strings.EqualFold(project.Name, name) {
			return project.ID
		}
	}
	return ""
}

// GetProjectIDByPathFromProjects returns the Todoist project ID for a given project path using only projects data.
// This is a lightweight alternative to GetProjectIDByPath that doesn't require the full TodoistData structure.
//   - pathname: the project pathname as entered by the user (e.g., "Work/Reports")
//   - projects: slice of TodoistProject structs
//
// Returns the project ID as a string, or empty string if not found.
func GetProjectIDByPathFromProjects(pathname string, projects []TodoistProject) string {
	names := strings.Split(pathname, "/")
	var currentParentID string

	for i, name := range names {
		found := false
		for _, proj := range projects {
			// Match by name and parent_id to ensure we traverse the correct path
			if strings.EqualFold(proj.Name, name) && proj.ParentID == currentParentID {
				currentParentID = proj.ID
				found = true
				// If this is the last name in the path, we've found our project
				if i == len(names)-1 {
					return proj.ID
				}
				break
			}
		}
		if !found {
			return ""
		}
	}

	return ""
}

// GetProjectByPathName returns a project and its canonical pathname.
//   - pathname: the project pathname as entered by the user
//   - project: pointer to ExportedProject struct as parsed by HierarchicalData
//
// Returns the canonical pathname and a pointer to the projectʼs ExportedProject struct.
func GetProjectByPathName(pathname string, project *ExportedProject) (string, *ExportedProject) {
	var p *ExportedProject
	var actualName, actualPathname string

	root := project
	names := strings.Split(pathname, "/")
	for _, name := range names {
		if actualName, p = GetProjectByName(name, root); p != nil {
			root = p
			actualPathname += actualName + "/"
		}
	}
	// Trim the trailing slash
	if len(actualPathname) > 0 {
		actualPathname = actualPathname[:len(actualPathname)-1]
	}
	return actualPathname, p
}

// GetProjectIDByPath returns the Todoist project ID for a given project path.
//   - pathname: the project pathname as entered by the user (e.g., "Work/Reports")
//   - todoistData: pointer to TodoistData struct
//
// Returns the project ID as a string, or empty string if not found.
func GetProjectIDByPath(pathname string, todoistData *TodoistData) (string, string) {
	// Build hierarchical structure
	hierarchicalData := HierarchicalData(todoistData)
	root := ExportedProject{Subprojects: hierarchicalData}
	root.Name = "Projects"

	// Find the project in the hierarchy
	actualPathname, exportedProject := GetProjectByPathName(pathname, &root)
	if exportedProject == nil {
		return "", ""
	}

	// Now find the corresponding TodoistProject by matching the name
	// We need to search through the hierarchy to find the right one
	// Build a path through the hierarchy to ensure we get the correct project
	names := strings.Split(pathname, "/")
	var currentParentID string

	for i, name := range names {
		found := false
		for _, proj := range todoistData.Projects {
			// Match by name and parent_id to ensure we traverse the correct path
			if strings.EqualFold(proj.Name, name) && proj.ParentID == currentParentID {
				currentParentID = proj.ID
				found = true
				// If this is the last name in the path, we've found our project
				if i == len(names)-1 {
					return proj.ID, actualPathname
				}
				break
			}
		}
		if !found {
			return "", ""
		}
	}

	return "", ""
}

// GetProjectByName returns a project by name.
// - name: the project name as entered by the user
// - project: pointer to ExportedProject struct as parsed by HierarchicalData
//
// Returns the canonical project name and a pointer to the projectʼs ExportedProject struct.
func GetProjectByName(name string, project *ExportedProject) (string, *ExportedProject) {
	if strings.EqualFold(project.Name, name) {
		return project.Name, project
	} else {
		if project.Subprojects != nil {
			for _, subproject := range project.Subprojects {
				if actualName, p := GetProjectByName(name, subproject); p != nil {
					return actualName, p
				}
			}
		}
	}
	return "", nil
}

// HierarchicalData converts TodoistData to a hierarchical structure of Exported* structs.
//   - todoistData: a pointer to a TodoistData struct
//
// Returns a slice of root ExportedProject pointers.
func HierarchicalData(todoistData *TodoistData) []*ExportedProject {
	// Persistent variable to hold the root ExportedProject references.
	var roots []*ExportedProject

	todoistProjects := todoistData.Projects

	// Map to hold references to each project by ID for easy lookup.
	var projectMap = make(map[string]*ExportedProject)

	// Initialize the projectMap with common Project fields and empty Subproject slices.
	for _, project := range todoistProjects {
		p := new(ExportedProject)
		// Copy common fields from TodoistProject to ExportedProject.
		p.Project = project.Project
		p.Subprojects = make([]*ExportedProject, 0)
		p.Sections = make([]*ExportedSection, 0)
		p.Tasks = make([]*ExportedTask, 0)
		projectMap[project.ID] = p
	}

	// Build the hierarchy by linking Subprojects to their parent Projects.
	for _, project := range todoistProjects {
		if project.ParentID == "" {
			// If there's no ParentID, it's a root project.
			roots = append(roots, projectMap[project.ID])
		} else {
			// Otherwise, it's a Subproject, so add it to its parent's Subprojects slice.
			projectMap[project.ParentID].Subprojects =
				append(projectMap[project.ParentID].Subprojects, projectMap[project.ID])
		}
	}

	todoistSections := todoistData.Sections

	// Map to hold references to each section by ID for easy lookup.
	var sectionMap = make(map[string]*ExportedSection)

	// Initialize sectionMap with common Section fields and empty Task slices.
	for _, section := range todoistSections {
		s := new(ExportedSection)
		// Copy common fields from TodoistSection to ExportedSection.
		s.Section = section.Section
		s.Tasks = make([]*ExportedTask, 0)
		sectionMap[section.ID] = s
	}

	// Add to the hierarchy by linking Sections to their parent Projects.
	for _, section := range todoistSections {
		projectMap[section.ProjectID].Sections =
			append(projectMap[section.ProjectID].Sections, sectionMap[section.ID])
	}

	todoistItems := todoistData.Items

	// Map to hold references to each item (task) by ID for easy lookup.
	var taskMap = make(map[string]*ExportedTask)

	// Initialize taskMap with common Task fields and empty Task slices.
	for _, item := range todoistItems {
		t := new(ExportedTask)
		t.Task = item.Task // Copy common fields from TodoistItem to ExportedTask

		if item.Duration != nil && item.Duration.Amount > 0 {
			t.Duration = new(Duration)
			// Copy common fields from duration to ExportedTask.
			*t.Duration = *item.Duration
		}
		if item.Due != nil && item.Due.Date != "" {
			t.Due = new(Due)
			// Copy common fields from due date to ExportedTask.
			*t.Due = *item.Due
		}

		t.Labeled = make([]*ExportedLabel, 0)
		t.Comments = make([]*ExportedComment, 0)
		taskMap[item.ID] = t
	}

	// Add to the hierarchy by linking Tasks to their parent Projects or Sections.
	for _, item := range todoistItems {
		if item.SectionID == "" {
			// If there's no SectionID, it's a task attached to the project.
			projectMap[item.ProjectID].Tasks =
				append(projectMap[item.ProjectID].Tasks, taskMap[item.ID])
		} else {
			// Otherwise, it's attached to a section, so add it to the section's Tasks slice.
			sectionMap[item.SectionID].Tasks =
				append(sectionMap[item.SectionID].Tasks, taskMap[item.ID])
		}
	}

	todoistLabels := todoistData.Labels

	// Map to hold references to each label by ID for easy lookup.
	var labelMap = make(map[string]*ExportedLabel)

	// Initialize labelMap with common Label fields.
	for _, label := range todoistLabels {
		l := new(ExportedLabel)
		// Copy common fields from TodoistLabel to ExportedLabel
		l.Label = label.Label
		labelMap[label.Name] = l
	}

	// Add to the hierarchy by linking Labels to their respective Tasks.
	for _, item := range todoistItems {
		for _, label := range item.Labels {
			if label != "" {
				taskMap[item.ID].Labeled = append(taskMap[item.ID].Labeled, labelMap[label])
			}
		}
	}

	todoistComments := todoistData.Comments

	// Map to hold references to each comment by ID for easy lookup.
	var commentMap = make(map[string]*ExportedComment)

	// Initialize commentMap with common Comment fields.
	for _, comment := range todoistComments {
		c := new(ExportedComment)
		// Copy common fields from TodoistComment to ExportedComment
		c.Comment = comment.Comment
		commentMap[comment.ID] = c
	}

	// Add to the hierarchy by linking Comments to their respective Tasks or Projects.
	for _, comment := range todoistComments {
		if comment.TaskID != "" {
			// Task comment
			if task, exists := taskMap[comment.TaskID]; exists {
				task.Comments = append(task.Comments, commentMap[comment.ID])
			}
		} else if comment.ProjectID != "" {
			// Project comment
			if project, exists := projectMap[comment.ProjectID]; exists {
				project.Comments = append(project.Comments, commentMap[comment.ID])
			}
		}
	}

	return roots
}
