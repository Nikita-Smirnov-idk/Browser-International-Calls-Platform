package calls

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/domain"
)

type StartCallInput struct {
	UserID      string
	PhoneNumber string
}

type StartCallOutput struct {
	CallID    string
	StartTime time.Time
}

type StartCallUseCase struct {
	callRepo domain.CallRepository
}

func NewStartCallUseCase(callRepo domain.CallRepository) *StartCallUseCase {
	return &StartCallUseCase{callRepo: callRepo}
}

func (uc *StartCallUseCase) Execute(ctx context.Context, input StartCallInput) (*StartCallOutput, error) {
	if input.UserID == "" {
		return nil, errors.New("user_id is required")
	}

	if input.PhoneNumber == "" {
		return nil, errors.New("phone_number is required")
	}

	call := &domain.Call{
		UserID:      input.UserID,
		PhoneNumber: input.PhoneNumber,
		StartTime:   time.Now(),
		Duration:    0,
		Status:      domain.CallStatusInitiated,
	}

	if err := uc.callRepo.Create(ctx, call); err != nil {
		slog.Error("failed to create call", "error", err, "user_id", input.UserID)
		return nil, errors.New("failed to create call")
	}

	slog.Info("call created", "call_id", call.ID, "user_id", input.UserID, "phone", input.PhoneNumber)

	return &StartCallOutput{
		CallID:    call.ID,
		StartTime: call.StartTime,
	}, nil
}
