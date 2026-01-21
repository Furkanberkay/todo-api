package todo

type CreateTodoInput struct {
	Name        string
	Description string
}

type UpdateTodoInput struct {
	ID          uint
	Name        string
	Description string
	Completed   bool
}

type PatchTodoInput struct {
	ID          uint
	Name        *string
	Description *string
	Completed   *bool
}
