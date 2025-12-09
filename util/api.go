package util

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	TodoistBaseURL = "https://api.todoist.com/api/v1"
	TodoistSyncURL = "https://api.todoist.com/api/v1/sync"
)

// ValidColors are the allowed color values for projects
var ValidColors = map[string]bool{
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

// SyncResponse represents the Sync API response
type SyncResponse struct {
	SyncToken    string            `json:"sync_token"`
	FullSync     bool              `json:"full_sync"`
	Projects     []TodoistProject  `json:"projects"`
	Sections     []TodoistSection  `json:"sections"`
	Items        []TodoistItem     `json:"items"`
	Labels       []TodoistLabel    `json:"labels"`
	Notes        []TodoistComment  `json:"notes"`
	ProjectNotes []TodoistComment  `json:"project_notes"`
}

// makeSyncRequest makes a POST request to the Sync API endpoint.
//   - token: the Todoist API token
//   - syncToken: the sync token for incremental sync, or "*" for full sync
//   - resourceTypes: the list of resource types to fetch
//
// Returns a SyncResponse and an error if the request fails.
func makeSyncRequest(token, syncToken string, resourceTypes []string) (*SyncResponse, error) {
	client := &http.Client{}

	// Construct proper JSON array for resource_types
	resourceTypesJSON, err := json.Marshal(resourceTypes)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resource types: %w", err)
	}
	formData := fmt.Sprintf("sync_token=%s&resource_types=%s", syncToken, string(resourceTypesJSON))

	req, err := http.NewRequest("POST", TodoistSyncURL, strings.NewReader(formData))
	if err != nil {
		return nil, fmt.Errorf("failed to create sync request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make sync request: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			Warn("Failed to close response body", cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("sync request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read sync response body: %w", err)
	}

	var syncResp SyncResponse
	if err := json.Unmarshal(body, &syncResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal sync response: %w", err)
	}

	return &syncResp, nil
}

// GetTodoistData retrieves data from the Todoist Sync API with caching.
//   - token: the Todoist API token
//
// Returns a pointer to a TodoistData struct with the data.
// Uses a local Protobuf cache for performance. On first run, performs a full sync.
// On subsequent runs, performs incremental sync to fetch only changes.
func GetTodoistData(token string) *TodoistData {
	if token == "" {
		Die("Missing Todoist token", nil)
	}

	// 1. Try to load cache
	cached, err := LoadCache()
	if err != nil {
		Warn("Failed to load cache, will perform full sync", err)
	}

	// 2. Determine sync token
	syncToken := "*" // Full sync by default
	if cached != nil && cached.SyncToken != "" {
		syncToken = cached.SyncToken
	}

	// 3. Make Sync API request
	resourceTypes := []string{"projects", "sections", "items", "labels", "notes", "project_notes"}
	syncResp, err := makeSyncRequest(token, syncToken, resourceTypes)
	if err != nil {
		// If we have cached data and network fails, warn and use cache
		if cached != nil {
			Warn("Failed to sync, using cached data", err)
			return convertCachedToTodoistData(cached)
		}
		Die("Failed to sync", err)
	}

	// 4. Merge or replace data
	var todoistData *TodoistData
	if syncResp.FullSync || cached == nil {
		// Full sync: use response directly
		todoistData = &TodoistData{
			Projects: syncResp.Projects,
			Sections: syncResp.Sections,
			Items:    syncResp.Items,
			Labels:   syncResp.Labels,
			Comments: append(syncResp.Notes, syncResp.ProjectNotes...),
		}
	} else {
		// Incremental sync: merge with cached data
		cachedData := convertCachedToTodoistData(cached)
		todoistData = mergeData(cachedData, syncResp)
	}

	// 5. Update cache
	newCache := convertTodoistDataToCached(todoistData, syncResp.SyncToken)
	if err := SaveCache(newCache); err != nil {
		Warn("Failed to save cache", err)
		// Continue anyway - not fatal
	}

	// 6. Return data
	return todoistData
}


// TaskResponse represents a task creation response from the API
type TaskResponse struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	ProjectID string `json:"project_id"`
}

// DateParams holds the parsed date parameters for task creation
type DateParams struct {
	DueDate     string // YYYY-MM-DD format
	DueDateTime string // YYYY-MM-DDTHH:MM:SS format (no timezone)
	DueString   string // Natural language string
	DueLang     string // Language code, default "en"
}

