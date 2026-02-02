package calls

import (
	"context"

	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/domain"
)

type StartCallInput struct {
	UserID      string
	CountryCode string
	PhoneNumber string
}

type StartCallOutput struct {
	CallID    string
	RTCConfig map[string]interface{}
}

type StartCallUseCase struct {
	callRepo domain.CallRepository
}

func NewStartCallUseCase(callRepo domain.CallRepository) *StartCallUseCase {
	return &StartCallUseCase{callRepo: callRepo}
}

func (uc *StartCallUseCase) Execute(ctx context.Context, input StartCallInput) (*StartCallOutput, error) {
	_ = input
	_ = uc.callRepo
	_ = domain.Call{}
	return nil, nil
}
