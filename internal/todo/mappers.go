package todo

func UpdateRequestToUpdateInput(todoID int, updateDTO *UpdateRequest) *UpdateTodoInput {
	return &UpdateTodoInput{
		ID:          uint(todoID),
		Name:        updateDTO.Name,
		Description: updateDTO.Description,
		Completed:   *updateDTO.Completed,
		Version:     updateDTO.Version,
	}
}
