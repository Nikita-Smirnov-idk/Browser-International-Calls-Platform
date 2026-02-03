package voip

import (
	"errors"

	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/domain"
)

var (
	ErrVoIPServiceUnavailable = errors.New("voip service unavailable")
	ErrInvalidPhoneNumber     = errors.New("invalid phone number")
	ErrSessionNotFound        = errors.New("session not found")
	ErrCallAlreadyActive      = errors.New("call already active")
	ErrUnauthorized           = errors.New("unauthorized")
)

type Client interface {
	domain.VoIPService
	Close() error
}

type Config struct {
	Provider   string
	AccountSID string
	AuthToken  string
	APIKey     string
	FromNumber string
}

func NewClient(cfg *Config) (Client, error) {
	if cfg == nil {
		return nil, errors.New("config is required")
	}

	switch cfg.Provider {
	case "twilio":
		return NewTwilioClient(cfg)
	case "mock":
		return NewMockClient(cfg)
	default:
		return nil, errors.New("unsupported voip provider: " + cfg.Provider)
	}
}

