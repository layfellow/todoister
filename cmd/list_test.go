package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/layfellow/todoister/util"
)

// createTestData creates simple mock Todoist data for testing
func createTestData() *util.TodoistData {
	return &util.TodoistData{
		Projects: []util.TodoistProject{
			{Project: util.Project{Name: "Alpha"}, ID: "1", ParentID: ""},
			{Project: util.Project{Name: "Beta"}, ID: "2", ParentID: ""},
		},
		Sections: []util.TodoistSection{},
		Items: []util.TodoistItem{
			// Tasks for Beta project (ID: 2)
			{
				Task: util.Task{
					Content:    "Beta item 1",
					ChildOrder: 1,
				},
				ID:        "task1",
				ProjectID: "2",
			},
			{
				Task: util.Task{
					Content:    "Beta item 2",
					ChildOrder: 2,
				},
				ID:        "task2",
				ProjectID: "2",
			},
		},
		Labels: []util.TodoistLabel{},
	}
}

func TestListCommand(t *testing.T) {
	// Create test data with simple projects
	testData := createTestData()
	hierarchicalData := util.HierarchicalData(testData)

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Create the root project structure as the list command does
	project := util.ExportedProject{Subprojects: hierarchicalData}
	project.Name = "Projects"

	// Walk the project tree
	walkProject(&project, 0)

	// Restore stdout and read captured output
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)

	// Expected output for simple test data
	expected := `# Projects
  # Alpha
  # Beta
`

	actual := buf.String()

	if actual != expected {
		t.Errorf("List output mismatch.\nExpected:\n%s\nActual:\n%s", expected, actual)
	}
}

func TestTasksCommand(t *testing.T) {
	// Create test data with tasks for Beta
	testData := createTestData()
	hierarchicalData := util.HierarchicalData(testData)

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Create the root project structure as the tasks command does
	project := util.ExportedProject{Subprojects: hierarchicalData}
	project.Name = "Projects"

	// Get the Beta project and print its tasks
	if actualPathname, p := util.GetProjectByPathName("Beta", &project); p != nil && p.Tasks != nil {
		fmt.Printf("\n# %s\n\n", actualPathname)
		printTasks(p.Tasks)
	}

	// Restore stdout and read captured output
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)

	// Expected output for Beta project
	expected := `
# Beta

  - Beta item 1
  - Beta item 2
`

	actual := buf.String()

	if actual != expected {
		t.Errorf("Tasks output mismatch.\nExpected:\n%q\nActual:\n%q", expected, actual)
	}
}
