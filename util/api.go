package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	TodoistURL = "https://api.todoist.com/sync/v9/sync"
)

// GetTodoistData retrieves data from the Todoist API.
//   - token: the Todoist API token
//
// Returns a pointer to a TodoistData struct with the data.
func GetTodoistData(token string) *TodoistData {
	client := &http.Client{}
	req, err := http.NewRequest("POST", TodoistURL, nil)
	if err != nil {
		Die("Failed to create request", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	q := req.URL.Query()
	q.Add("resource_types", `["projects", "sections", "items", "labels", "notes", "reminders"]`)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		Die("Failed to make request", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			Warn("Failed to close response body", cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		Die(fmt.Sprintf("Unexpected status code %d", resp.StatusCode), nil)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		Die("Failed to read response body", err)
	}

	var todoistData TodoistData
	if err := json.Unmarshal(body, &todoistData); err != nil {
		Die("Failed to unmarshal JSON", err)
	}
	return &todoistData
}
