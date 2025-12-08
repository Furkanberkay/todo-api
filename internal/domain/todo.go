package domain

import (
	"context"

	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	Name        string `gorm:"not null;size=200"`
	Description string `gorm:"not null;size:300"`
	Completed   bool   `gorm:"default:false"`
	UserID      uint   `gorm:"not null"`
}

type TodoRepository interface {
	GetTodos(ctx context.Context, page int, limit int, userID uint) ([]Todo, int, error)
}
