package history

import (
	"context"
	"time"

	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/domain"
)

type CallHistoryItem struct {
	CallID          string    `json:"callId"`
	PhoneNumber     string    `json:"phoneNumber"`
	StartedAt       time.Time `json:"startedAt"`
	DurationSeconds int       `json:"durationSeconds"`
	Status          string    `json:"status"`
}

type ListHistoryInput struct {
	UserID string
}

type ListHistoryUseCase struct {
	callRepo domain.CallRepository
}

func NewListHistoryUseCase(callRepo domain.CallRepository) *ListHistoryUseCase {
	return &ListHistoryUseCase{callRepo: callRepo}
}

func (uc *ListHistoryUseCase) Execute(ctx context.Context, input ListHistoryInput) ([]*CallHistoryItem, error) {
	_ = input
	_ = uc.callRepo
	return nil, nil
}
