package domain

import "time"

type CallStatus string

const (
	CallStatusInitiated CallStatus = "initiated"
	CallStatusConnecting CallStatus = "connecting"
	CallStatusActive    CallStatus = "active"
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
	SessionID   string
	SDPOffer    string
	SDPAnswer   string
}
