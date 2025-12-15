package todo

import (
	"time"
	"todoApp3/internal/domain"
)

type CreateTodoRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=30"`
	Description string `json:"description" validate:"required,min=3,max=100"`
}

type PatchTodoRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=3"`
	Description *string `json:"description" validate:"omitempty,min=3"`
	Completed   *bool   `json:"completed" validate:"omitempty"`
}

type UpdateRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=30"`
	Description string `json:"description" validate:"required,min=3,max=100"`
	Completed   *bool  `json:"completed" validate:"required"`
}

type PaginationMeta struct {
	CurrentPage int   `json:"currentPage"`
	TotalPages  int   `json:"totalPages"`
	TotalItems  int64 `json:"totalItems"`
	Limit       int   `json:"limit"`
}

type PaginatedResponse struct {
	Data interface{}    `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

type TodoResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
}

func toTodoResponseList(todos []domain.Todo) []TodoResponse {
	response := make([]TodoResponse, len(todos))

	for i, t := range todos {
		response[i] = TodoResponse{
			ID:          t.ID,
			Name:        t.Name,
			Description: t.Description,
			Completed:   t.Completed,
			CreatedAt:   t.CreatedAt,
		}
	}

	return response
}
