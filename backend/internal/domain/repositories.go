package domain

import "context"

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
}

type CallRepository interface {
	Create(ctx context.Context, call *Call) error
	Update(ctx context.Context, call *Call) error
	GetByID(ctx context.Context, id string) (*Call, error)
	ListByUserID(ctx context.Context, userID string) ([]*Call, error)
}
