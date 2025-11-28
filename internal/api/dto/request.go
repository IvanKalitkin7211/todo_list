package dto

type TaskRequestDTO struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Status   string `json:"status,omitempty"`
	Priority string `json:"priority,omitempty"`
	DueDate  string `json:"due_date,omitempty"`
}
