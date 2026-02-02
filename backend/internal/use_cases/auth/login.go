package auth

import (
	"context"

	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/domain"
)

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	AccessToken string
	UserID      string
}

type LoginUseCase struct {
	userRepo domain.UserRepository
}

func NewLoginUseCase(userRepo domain.UserRepository) *LoginUseCase {
	return &LoginUseCase{userRepo: userRepo}
}

func (uc *LoginUseCase) Execute(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	_ = input
	_ = uc.userRepo
	return nil, nil
}
