package handlers

type CreateTodoRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateTodoRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Completed   *bool  `json:"completed"`
}
