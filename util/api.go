package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	TodoistBaseURL = "https://api.todoist.com/api/v1"
)

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Results    json.RawMessage `json:"results"`
	NextCursor string          `json:"next_cursor"`
}

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

// getAllPaginated fetches all items from an endpoint, handling both paginated and non-paginated responses
func getAllPaginated(token, endpoint string) ([]byte, error) {
	body, err := makeRequest(token, endpoint)
	if err != nil {
		return nil, err
	}

	// Try to detect if this is a paginated response or a direct array
	// Check if it starts with '[' (array) or '{' (object with pagination)
	trimmed := strings.TrimSpace(string(body))
	if len(trimmed) > 0 && trimmed[0] == '[' {
		// Direct array response, no pagination
		return body, nil
	}

	// Assume paginated response
	var allResults []json.RawMessage
	cursor := ""
	maxPages := 100 // Safety limit to prevent infinite loops

	for page := 0; page < maxPages; page++ {
		url := endpoint
		if cursor != "" {
			separator := "?"
			if strings.Contains(endpoint, "?") {
				separator = "&"
			}
			url = fmt.Sprintf("%s%scursor=%s", endpoint, separator, cursor)
		}

		if cursor != "" {
			// Fetch next page
			body, err = makeRequest(token, url)
			if err != nil {
				return nil, err
			}
		}

		var paginated PaginatedResponse
		if err := json.Unmarshal(body, &paginated); err != nil {
			// If we can't unmarshal as paginated, return what we have
			if len(allResults) == 0 {
				return nil, fmt.Errorf("failed to unmarshal response: %w", err)
			}
			break
		}

		// Check if Results array actually contains items
		var resultItems []json.RawMessage
		if err := json.Unmarshal(paginated.Results, &resultItems); err != nil {
			// Can't parse results as array
			if len(allResults) == 0 {
				return nil, fmt.Errorf("failed to unmarshal results array: %w", err)
			}
			break
		}

		if len(resultItems) > 0 {
			allResults = append(allResults, paginated.Results)
		} else {
			// Empty results array, we're done
			break
		}

		if paginated.NextCursor == "" || paginated.NextCursor == cursor {
			// No more pages or cursor didn't change (safety check)
			break
		}
		cursor = paginated.NextCursor
	}

	// Combine all results into a single array
	if len(allResults) == 0 {
		return []byte("[]"), nil
	}
	if len(allResults) == 1 {
		return allResults[0], nil
	}

	// Merge multiple result arrays
	var combined []json.RawMessage
	for _, result := range allResults {
		var items []json.RawMessage
		if err := json.Unmarshal(result, &items); err != nil {
			return nil, fmt.Errorf("failed to unmarshal results array: %w", err)
		}
		combined = append(combined, items...)
	}

	return json.Marshal(combined)
}

// GetProjects retrieves only projects data from the Todoist API.
// This is a lightweight alternative to GetTodoistData when only project information is needed.
func GetProjects(token string) []TodoistProject {
	if token == "" {
		Die("Missing Todoist token", nil)
	}

	// Get all projects
	body, err := getAllPaginated(token, "/projects")
	if err != nil {
		Die("Failed to get projects", err)
	}

	var projects []TodoistProject
	if err := json.Unmarshal(body, &projects); err != nil {
		Die("Failed to unmarshal projects", err)
	}

	return projects
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
	body, err := getAllPaginated(token, "/projects")
	if err != nil {
		Die("Failed to get projects", err)
	}
	if err := json.Unmarshal(body, &todoistData.Projects); err != nil {
		Die("Failed to unmarshal projects", err)
	}

	// Get all sections
	body, err = getAllPaginated(token, "/sections")
	if err != nil {
		Die("Failed to get sections", err)
	}
	if err := json.Unmarshal(body, &todoistData.Sections); err != nil {
		Die("Failed to unmarshal sections", err)
	}

	// Get all tasks
	body, err = getAllPaginated(token, "/tasks")
	if err != nil {
		Die("Failed to get tasks", err)
	}
	if err := json.Unmarshal(body, &todoistData.Items); err != nil {
		Die("Failed to unmarshal tasks", err)
	}

	// Get all labels
	body, err = getAllPaginated(token, "/labels")
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
		body, err = getAllPaginated(token, fmt.Sprintf("/comments?project_id=%s", project.ID))
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
