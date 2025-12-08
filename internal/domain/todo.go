package domain

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/plugin/optimisticlock"
)

type Todo struct {
	gorm.Model
	Name        string `gorm:"not null;size=200"`
	Description string `gorm:"not null;size:300"`
	Completed   bool   `gorm:"default:false"`
	UserID      uint   `gorm:"not null"`
	Version     optimisticlock.Version
}

type TodoRepository interface {
	GetTodos(ctx context.Context, page int, limit int, userID uint) ([]Todo, int, error)
	GetTodoByID(ctx context.Context, userID uint, todoID uint) (*Todo, error)
	CreateTodo(ctx context.Context, todo *Todo) error
	UpdateTodo(ctx context.Context, todo *Todo) error
}

var ErrTodoNotFound = errors.New("todo not found")
var ErrConflict = errors.New("data has been modified by another user")
