package domain

import "time"

type CallStatus string

const (
	CallStatusCompleted CallStatus = "completed"
	CallStatusFailed    CallStatus = "failed"
	CallStatusCanceled  CallStatus = "canceled"
)

type Call struct {
	ID             string
	UserID         string
	CountryCode    string
	PhoneNumber    string
	StartedAt      time.Time
	DurationSeconds int
	Status         CallStatus
}
