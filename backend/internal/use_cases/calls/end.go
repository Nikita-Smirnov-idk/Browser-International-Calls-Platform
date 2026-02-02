package calls

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/domain"
)

type EndCallInput struct {
	UserID string
	CallID string
}

type EndCallUseCase struct {
	callRepo domain.CallRepository
}

func NewEndCallUseCase(callRepo domain.CallRepository) *EndCallUseCase {
	return &EndCallUseCase{callRepo: callRepo}
}

func (uc *EndCallUseCase) Execute(ctx context.Context, input EndCallInput) error {
	if input.CallID == "" {
		return errors.New("call_id is required")
	}

	call, err := uc.callRepo.GetByID(ctx, input.CallID)
	if err != nil {
		slog.Error("failed to get call", "error", err, "call_id", input.CallID)
		return errors.New("failed to get call")
	}

	if call == nil {
		return errors.New("call not found")
	}

	if call.UserID != input.UserID {
		return errors.New("unauthorized")
	}

	duration := int(time.Since(call.StartTime).Seconds())
	call.Duration = duration
	call.Status = domain.CallStatusCompleted

	if err := uc.callRepo.Update(ctx, call); err != nil {
		slog.Error("failed to update call", "error", err, "call_id", input.CallID)
		return errors.New("failed to update call")
	}

	slog.Info("call ended", "call_id", call.ID, "user_id", input.UserID, "duration", duration)

	return nil
}
