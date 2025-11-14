package service

import (
	"context"
	"errors"
	"strings"
	"todo-api/domain"
)

type CreateTodoRequest struct {
	Name        *string
	Description *string
}
type UpdateTodoRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Completed   *bool   `json:"completed"`
}
type PatchTodoRequest struct {
	Name        *string
	Description *string
	Completed   *bool
}

type TodoService struct {
	repo domain.TodoRepository
}

func NewTodoService(repository domain.TodoRepository) *TodoService {
	return &TodoService{repo: repository}
}

func (s *TodoService) GetTodos(ctx context.Context) ([]domain.Todo, error) {
	return s.repo.GetTodos(ctx)
}
func (s *TodoService) GetTodoByID(ctx context.Context, id int) (*domain.Todo, error) {
	return s.repo.GetTodoByID(ctx, id)
}
func (s *TodoService) CreateTodo(ctx context.Context, createDto CreateTodoRequest) (*domain.Todo, error) {

	if createDto.Name == nil || createDto.Description == nil {
		return nil, domain.ErrValidation
	}
	if strings.TrimSpace(*createDto.Name) == "" || strings.TrimSpace(*createDto.Description) == "" {
		return nil, domain.ErrValidation
	}
	todo := domain.Todo{
		Name:        *createDto.Name,
		Description: *createDto.Description,
		Completed:   false,
	}
	err := s.repo.CreateTodo(ctx, &todo)
	if err != nil {
		return nil, err
	}
	return &todo, nil

}
func (s *TodoService) UpdateTodo(ctx context.Context, updateTodoDto *UpdateTodoRequest, id int) (*domain.Todo, error) {

	if updateTodoDto.Name == nil {
		return nil, domain.ErrNameValidation
	}

	if updateTodoDto.Description == nil {
		return nil, domain.ErrDescriptionValidation
	}

	if updateTodoDto.Completed == nil {
		return nil, domain.ErrCompletedValidation
	}

	updateTodo := domain.Todo{
		Id:          id,
		Name:        *updateTodoDto.Name,
		Description: *updateTodoDto.Description,
		Completed:   *updateTodoDto.Completed,
	}

	updateTodoResult, err := s.repo.UpdateTodo(ctx, &updateTodo)
	if err != nil {
		if errors.Is(err, domain.ErrTodoNotFound) {
			return nil, domain.ErrTodoNotFound
		}
		return nil, domain.ErrInternal
	}
	return updateTodoResult, nil

}
func (s *TodoService) DeleteTodo(ctx context.Context, id int) error {

	return s.repo.DeleteTodo(ctx, id)
}
func (s *TodoService) PatchTodo(ctx context.Context, patchTodo *PatchTodoRequest, id int) (*domain.Todo, error) {
	todo, err := s.repo.GetTodoByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrTodoNotFound) {
			return nil, domain.ErrTodoNotFound
		}
		return nil, domain.ErrInternal
	}
	if patchTodo.Name == nil && patchTodo.Description == nil && patchTodo.Completed == nil {
		return nil, domain.ErrValidation
	}
	if patchTodo.Name != nil {
		todo.Name = *patchTodo.Name
	}
	if patchTodo.Description != nil {
		todo.Description = *patchTodo.Description
	}
	if patchTodo.Completed != nil {
		todo.Completed = *patchTodo.Completed
	}
	updatedTodo, errUpdate := s.repo.UpdateTodo(ctx, todo)
	if errUpdate != nil {
		if errors.Is(errUpdate, domain.ErrTodoNotFound) {
			return nil, domain.ErrTodoNotFound
		}
		return nil, domain.ErrInternal
	}
	return updatedTodo, nil
}
