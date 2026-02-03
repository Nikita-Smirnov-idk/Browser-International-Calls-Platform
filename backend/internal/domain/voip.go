package domain

import (
	"context"
	"time"
)

type VoIPService interface {
	InitiateCall(ctx context.Context, phoneNumber string) (*CallSession, error)
	TerminateCall(ctx context.Context, sessionID string) error
	GetSessionStatus(ctx context.Context, sessionID string) (SessionStatus, error)
}

type CallSession struct {
	SessionID    string
	PhoneNumber  string
	SDPOffer     string
	Status       SessionStatus
	CreatedAt    time.Time
	ExpiresAt    time.Time
}

type SessionStatus string

const (
	SessionStatusInitialized SessionStatus = "initialized"
	SessionStatusConnecting  SessionStatus = "connecting"
	SessionStatusActive      SessionStatus = "active"
	SessionStatusCompleted   SessionStatus = "completed"
	SessionStatusFailed      SessionStatus = "failed"
)

type WebRTCConfig struct {
	IceServers []IceServer `json:"iceServers"`
}

type IceServer struct {
	URLs       []string `json:"urls"`
	Username   string   `json:"username,omitempty"`
	Credential string   `json:"credential,omitempty"`
}

