package auth

import (
	"context"
	"errors"
	"log/slog"
	"regexp"
	"time"

	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Email    string
	Password string
}

type RegisterOutput struct {
	UserID string
	Email  string
}

type RegisterUseCase struct {
	userRepo domain.UserRepository
}

func NewRegisterUseCase(userRepo domain.UserRepository) *RegisterUseCase {
	return &RegisterUseCase{userRepo: userRepo}
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func (uc *RegisterUseCase) Execute(ctx context.Context, input RegisterInput) (*RegisterOutput, error) {
	if !emailRegex.MatchString(input.Email) {
		return nil, errors.New("invalid email format")
	}

	if len(input.Password) < 6 {
		return nil, errors.New("password must be at least 6 characters")
	}

	existing, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		slog.Error("failed to check existing user", "error", err)
		return nil, errors.New("failed to check existing user")
	}

	if existing != nil {
		return nil, errors.New("user with this email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("failed to hash password", "error", err)
		return nil, errors.New("failed to hash password")
	}

	user := &domain.User{
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		slog.Error("failed to create user", "error", err, "email", input.Email)
		return nil, errors.New("failed to create user")
	}

	slog.Info("user registered successfully", "user_id", user.ID, "email", user.Email)

	return &RegisterOutput{
		UserID: user.ID,
		Email:  user.Email,
	}, nil
}
