package cmd

import (
	"testing"

	"github.com/layfellow/todoister/util"
)

func TestDeleteProjectPathParsing(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		projects    []util.TodoistProject
		expectedID  string
		expectFound bool
	}{
		{
			name: "root level project",
			path: "Shopping",
			projects: []util.TodoistProject{
				{Project: util.Project{Name: "Shopping"}, ID: "123", ParentID: ""},
			},
			expectedID:  "123",
			expectFound: true,
		},
		{
			name: "nested project",
			path: "Work/Reports",
			projects: []util.TodoistProject{
				{Project: util.Project{Name: "Work"}, ID: "100", ParentID: ""},
				{Project: util.Project{Name: "Reports"}, ID: "200", ParentID: "100"},
			},
			expectedID:  "200",
			expectFound: true,
		},
		{
			name: "deeply nested project",
			path: "Work/Projects/Q1",
			projects: []util.TodoistProject{
				{Project: util.Project{Name: "Work"}, ID: "100", ParentID: ""},
				{Project: util.Project{Name: "Projects"}, ID: "200", ParentID: "100"},
				{Project: util.Project{Name: "Q1"}, ID: "300", ParentID: "200"},
			},
			expectedID:  "300",
			expectFound: true,
		},
		{
			name: "project not found",
			path: "NonExistent",
			projects: []util.TodoistProject{
				{Project: util.Project{Name: "Shopping"}, ID: "123", ParentID: ""},
			},
			expectedID:  "",
			expectFound: false,
		},
		{
			name: "parent not found",
			path: "NonExistent/Reports",
			projects: []util.TodoistProject{
				{Project: util.Project{Name: "Work"}, ID: "100", ParentID: ""},
				{Project: util.Project{Name: "Reports"}, ID: "200", ParentID: "100"},
			},
			expectedID:  "",
			expectFound: false,
		},
		{
			name: "case insensitive match",
			path: "SHOPPING",
			projects: []util.TodoistProject{
				{Project: util.Project{Name: "Shopping"}, ID: "123", ParentID: ""},
			},
			expectedID:  "123",
			expectFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.GetProjectIDByPathFromProjects(tt.path, tt.projects)
			if tt.expectFound {
				if result != tt.expectedID {
					t.Errorf("Expected project ID '%s', got '%s'", tt.expectedID, result)
				}
			} else {
				if result != "" {
					t.Errorf("Expected empty project ID, got '%s'", result)
				}
			}
		})
	}
}

func TestDeleteProjectForceFlagDefault(t *testing.T) {
	// Verify the force flag defaults to false
	if forceDelete {
		t.Error("Expected forceDelete to default to false")
	}
}

func TestDeleteTaskFindByPrefix(t *testing.T) {
	tests := []struct {
		name          string
		projectID     string
		prefix        string
		items         []util.TodoistItem
		expectedCount int
	}{
		{
			name:      "single match",
			projectID: "100",
			prefix:    "Buy milk",
			items: []util.TodoistItem{
				{Task: util.Task{Content: "Buy milk"}, ID: "1", ProjectID: "100"},
				{Task: util.Task{Content: "Buy groceries"}, ID: "2", ProjectID: "100"},
			},
			expectedCount: 1,
		},
		{
			name:      "multiple matches",
			projectID: "100",
			prefix:    "Buy",
			items: []util.TodoistItem{
				{Task: util.Task{Content: "Buy milk"}, ID: "1", ProjectID: "100"},
				{Task: util.Task{Content: "Buy groceries"}, ID: "2", ProjectID: "100"},
				{Task: util.Task{Content: "Buy coffee"}, ID: "3", ProjectID: "100"},
			},
			expectedCount: 3,
		},
		{
			name:      "no match",
			projectID: "100",
			prefix:    "Sell",
			items: []util.TodoistItem{
				{Task: util.Task{Content: "Buy milk"}, ID: "1", ProjectID: "100"},
			},
			expectedCount: 0,
		},
		{
			name:      "case insensitive match",
			projectID: "100",
			prefix:    "BUY MILK",
			items: []util.TodoistItem{
				{Task: util.Task{Content: "Buy milk"}, ID: "1", ProjectID: "100"},
			},
			expectedCount: 1,
		},
		{
			name:      "different project",
			projectID: "100",
			prefix:    "Buy",
			items: []util.TodoistItem{
				{Task: util.Task{Content: "Buy milk"}, ID: "1", ProjectID: "200"},
			},
			expectedCount: 0,
		},
		{
			name:      "completed task excluded",
			projectID: "100",
			prefix:    "Buy",
			items: []util.TodoistItem{
				{Task: util.Task{Content: "Buy milk", CompletedAt: "2024-01-01"}, ID: "1", ProjectID: "100"},
			},
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todoistData := &util.TodoistData{
				Items: tt.items,
			}
			matches := util.FindTasksByPrefix(tt.projectID, tt.prefix, todoistData)
			if len(matches) != tt.expectedCount {
				t.Errorf("Expected %d matches, got %d", tt.expectedCount, len(matches))
			}
		})
	}
}
