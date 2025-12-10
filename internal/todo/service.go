package todo

import (
	"context"
	"log/slog"
	"todoApp3/internal/domain"
)

type Service struct {
	repository domain.TodoRepository
	logger     *slog.Logger
}

func NewTodoService(repository domain.TodoRepository, logger *slog.Logger) *Service {
	return &Service{
		repository: repository,
		logger:     logger,
	}
}

func (s *Service) GetTodos(ctx context.Context, page int, limit int, userID uint) ([]domain.Todo, int, error) {
	if page <= 0 {
		page = 1
	}

	if limit <= 0 || limit > 30 {
		limit = 30
	}

	todos, totalCount, err := s.repository.GetTodos(ctx, page, limit, userID)

	if err != nil {
		return nil, 0, err
	}

	return todos, totalCount, nil
}

func (s *Service) GetTodoByID(ctx context.Context, userID uint, todoID uint) (*domain.Todo, error) {
	return s.repository.GetTodoByID(ctx, userID, todoID)
}

func (s *Service) CreateTodo(ctx context.Context, input *CreateTodoInput, userID uint) (*domain.Todo, error) {
	domainTodo := domain.Todo{
		Name:        input.Name,
		Description: input.Description,
		Completed:   false,
		UserID:      userID,
	}

	if err := s.repository.CreateTodo(ctx, &domainTodo); err != nil {
		return nil, err
	}

	return &domainTodo, nil

}

func (s *Service) UpdateTodo(ctx context.Context, input *UpdateTodoInput, userID uint) (*domain.Todo, error) {

	todo, err := s.repository.GetTodoByID(ctx, userID, input.ID)

	if err != nil {
		return nil, err
	}

	if input.Version != todo.Version.Int64 {
		return nil, domain.ErrConflict
	}

	todo.Name = input.Name
	todo.Description = input.Description
	todo.Completed = input.Completed

	if err := s.repository.UpdateTodo(ctx, todo); err != nil {
		return nil, err
	}

	todo.Version.Int64++
	return todo, nil
}

func (s *Service) DeleteTodo(ctx context.Context, userID uint, todoID uint) error {

	return s.repository.DeleteTodo(ctx, userID, todoID)
}

func (s *Service) PatchTodo(ctx context.Context, userID uint, input *PatchTodoInput) (*domain.Todo, error) {

	if input.Name == nil && input.Description == nil && input.Completed == nil {
		return nil, domain.InvalidInput
	}

	todo, err := s.GetTodoByID(ctx, userID, input.ID)
	if err != nil {
		return nil, err
	}
	if input.Name != nil {
		todo.Name = *input.Name
	}
	if input.Description != nil {
		todo.Description = *input.Description
	}
	if input.Completed != nil {
		todo.Completed = *input.Completed
	}

	if input.Version != todo.Version.Int64 {
		return nil, domain.ErrConflict
	}

	if err := s.repository.UpdateTodo(ctx, todo); err != nil {
		return nil, err
	}
	todo.Version.Int64++
	return todo, nil
}
