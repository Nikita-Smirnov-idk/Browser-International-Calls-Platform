package auth

import (
	"context"
	"log/slog"
)

type LogoutInput struct {
	UserID string
	Token  string
}

type LogoutUseCase struct{}

func NewLogoutUseCase() *LogoutUseCase {
	return &LogoutUseCase{}
}

func (uc *LogoutUseCase) Execute(ctx context.Context, input LogoutInput) error {
	slog.Info("user logged out", "user_id", input.UserID)
	return nil
}
