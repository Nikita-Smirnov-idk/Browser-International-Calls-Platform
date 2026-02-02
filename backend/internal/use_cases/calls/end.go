package calls

import (
	"context"

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
	_ = input
	_ = uc.callRepo
	return nil
}
