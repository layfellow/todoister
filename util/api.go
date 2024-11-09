package util

const (
	TodoistURL = "https://api.todoist.com/sync/v9/sync"
)

var TodoistToken string

type TodoistResponse struct {
	Projects []Project `json:"projects"`
	//Sections  []Section  `json:"sections"`
	//Items     []Item     `json:"items"`
	//Labels    []Label    `json:"labels"`
	//Notes     []Note     `json:"notes"`
	//Reminders []Reminder `json:"reminders"`
}
