package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	// Get version from environment variable
	Version = os.Getenv("VERSION")
	if Version == "" {
		t.Fatal("VERSION environment variable is not set")
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Reset root command state
	RootCmd.SetArgs([]string{"version"})

	// Execute the command
	err := RootCmd.Execute()
	if err != nil {
		t.Fatalf("Failed to execute version command: %v", err)
	}

	// Restore stdout and read captured output
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)

	// Check the output
	expected := fmt.Sprintf("%s v%s\n", VersionText, Version)
	actual := buf.String()

	if actual != expected {
		t.Errorf("Version output mismatch.\nExpected: %q\nActual: %q", expected, actual)
	}
}
