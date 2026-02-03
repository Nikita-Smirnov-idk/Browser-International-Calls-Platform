package calls

import (
	"context"
	"errors"
	"log/slog"
	"regexp"
	"time"

	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/domain"
)

var e164Re = regexp.MustCompile(`^\+[1-9]\d{6,14}$`)

type InitiateCallInput struct {
	UserID      string
	PhoneNumber string
}

type InitiateCallOutput struct {
	CallID     string
	SessionID  string
	SDPOffer   string
	Status     string
	StartTime  time.Time
	VoiceToken string
}

type VoiceTokenGenerator interface {
	GetToken(identity string, ttlSec int) (string, error)
}

type InitiateCallUseCase struct {
	callRepo       domain.CallRepository
	voipService    domain.VoIPService
	tokenGenerator VoiceTokenGenerator
}

func NewInitiateCallUseCase(callRepo domain.CallRepository, voipService domain.VoIPService, tokenGenerator VoiceTokenGenerator) *InitiateCallUseCase {
	return &InitiateCallUseCase{
		callRepo:       callRepo,
		voipService:    voipService,
		tokenGenerator: tokenGenerator,
	}
}

func (uc *InitiateCallUseCase) Execute(ctx context.Context, input InitiateCallInput) (*InitiateCallOutput, error) {
	if input.UserID == "" {
		return nil, errors.New("user_id is required")
	}

	if input.PhoneNumber == "" {
		return nil, errors.New("phone_number is required")
	}

	if !e164Re.MatchString(input.PhoneNumber) {
		return nil, domain.ErrInvalidPhoneNumber
	}

	if uc.tokenGenerator != nil {
		call := &domain.Call{
			UserID:      input.UserID,
			PhoneNumber: input.PhoneNumber,
			StartTime:   time.Now(),
			Duration:    0,
			Status:      domain.CallStatusConnecting,
			SessionID:   "voice_sdk",
			SDPOffer:    "",
		}
		if err := uc.callRepo.Create(ctx, call); err != nil {
			slog.Error("failed to create call record", "error", err, "user_id", input.UserID)
			return nil, errors.New("failed to create call record")
		}
		token, err := uc.tokenGenerator.GetToken(input.UserID, 3600)
		if err != nil {
			slog.Error("failed to generate voice token", "error", err, "user_id", input.UserID)
			return nil, errors.New("failed to generate voice token")
		}
		slog.Info("call initiated with voice sdk", "call_id", call.ID, "user_id", input.UserID, "phone", input.PhoneNumber)
		return &InitiateCallOutput{
			CallID:     call.ID,
			SessionID:  call.SessionID,
			SDPOffer:   "",
			Status:     string(call.Status),
			StartTime:  call.StartTime,
			VoiceToken: token,
		}, nil
	}

	session, err := uc.voipService.InitiateCall(ctx, input.PhoneNumber)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidPhoneNumber) {
			return nil, err
		}
		slog.Error("failed to initiate voip call",
			"error", err,
			"user_id", input.UserID,
			"phone", input.PhoneNumber)
		return nil, errors.New("failed to initiate call")
	}

	call := &domain.Call{
		UserID:      input.UserID,
		PhoneNumber: input.PhoneNumber,
		StartTime:   time.Now(),
		Duration:    0,
		Status:      domain.CallStatusConnecting,
		SessionID:   session.SessionID,
		SDPOffer:    session.SDPOffer,
	}

	if err := uc.callRepo.Create(ctx, call); err != nil {
		slog.Error("failed to create call record", 
			"error", err, 
			"user_id", input.UserID,
			"session_id", session.SessionID)
		return nil, errors.New("failed to create call record")
	}

	slog.Info("call initiated successfully", 
		"call_id", call.ID, 
		"user_id", input.UserID, 
		"session_id", session.SessionID,
		"phone", input.PhoneNumber)

	return &InitiateCallOutput{
		CallID:    call.ID,
		SessionID: session.SessionID,
		SDPOffer:  session.SDPOffer,
		Status:    string(call.Status),
		StartTime: call.StartTime,
	}, nil
}

