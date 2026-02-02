package postgres

import (
	"context"

	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/domain"
)

type CallRepository struct{}

func NewCallRepository() *CallRepository {
	return &CallRepository{}
}

func (r *CallRepository) Create(ctx context.Context, call *domain.Call) error {
	_ = ctx
	_ = call
	return nil
}

func (r *CallRepository) Update(ctx context.Context, call *domain.Call) error {
	_ = ctx
	_ = call
	return nil
}

func (r *CallRepository) GetByID(ctx context.Context, id string) (*domain.Call, error) {
	_ = ctx
	_ = id
	return nil, nil
}

func (r *CallRepository) ListByUserID(ctx context.Context, userID string) ([]*domain.Call, error) {
	_ = ctx
	_ = userID
	return nil, nil
}
