package util

type Project struct {
	ID   string `json:"v2_id"`
	Name string `json:"name"`
}

type Section struct {
	ID        int    `json:"id"`
	ProjectID int    `json:"project_id"`
	Name      string `json:"name"`
}

type Item struct {
	ID        int    `json:"id"`
	ProjectID int    `json:"project_id"`
	SectionID int    `json:"section_id"`
	Content   string `json:"content"`
	Due       *Due   `json:"due"`
}

type Due struct {
	Date string `json:"date"`
}

type Label struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Note struct {
	ID      int    `json:"id"`
	ItemID  int    `json:"item_id"`
	Content string `json:"content"`
}

type Reminder struct {
	ID     int  `json:"id"`
	ItemID int  `json:"item_id"`
	Due    *Due `json:"due"`
}
