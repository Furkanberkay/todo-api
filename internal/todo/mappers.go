package todo

func UpdateRequestToUpdateInput(todoID uint, updateDTO *UpdateRequest) *UpdateTodoInput {
	return &UpdateTodoInput{
		ID:          todoID,
		Name:        updateDTO.Name,
		Description: updateDTO.Description,
		Completed:   *updateDTO.Completed,
	}
}