// TaskCreateRequest represents the request body for creating a task via REST API
type TaskCreateRequest struct {
	Content     string `json:"content"`
	ProjectID   string `json:"project_id,omitempty"`
	DueDate     string `json:"due_date,omitempty"`
	DueDateTime string `json:"due_datetime,omitempty"`
	DueString   string `json:"due_string,omitempty"`
	DueLang     string `json:"due_lang,omitempty"`
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

// ParseDateInput analyzes a date string and returns appropriate API parameters
// Returns nil if input is empty, indicating no due date
func ParseDateInput(dateInput string) (*DateParams, error) {
	dateInput = strings.TrimSpace(dateInput)

	if dateInput == "" {
		return nil, nil
	}

	// Try date-only format (YYYY-MM-DD)
	if t, err := time.Parse("2006-01-02", dateInput); err == nil {
		return &DateParams{DueDate: t.Format("2006-01-02")}, nil
	}

	// Try datetime formats (in order of specificity)
	dateTimeFormats := []string{
		"2006-01-02T15:04:05", // ISO with seconds
		"2006-01-02T15:04",    // ISO without seconds
		"2006-01-02 15:04:05", // Space-separated with seconds
		"2006-01-02 15:04",    // Space-separated without seconds
	}

	for _, format := range dateTimeFormats {
		if t, err := time.Parse(format, dateInput); err == nil {
			// Normalize to YYYY-MM-DDTHH:MM:SS (no timezone)
			normalized := t.Format("2006-01-02T15:04:05")
			return &DateParams{DueDateTime: normalized}, nil
		}
	}

	// Fallback to natural language
	return &DateParams{
		DueString: dateInput,
		DueLang:   "en",
	}, nil
}

// CreateTask makes a POST request to create a new task using the REST API v1
func CreateTask(token, content, projectID string, dateParams *DateParams) (*TaskResponse, error) {
	client := &http.Client{}

	reqBody := TaskCreateRequest{
		Content:   content,
		ProjectID: projectID,
	}

	// Add due date parameters if provided
	if dateParams != nil {
		reqBody.DueDate = dateParams.DueDate
		reqBody.DueDateTime = dateParams.DueDateTime
		reqBody.DueString = dateParams.DueString
		reqBody.DueLang = dateParams.DueLang
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/tasks", TodoistBaseURL)
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
			Warn("Failed to close response body", cerr)
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

	var task TaskResponse
	if err := json.Unmarshal(body, &task); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &task, nil
}

// CreateProject makes a POST request to create a new project
func CreateProject(token, name, parentID, color string) (*ProjectResponse, error) {
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

	url := fmt.Sprintf("%s/projects", TodoistBaseURL)
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
			Warn("Failed to close response body", cerr)
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

// CompleteTask closes/completes a task using the Sync API item_close command.
//   - token: Todoist API token
//   - taskID: The task ID to close
//
// Returns an error if the request fails.
func CompleteTask(token, taskID string) error {
	client := &http.Client{}

	// Generate UUID for command
	uuid := generateUUID()

	// Build Sync API command
	commands := []map[string]interface{}{
		{
			"type": "item_close",
			"uuid": uuid,
			"args": map[string]string{
				"id": taskID,
			},
		},
	}

	commandsJSON, err := json.Marshal(commands)
	if err != nil {
		return fmt.Errorf("failed to marshal commands: %w", err)
	}

	// Build form data
	formData := fmt.Sprintf("commands=%s", url.QueryEscape(string(commandsJSON)))

	req, err := http.NewRequest("POST", TodoistSyncURL, strings.NewReader(formData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to complete task: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			Warn("Failed to close response body", cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// generateUUID generates a simple UUID for Sync API commands
func generateUUID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// crypto/rand.Read should never fail, but handle it anyway
		Warn("Failed to generate UUID, using fallback", err)
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// DeleteProject deletes a project using the Sync API project_delete command.
//   - token: Todoist API token
//   - projectID: The project ID to delete
//
// Returns an error if the request fails.
// Note: This deletes the project and all its descendants.
func DeleteProject(token, projectID string) error {
	client := &http.Client{}

	// Generate UUID for command
	uuid := generateUUID()

	// Build Sync API command
	commands := []map[string]interface{}{
		{
			"type": "project_delete",
			"uuid": uuid,
			"args": map[string]string{
				"id": projectID,
			},
		},
	}

	commandsJSON, err := json.Marshal(commands)
	if err != nil {
		return fmt.Errorf("failed to marshal commands: %w", err)
	}

	// Build form data
	formData := fmt.Sprintf("commands=%s", url.QueryEscape(string(commandsJSON)))

	req, err := http.NewRequest("POST", TodoistSyncURL, strings.NewReader(formData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			Warn("Failed to close response body", cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
