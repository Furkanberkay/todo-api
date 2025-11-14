package service

import (
	"context"
	"strings"
	"todo-api/domain"
	"todo-api/repository"
)

type CreateTodoRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type TodoService struct {
	repo repository.SqliteTodoRepository
}

func NewTodoService(repository repository.SqliteTodoRepository) *TodoService {
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
