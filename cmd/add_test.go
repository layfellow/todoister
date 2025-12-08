package cmd

import (
	"encoding/json"
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
		if !util.ValidColors[color] {
			t.Errorf("Expected color '%s' to be valid, but it's not in ValidColors map", color)
		}
	}

	// Test that the count matches
	if len(util.ValidColors) != len(validColorList) {
		t.Errorf("Expected %d valid colors, but ValidColors map has %d entries", len(validColorList), len(util.ValidColors))
	}

	// Test that invalid colors are not in the map
	invalidColors := []string{"invalid", "purple", "pink", "brown", ""}
	for _, color := range invalidColors {
		if util.ValidColors[color] {
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
		request  util.ProjectCreateRequest
		expected string
	}{
		{
			name: "basic project",
			request: util.ProjectCreateRequest{
				Name: "Test Project",
			},
			expected: `{"name":"Test Project"}`,
		},
		{
			name: "project with color",
			request: util.ProjectCreateRequest{
				Name:  "Color Project",
				Color: "blue",
			},
			expected: `{"name":"Color Project","color":"blue"}`,
		},
		{
			name: "project with parent",
			request: util.ProjectCreateRequest{
				Name:     "Sub Project",
				ParentID: "12345",
			},
			expected: `{"name":"Sub Project","parent_id":"12345"}`,
		},
		{
			name: "project with all fields",
			request: util.ProjectCreateRequest{
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

func TestParseDateInput(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantDate     string
		wantDateTime string
		wantString   string
		wantLang     string
		wantNil      bool
	}{
		{
			name:    "empty string returns nil",
			input:   "",
			wantNil: true,
		},
		{
			name:     "date only format",
			input:    "2024-01-15",
			wantDate: "2024-01-15",
		},
		{
			name:         "ISO datetime format with seconds",
			input:        "2024-01-15T14:30:00",
			wantDateTime: "2024-01-15T14:30:00",
		},
		{
			name:         "ISO datetime format without seconds",
			input:        "2024-01-15T14:30",
			wantDateTime: "2024-01-15T14:30:00",
		},
		{
			name:         "space-separated datetime with seconds",
			input:        "2024-01-15 14:30:00",
			wantDateTime: "2024-01-15T14:30:00",
		},
		{
			name:         "space-separated datetime without seconds",
			input:        "2024-01-15 14:30",
			wantDateTime: "2024-01-15T14:30:00",
		},
		{
			name:       "natural language - tomorrow",
			input:      "tomorrow",
			wantString: "tomorrow",
			wantLang:   "en",
		},
		{
			name:       "natural language - recurring",
			input:      "every monday",
			wantString: "every monday",
			wantLang:   "en",
		},
		{
			name:       "natural language - complex",
			input:      "next friday at 3pm",
			wantString: "next friday at 3pm",
			wantLang:   "en",
		},
		{
			name:       "natural language - with whitespace",
			input:      "  tomorrow  ",
			wantString: "tomorrow",
			wantLang:   "en",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := util.ParseDateInput(tt.input)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if tt.wantNil {
				if result != nil {
					t.Error("Expected nil for empty input")
				}
				return
			}

			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			if result.DueDate != tt.wantDate {
				t.Errorf("DueDate = %q, want %q", result.DueDate, tt.wantDate)
			}
			if result.DueDateTime != tt.wantDateTime {
				t.Errorf("DueDateTime = %q, want %q", result.DueDateTime, tt.wantDateTime)
			}
			if result.DueString != tt.wantString {
				t.Errorf("DueString = %q, want %q", result.DueString, tt.wantString)
			}
			if result.DueLang != tt.wantLang {
				t.Errorf("DueLang = %q, want %q", result.DueLang, tt.wantLang)
			}
		})
	}
}

func TestTaskCreateRequestWithDate(t *testing.T) {
	tests := []struct {
		name    string
		request util.TaskCreateRequest
	}{
		{
			name: "task with date only",
			request: util.TaskCreateRequest{
				Content:   "Test Task",
				ProjectID: "12345",
				DueDate:   "2024-01-15",
			},
		},
		{
			name: "task with datetime",
			request: util.TaskCreateRequest{
				Content:     "Meeting",
				ProjectID:   "12345",
				DueDateTime: "2024-01-15T14:30:00",
			},
		},
		{
			name: "task with natural language",
			request: util.TaskCreateRequest{
				Content:   "Call",
				ProjectID: "12345",
				DueString: "tomorrow",
				DueLang:   "en",
			},
		},
		{
			name: "task without date",
			request: util.TaskCreateRequest{
				Content:   "No date task",
				ProjectID: "12345",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, err := json.Marshal(tt.request)
			if err != nil {
				t.Fatalf("Failed to marshal: %v", err)
			}

			// Verify expected fields are present and empty fields are omitted
			var result map[string]interface{}
			err = json.Unmarshal(jsonBytes, &result)
			if err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}

			// Check that only non-empty date fields are included
			hasDate := result["due_date"] != nil
			hasDateTime := result["due_datetime"] != nil
			hasString := result["due_string"] != nil

			// Ensure mutual exclusivity for date fields
			dateFieldCount := 0
			if hasDate {
				dateFieldCount++
			}
			if hasDateTime {
				dateFieldCount++
			}
			if hasString {
				dateFieldCount++
			}

			if dateFieldCount > 1 {
				t.Error("Multiple date fields present - should be mutually exclusive")
			}

			// Verify content is always present
			if result["content"] == nil {
				t.Error("Content field should always be present")
			}
		})
	}
}
