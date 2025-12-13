package domain

import (
	"time"

	"gorm.io/gorm"
)

type RefreshToken struct {
	gorm.Model
	UserID    uint      `gorm:"index;not null"`
	TokenHash string    `gorm:"size:64;uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"index;not null"`
	Revoked   *time.Time
}
