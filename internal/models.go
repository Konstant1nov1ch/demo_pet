package internal

type ToDoList struct {
	ID       string `json:"id" bson:"_id,omitempty"`
	Task     string `json:"task" bson:"task"`
	Deadline string `json:"deadline" bson:"deadline"`
	Status   bool   `json:"status" bson:"status"`
}

type CreateListDTO struct {
	Task     string `json:"task"`
	Deadline string `json:"deadline"`
	Status   bool   `json:"status"`
}
