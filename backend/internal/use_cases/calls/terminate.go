package calls

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/domain"
)

type TerminateCallInput struct {
	UserID string
	CallID string
}

type TerminateCallOutput struct {
	CallID   string
	Duration int
	Status   string
}

type TerminateCallUseCase struct {
	callRepo    domain.CallRepository
	voipService domain.VoIPService
}

func NewTerminateCallUseCase(callRepo domain.CallRepository, voipService domain.VoIPService) *TerminateCallUseCase {
	return &TerminateCallUseCase{
		callRepo:    callRepo,
		voipService: voipService,
	}
}

func (uc *TerminateCallUseCase) Execute(ctx context.Context, input TerminateCallInput) (*TerminateCallOutput, error) {
	if input.CallID == "" {
		return nil, errors.New("call_id is required")
	}

	if input.UserID == "" {
		return nil, errors.New("user_id is required")
	}

	call, err := uc.callRepo.GetByID(ctx, input.CallID)
	if err != nil {
		slog.Error("failed to get call", "error", err, "call_id", input.CallID)
		return nil, errors.New("failed to get call")
	}

	if call == nil {
		return nil, errors.New("call not found")
	}

	if call.UserID != input.UserID {
		slog.Warn("unauthorized call termination attempt", 
			"call_id", input.CallID, 
			"user_id", input.UserID,
			"call_user_id", call.UserID)
		return nil, errors.New("unauthorized")
	}

	if call.SessionID != "" {
		if err := uc.voipService.TerminateCall(ctx, call.SessionID); err != nil {
			slog.Warn("failed to terminate voip session", 
				"error", err, 
				"session_id", call.SessionID)
		}
	}

	duration := int(time.Since(call.StartTime).Seconds())
	call.Duration = duration
	call.Status = domain.CallStatusCompleted

	if err := uc.callRepo.Update(ctx, call); err != nil {
		slog.Error("failed to update call", 
			"error", err, 
			"call_id", input.CallID)
		return nil, errors.New("failed to update call")
	}

	slog.Info("call terminated successfully", 
		"call_id", call.ID, 
		"user_id", input.UserID, 
		"session_id", call.SessionID,
		"duration", duration)

	return &TerminateCallOutput{
		CallID:   call.ID,
		Duration: duration,
		Status:   string(call.Status),
	}, nil
}

