package domain

import "time"

type CallStatus string

const (
	CallStatusInitiated CallStatus = "initiated"
	CallStatusCompleted CallStatus = "completed"
	CallStatusFailed    CallStatus = "failed"
	CallStatusCanceled  CallStatus = "canceled"
)

type Call struct {
	ID          string
	UserID      string
	PhoneNumber string
	StartTime   time.Time
	Duration    int
	Status      CallStatus
	CreatedAt   time.Time
}
