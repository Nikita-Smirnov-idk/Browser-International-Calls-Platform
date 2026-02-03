package calls

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/domain"
)

type InitiateCallInput struct {
	UserID      string
	PhoneNumber string
}

type InitiateCallOutput struct {
	CallID    string
	SessionID string
	SDPOffer  string
	Status    string
	StartTime time.Time
}

type InitiateCallUseCase struct {
	callRepo    domain.CallRepository
	voipService domain.VoIPService
}

func NewInitiateCallUseCase(callRepo domain.CallRepository, voipService domain.VoIPService) *InitiateCallUseCase {
	return &InitiateCallUseCase{
		callRepo:    callRepo,
		voipService: voipService,
	}
}

func (uc *InitiateCallUseCase) Execute(ctx context.Context, input InitiateCallInput) (*InitiateCallOutput, error) {
	if input.UserID == "" {
		return nil, errors.New("user_id is required")
	}

	if input.PhoneNumber == "" {
		return nil, errors.New("phone_number is required")
	}

	session, err := uc.voipService.InitiateCall(ctx, input.PhoneNumber)
	if err != nil {
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

