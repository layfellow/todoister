package util

// TodoistData as returned by the API
type TodoistData struct {
	Projects     []TodoistProject     `json:"projects"`
	Sections     []TodoistSection     `json:"sections"`
	Items        []TodoistItem        `json:"items"`
	Labels       []TodoistLabel       `json:"labels"`
	Notes        []TodoistNote        `json:"notes"`
	ProjectNotes []TodoistProjectNote `json:"project_notes"`
	Reminders    []TodoistReminder    `json:"reminders"`
}

// Projects

type Project struct {
	Name      string `json:"name"`
	Color     string `json:"color"`
	ViewStyle string `json:"view_style"`
}

type TodoistProject struct {
	Project
	ID         string `json:"v2_id"`
	ParentID   string `json:"v2_parent_id"`
	IsArchived bool   `json:"is_archived"`
	IsDeleted  bool   `json:"is_deleted"`
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
	ID           string `json:"v2_id"`
	ProjectID    string `json:"v2_project_id"`
	SectionOrder int    `json:"section_order"`
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
	ID        string   `json:"v2_id"`
	ProjectID string   `json:"v2_project_id"`
	SectionID string   `json:"v2_section_id"`
	Labels    []string `json:"labels"`
	Duration  Duration `json:"duration"`
	Due       Due      `json:"due"`
}

type ExportedTask struct {
	Task
	Labeled   []*ExportedLabel    `json:"labeled"`
	Comments  []*ExportedComment  `json:"comments"`
	Reminders []*ExportedReminder `json:"reminders"`
	Duration  *Duration           `json:"duration"`
	Due       *Due                `json:"due"`
}

// Labels

type Label struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type TodoistLabel struct {
	Label
	ID        string `json:"id"`
	IsDeleted bool   `json:"is_deleted"`
}

type ExportedLabel struct {
	Label
}

// Comments (aka Notes)

type Comment struct {
	Content   string `json:"content"`
	IsDeleted bool   `json:"is_deleted"`
}

type TodoistNote struct {
	Comment
	ID        string `json:"v2_id"`
	ItemID    string `json:"v2_item_id"`
	ProjectID string `json:"v2_project_id"`
}

type TodoistProjectNote struct {
	Comment
	ID        string `json:"v2_id"`
	ProjectID string `json:"v2_project_id"`
}

type ExportedComment struct {
	Comment
}

// Reminders

type Reminder struct {
	Type      string `json:"type"`
	IsDeleted bool   `json:"is_deleted"`
}

type TodoistReminder struct {
	Reminder
	ID     string `json:"v2_id"`
	ItemID string `json:"v2_item_id"`
	Due    Due    `json:"due"`
}

type ExportedReminder struct {
	Reminder
	Due *Due `json:"due"`
}

// Due dates

type Due struct {
	Date        string `json:"date"`
	IsRecurring bool   `json:"is_recurring"`
	Lang        string `json:"lang"`
	String      string `json:"string"`
	Timezone    string `json:"timezone"`
}

// Duration

type Duration struct {
	Amount int    `json:"amount"`
	Unit   string `json:"unit"`
}

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

		if item.Duration.Amount > 0 {
			t.Duration = new(Duration)
			// Copy common fields from duration to ExportedTask.
			*t.Duration = item.Duration
		}
		if item.Due.Date != "" {
			t.Due = new(Due)
			// Copy common fields from due date to ExportedTask.
			*t.Due = item.Due
		}

		t.Labeled = make([]*ExportedLabel, 0)
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

	todoistNotes := todoistData.Notes

	// Map to hold references to each note (comment) by ID for easy lookup.
	var commentMap = make(map[string]*ExportedComment)

	// Initialize commentMap with common Comment fields.
	for _, note := range todoistNotes {
		c := new(ExportedComment)
		// Copy common fields from TodoistNote to ExportedComment
		c.Comment = note.Comment
		commentMap[note.ID] = c
	}

	// Add to the hierarchy by linking Comments to their respective Tasks.
	for _, note := range todoistNotes {
		taskMap[note.ItemID].Comments =
			append(taskMap[note.ItemID].Comments, commentMap[note.ID])
	}

	todoistProjectNotes := todoistData.ProjectNotes

	// Map to hold references to each project note (comment) by ID for easy lookup.
	var projectCommentMap = make(map[string]*ExportedComment)

	// Initialize projectCommentMap with common Comment fields.
	for _, note := range todoistProjectNotes {
		c := new(ExportedComment)
		// Copy common fields from TodoistProjectNote to ExportedComment
		c.Comment = note.Comment
		projectCommentMap[note.ID] = c
	}

	// Add to the hierarchy by linking Comments to their respective Projects.
	for _, note := range todoistProjectNotes {
		projectMap[note.ProjectID].Comments =
			append(projectMap[note.ProjectID].Comments, projectCommentMap[note.ID])
	}

	todoistReminders := todoistData.Reminders

	// Map to hold references to each reminder by ID for easy lookup.
	var reminderMap = make(map[string]*ExportedReminder)

	// Initialize reminderMap with common Reminder fields.
	for _, reminder := range todoistReminders {
		r := new(ExportedReminder)
		// Copy common fields from TodoistReminder to ExportedReminder
		r.Reminder = reminder.Reminder

		if reminder.Due.Date != "" {
			r.Due = new(Due)
			// Copy common fields from due date to ExportedReminder
			*r.Due = reminder.Due
		}
		reminderMap[reminder.ID] = r
	}

	// Add to the hierarchy by linking Reminders to their respective Tasks.
	for _, reminder := range todoistReminders {
		taskMap[reminder.ItemID].Reminders =
			append(taskMap[reminder.ItemID].Reminders, reminderMap[reminder.ID])
	}

	return roots
}
