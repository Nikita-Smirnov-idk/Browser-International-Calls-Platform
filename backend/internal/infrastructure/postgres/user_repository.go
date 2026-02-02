package postgres

import (
	"context"

	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/domain"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	_ = ctx
	_ = user
	return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	_ = ctx
	_ = email
	return nil, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	_ = ctx
	_ = id
	return nil, nil
}
