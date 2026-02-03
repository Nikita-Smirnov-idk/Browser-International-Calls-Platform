package voip

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/domain"
)

type MockClient struct {
	sessionManager *SessionManager
}

func NewMockClient(cfg *Config) (*MockClient, error) {
	sessionManager := NewSessionManager()

	slog.Info("mock voip client initialized")

	return &MockClient{
		sessionManager: sessionManager,
	}, nil
}

func (c *MockClient) InitiateCall(ctx context.Context, phoneNumber string) (*domain.CallSession, error) {
	if phoneNumber == "" {
		return nil, domain.ErrInvalidPhoneNumber
	}

	sessionID := fmt.Sprintf("mock_sess_%d", time.Now().UnixNano())

	session := &domain.CallSession{
		SessionID:   sessionID,
		PhoneNumber: phoneNumber,
		SDPOffer:    generateMockSDP(),
		Status:      domain.SessionStatusInitialized,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(5 * time.Minute),
	}

	c.sessionManager.AddSession(session)

	slog.Info("mock call initiated", "session_id", sessionID, "phone", phoneNumber)

	return session, nil
}

func (c *MockClient) TerminateCall(ctx context.Context, sessionID string) error {
	session := c.sessionManager.GetSession(sessionID)
	if session == nil {
		return ErrSessionNotFound
	}

	session.Status = domain.SessionStatusCompleted
	c.sessionManager.RemoveSession(sessionID)

	slog.Info("mock call terminated", "session_id", sessionID)

	return nil
}

func (c *MockClient) GetSessionStatus(ctx context.Context, sessionID string) (domain.SessionStatus, error) {
	session := c.sessionManager.GetSession(sessionID)
	if session == nil {
		return "", ErrSessionNotFound
	}

	return session.Status, nil
}

func (c *MockClient) Close() error {
	c.sessionManager.Close()
	return nil
}

