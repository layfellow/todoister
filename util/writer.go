package util

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func WriteHierarchicalData(roots []*ExportedProject, fileName string) error {

	outputFile, err := os.Create(fileName)
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

	var projects = make(map[string][]*ExportedProject)
	projects["projects"] = roots

	if err := encoder.Encode(projects); err != nil {
		log.Fatalf("Failed to write JSON to file: %v", err)
	}

	fmt.Printf("Successfully written %s\n", fileName)
	return nil
}
