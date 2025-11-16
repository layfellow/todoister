package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	TodoistBaseURL = "https://api.todoist.com/api/v1"
)

// makeRequest makes an HTTP GET request to the Todoist API.
func makeRequest(token, endpoint string) ([]byte, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s%s", TodoistBaseURL, endpoint)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			Warn("Failed to close response body", cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

// GetTodoistData retrieves data from the Todoist unified API v1.
//   - token: the Todoist API token
//
// Returns a pointer to a TodoistData struct with the data.
func GetTodoistData(token string) *TodoistData {
	if token == "" {
		Die("Missing Todoist token", nil)
	}

	var todoistData TodoistData

	// Get all projects
	body, err := makeRequest(token, "/projects")
	if err != nil {
		Die("Failed to get projects", err)
	}
	if err := json.Unmarshal(body, &todoistData.Projects); err != nil {
		Die("Failed to unmarshal projects", err)
	}

	// Get all sections
	body, err = makeRequest(token, "/sections")
	if err != nil {
		Die("Failed to get sections", err)
	}
	if err := json.Unmarshal(body, &todoistData.Sections); err != nil {
		Die("Failed to unmarshal sections", err)
	}

	// Get all tasks
	body, err = makeRequest(token, "/tasks")
	if err != nil {
		Die("Failed to get tasks", err)
	}
	if err := json.Unmarshal(body, &todoistData.Items); err != nil {
		Die("Failed to unmarshal tasks", err)
	}

	// Get all labels
	body, err = makeRequest(token, "/labels")
	if err != nil {
		Die("Failed to get labels", err)
	}
	if err := json.Unmarshal(body, &todoistData.Labels); err != nil {
		Die("Failed to unmarshal labels", err)
	}

	// Get all comments (both task and project comments)
	// We need to get comments for all projects
	todoistData.Comments = make([]TodoistComment, 0)
	for _, project := range todoistData.Projects {
		body, err = makeRequest(token, fmt.Sprintf("/comments?project_id=%s", project.ID))
		if err != nil {
			Warn(fmt.Sprintf("Failed to get comments for project %s", project.ID), err)
			continue
		}
		var projectComments []TodoistComment
		if err := json.Unmarshal(body, &projectComments); err != nil {
			Warn(fmt.Sprintf("Failed to unmarshal comments for project %s", project.ID), err)
			continue
		}
		todoistData.Comments = append(todoistData.Comments, projectComments...)
	}

	return &todoistData
}
