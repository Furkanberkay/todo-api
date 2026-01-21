package domain

import "time"

type EmailJob struct {
	Email      string
	Name       string
	Surname    string
	Message    string
	EnqueuedAt time.Time
}
