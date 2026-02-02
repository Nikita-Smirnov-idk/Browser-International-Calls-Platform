package auth

import "context"

type LogoutInput struct {
	UserID string
	Token  string
}

type LogoutUseCase struct{}

func NewLogoutUseCase() *LogoutUseCase {
	return &LogoutUseCase{}
}

func (uc *LogoutUseCase) Execute(ctx context.Context, input LogoutInput) error {
	_ = input
	return nil
}
