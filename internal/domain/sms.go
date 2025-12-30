package domain

import "time"

type SmsJob struct {
	Email      string
	Name       string
	Surname    string
	Message    string
	EnqueuedAt time.Time
}
