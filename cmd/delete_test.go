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
