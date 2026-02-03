package voip

import (
	"fmt"
	"time"

	"github.com/twilio/twilio-go/client/jwt"
)

const defaultVoiceTokenTTL = 3600

type TokenConfig struct {
	AccountSid   string
	APIKeySid    string
	APIKeySecret string
	TwimlAppSid  string
}

type TokenGenerator struct {
	cfg *TokenConfig
}

func NewTokenGenerator(cfg *TokenConfig) (*TokenGenerator, error) {
	if cfg == nil || cfg.AccountSid == "" || cfg.APIKeySid == "" || cfg.APIKeySecret == "" || cfg.TwimlAppSid == "" {
		return nil, fmt.Errorf("voice token config incomplete")
	}
	return &TokenGenerator{cfg: cfg}, nil
}

func (g *TokenGenerator) GetToken(identity string, ttlSec int) (string, error) {
	if ttlSec <= 0 {
		ttlSec = defaultVoiceTokenTTL
	}
	params := jwt.AccessTokenParams{
		AccountSid:    g.cfg.AccountSid,
		SigningKeySid: g.cfg.APIKeySid,
		Secret:        g.cfg.APIKeySecret,
		Identity:      identity,
		Ttl:           float64(ttlSec),
		Nbf:           float64(time.Now().Unix()),
	}
	token := jwt.CreateAccessToken(params)
	voiceGrant := &jwt.VoiceGrant{
		Outgoing: jwt.Outgoing{
			ApplicationSid: g.cfg.TwimlAppSid,
		},
	}
	token.AddGrant(voiceGrant)
	return token.ToJwt()
}
