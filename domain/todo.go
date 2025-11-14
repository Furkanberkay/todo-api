package domain

import (
	"context"
	"errors"
)

type Todo struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type TodoRepository interface {
	GetTodos(ctx context.Context) ([]Todo, error)
	UpdateTodo(ctx context.Context, todo *Todo) (*Todo, error)
	CreateTodo(ctx context.Context, todo *Todo) error
	DeleteTodo(ctx context.Context, id int) error
	GetTodoByID(ctx context.Context, id int) (*Todo, error)
}

var (
	ErrInternal     = errors.New("Server Internel Error")
	ErrTodoNotFound = errors.New("Todo not found")
	ErrValidation   = errors.New("validation error")
)
