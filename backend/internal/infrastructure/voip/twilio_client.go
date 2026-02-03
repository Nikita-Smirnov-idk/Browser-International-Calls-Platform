package voip

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/domain"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type TwilioClient struct {
	client         *twilio.RestClient
	fromNumber     string
	sessionManager *SessionManager
}

func NewTwilioClient(cfg *Config) (*TwilioClient, error) {
	if cfg.AccountSID == "" || cfg.AuthToken == "" {
		return nil, fmt.Errorf("twilio credentials are required")
	}

	if cfg.FromNumber == "" {
		return nil, fmt.Errorf("from_number is required")
	}

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: cfg.AccountSID,
		Password: cfg.AuthToken,
	})

	sessionManager := NewSessionManager()

	return &TwilioClient{
		client:         client,
		fromNumber:     cfg.FromNumber,
		sessionManager: sessionManager,
	}, nil
}

func (c *TwilioClient) InitiateCall(ctx context.Context, phoneNumber string) (*domain.CallSession, error) {
	if phoneNumber == "" {
		return nil, domain.ErrInvalidPhoneNumber
	}

	sessionID := generateSessionID()

	params := &openapi.CreateCallParams{}
	params.SetTo(phoneNumber)
	params.SetFrom(c.fromNumber)
	params.SetUrl("http://demo.twilio.com/docs/voice.xml")

	resp, err := c.client.Api.CreateCall(params)
	if err != nil {
		slog.Error("failed to create twilio call", "error", err, "phone", phoneNumber)
		if isTwilioInvalidNumberError(err) {
			return nil, domain.ErrInvalidPhoneNumber
		}
		return nil, ErrVoIPServiceUnavailable
	}

	session := &domain.CallSession{
		SessionID:   sessionID,
		PhoneNumber: phoneNumber,
		SDPOffer:    generateMockSDP(),
		Status:      domain.SessionStatusInitialized,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(5 * time.Minute),
	}

	c.sessionManager.AddSession(session)

	slog.Info("twilio call initiated", 
		"session_id", sessionID, 
		"twilio_call_sid", *resp.Sid,
		"phone", phoneNumber)

	return session, nil
}

func (c *TwilioClient) TerminateCall(ctx context.Context, sessionID string) error {
	session := c.sessionManager.GetSession(sessionID)
	if session == nil {
		return ErrSessionNotFound
	}

	session.Status = domain.SessionStatusCompleted
	c.sessionManager.RemoveSession(sessionID)

	slog.Info("call terminated", "session_id", sessionID)

	return nil
}

func (c *TwilioClient) GetSessionStatus(ctx context.Context, sessionID string) (domain.SessionStatus, error) {
	session := c.sessionManager.GetSession(sessionID)
	if session == nil {
		return "", ErrSessionNotFound
	}

	return session.Status, nil
}

func (c *TwilioClient) Close() error {
	c.sessionManager.Close()
	return nil
}

func generateSessionID() string {
	return fmt.Sprintf("sess_%d", time.Now().UnixNano())
}

func isTwilioInvalidNumberError(err error) bool {
	s := err.Error()
	return strings.Contains(s, "21211") ||
		strings.Contains(s, "21614") ||
		strings.Contains(strings.ToLower(s), "not a valid") ||
		strings.Contains(strings.ToLower(s), "invalid phone")
}

func generateMockSDP() string {
	return `v=0
o=- 0 0 IN IP4 127.0.0.1
s=WebRTC Call
t=0 0
m=audio 9 UDP/TLS/RTP/SAVPF 0 8 101
c=IN IP4 0.0.0.0
a=rtcp:9 IN IP4 0.0.0.0
a=sendrecv`
}

