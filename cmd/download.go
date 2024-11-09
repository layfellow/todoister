package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"log"
	"net/http"
	"os"
	"todoister/util"
)

func init() {
	rootCmd.AddCommand(downloadCmd)
}

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a backup in JSON format from Todoist",
	Run: func(cmd *cobra.Command, args []string) {
		println("Running download command")

		client := &http.Client{}
		req, err := http.NewRequest("POST", util.TodoistURL, nil)
		if err != nil {
			log.Fatalf("Failed to create request: %v", err)
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", util.TodoistToken))
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

		var todoistData util.TodoistResponse
		if err := json.Unmarshal(body, &todoistData); err != nil {
			log.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		outputFile, err := os.Create("todoist_data.json")
		if err != nil {
			log.Fatalf("Failed to create output file: %v", err)
		}
		defer func() {
			if cerr := outputFile.Close(); cerr != nil {
				log.Printf("Failed to close output file: %v", cerr)
			}
		}()

		encoder := json.NewEncoder(outputFile)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(todoistData); err != nil {
			log.Fatalf("Failed to write JSON to file: %v", err)
		}

		fmt.Println("Data successfully written to todoist_data.json")
	},
}
