package auth

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	AccessToken string
	UserID      string
	Email       string
}

type LoginUseCase struct {
	userRepo   domain.UserRepository
	jwtService JWTService
}

type JWTService interface {
	GenerateToken(userID, email string) (string, error)
}

func NewLoginUseCase(userRepo domain.UserRepository, jwtService JWTService) *LoginUseCase {
	return &LoginUseCase{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	user, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		slog.Error("failed to get user by email", "error", err)
		return nil, errors.New("failed to get user")
	}

	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := uc.jwtService.GenerateToken(user.ID, user.Email)
	if err != nil {
		slog.Error("failed to generate token", "error", err, "user_id", user.ID)
		return nil, errors.New("failed to generate token")
	}

	slog.Info("user logged in successfully", "user_id", user.ID, "email", user.Email)

	return &LoginOutput{
		AccessToken: token,
		UserID:      user.ID,
		Email:       user.Email,
	}, nil
}
