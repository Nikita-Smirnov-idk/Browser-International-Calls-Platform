package auth

import (
	"context"

	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/domain"
)

type RegisterInput struct {
	Email    string
	Password string
}

type RegisterOutput struct {
	UserID string
}

type RegisterUseCase struct {
	userRepo domain.UserRepository
}

func NewRegisterUseCase(userRepo domain.UserRepository) *RegisterUseCase {
	return &RegisterUseCase{userRepo: userRepo}
}

func (uc *RegisterUseCase) Execute(ctx context.Context, input RegisterInput) (*RegisterOutput, error) {
	_ = input
	_ = uc.userRepo
	return nil, nil
}
