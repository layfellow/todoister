package util

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	TodoistURL = "https://api.todoist.com/sync/v9/sync"
)

var TodoistToken string

func GetTodoistData() (*TodoistData, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", TodoistURL, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", TodoistToken))
	req.Header.Set("Content-Type", "application/json")

	q := req.URL.Query()
	q.Add("resource_types", `["projects", "sections", "items", "labels", "notes", "reminders"]`)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to make request: %v", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("Failed to close response body: %v", cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	var todoistData TodoistData
	if err := json.Unmarshal(body, &todoistData); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	return &todoistData, nil
}
