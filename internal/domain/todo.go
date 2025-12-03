package domain

import "gorm.io/gorm"

type Todo struct {
	gorm.Model
	Name        string `gorm:"not null;size=200"`
	Description string `gorm:"not null;size:300"`
	Completed   bool   `gorm:"default:false"`
	UserID      uint   `gorm:"not null"`
}
