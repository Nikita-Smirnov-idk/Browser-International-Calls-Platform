package history

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/domain"
)

type CallHistoryItem struct {
	CallID      string    `json:"callId"`
	PhoneNumber string    `json:"phoneNumber"`
	StartTime   time.Time `json:"startTime"`
	Duration    int       `json:"duration"`
	Status      string    `json:"status"`
}

type ListHistoryInput struct {
	UserID   string
	Page     int
	Limit    int
	DateFrom *time.Time
	DateTo   *time.Time
}

type ListHistoryOutput struct {
	Calls []*CallHistoryItem `json:"calls"`
	Total int                `json:"total"`
	Page  int                `json:"page"`
	Limit int                `json:"limit"`
}

type ListHistoryUseCase struct {
	callRepo domain.CallRepository
}

func NewListHistoryUseCase(callRepo domain.CallRepository) *ListHistoryUseCase {
	return &ListHistoryUseCase{callRepo: callRepo}
}

func (uc *ListHistoryUseCase) Execute(ctx context.Context, input ListHistoryInput) (*ListHistoryOutput, error) {
	if input.UserID == "" {
		return nil, errors.New("user_id is required")
	}

	calls, err := uc.callRepo.ListByUserID(ctx, input.UserID)
	if err != nil {
		slog.Error("failed to get calls history", "error", err, "user_id", input.UserID)
		return nil, errors.New("failed to get calls history")
	}

	filteredCalls := calls
	if input.DateFrom != nil || input.DateTo != nil {
		filteredCalls = make([]*domain.Call, 0)
		for _, call := range calls {
			if input.DateFrom != nil && call.StartTime.Before(*input.DateFrom) {
				continue
			}
			if input.DateTo != nil && call.StartTime.After(*input.DateTo) {
				continue
			}
			filteredCalls = append(filteredCalls, call)
		}
	}

	total := len(filteredCalls)

	page := input.Page
	if page < 1 {
		page = 1
	}

	limit := input.Limit
	if limit < 1 {
		limit = 20
	}

	start := (page - 1) * limit
	end := start + limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedCalls := filteredCalls[start:end]

	items := make([]*CallHistoryItem, 0, len(paginatedCalls))
	for _, call := range paginatedCalls {
		items = append(items, &CallHistoryItem{
			CallID:      call.ID,
			PhoneNumber: call.PhoneNumber,
			StartTime:   call.StartTime,
			Duration:    call.Duration,
			Status:      string(call.Status),
		})
	}

	return &ListHistoryOutput{
		Calls: items,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}
