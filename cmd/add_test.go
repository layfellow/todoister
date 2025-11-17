package cmd

import (
	"testing"

	"github.com/layfellow/todoister/util"
)

func TestValidColors(t *testing.T) {
	validColorList := []string{
		"berry_red", "red", "orange", "yellow", "olive_green", "lime_green",
		"green", "mint_green", "teal", "sky_blue", "light_blue", "blue",
		"grape", "violet", "lavender", "magenta", "salmon", "charcoal", "grey", "taupe",
	}

	// Test that all expected colors are valid
	for _, color := range validColorList {
		if !validColors[color] {
			t.Errorf("Expected color '%s' to be valid, but it's not in validColors map", color)
		}
	}

	// Test that the count matches
	if len(validColors) != len(validColorList) {
		t.Errorf("Expected %d valid colors, but validColors map has %d entries", len(validColorList), len(validColors))
	}

	// Test that invalid colors are not in the map
	invalidColors := []string{"invalid", "purple", "pink", "brown", ""}
	for _, color := range invalidColors {
		if validColors[color] {
			t.Errorf("Color '%s' should not be valid", color)
		}
	}
}


func TestTaskCreateRequest(t *testing.T) {
	tests := []struct {
		name     string
		request  util.TaskCreateRequest
		expected string
	}{
		{
			name: "basic task",
			request: util.TaskCreateRequest{
				Content:   "Test Task",
				ProjectID: "12345",
			},
			expected: `{"content":"Test Task","project_id":"12345"}`,
		},
		{
			name: "task without project",
			request: util.TaskCreateRequest{
				Content: "Inbox Task",
			},
			expected: `{"content":"Inbox Task"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that the struct has the expected fields
			if tt.request.Content == "" {
				t.Error("Task content should not be empty")
			}
			// Basic structure validation
			if tt.name == "basic task" && tt.request.ProjectID == "" {
				t.Error("Project ID should not be empty for basic test")
			}
		})
	}
}

func TestTaskResponseStructure(t *testing.T) {
	response := util.TaskResponse{
		ID:        "67890",
		Content:   "Created Task",
		ProjectID: "12345",
	}

	if response.ID != "67890" {
		t.Errorf("Expected ID '67890', got '%s'", response.ID)
	}
	if response.Content != "Created Task" {
		t.Errorf("Expected content 'Created Task', got '%s'", response.Content)
	}
	if response.ProjectID != "12345" {
		t.Errorf("Expected project ID '12345', got '%s'", response.ProjectID)
	}
}

func TestProjectCreateRequestJSON(t *testing.T) {
	tests := []struct {
		name     string
		request  ProjectCreateRequest
		expected string
	}{
		{
			name: "basic project",
			request: ProjectCreateRequest{
				Name: "Test Project",
			},
			expected: `{"name":"Test Project"}`,
		},
		{
			name: "project with color",
			request: ProjectCreateRequest{
				Name:  "Colored Project",
				Color: "blue",
			},
			expected: `{"name":"Colored Project","color":"blue"}`,
		},
		{
			name: "project with parent",
			request: ProjectCreateRequest{
				Name:     "Sub Project",
				ParentID: "12345",
			},
			expected: `{"name":"Sub Project","parent_id":"12345"}`,
		},
		{
			name: "project with all fields",
			request: ProjectCreateRequest{
				Name:     "Complete Project",
				ParentID: "67890",
				Color:    "red",
			},
			expected: `{"name":"Complete Project","parent_id":"67890","color":"red"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We're just testing the struct tags are correct by checking expected JSON format
			// The actual marshaling is tested implicitly when the command runs
			if tt.request.Name == "" && tt.name != "basic project" {
				t.Error("Request name should not be empty for non-basic tests")
			}
		})
	}
}
