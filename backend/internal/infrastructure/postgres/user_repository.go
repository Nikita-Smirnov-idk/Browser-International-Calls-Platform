package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/domain"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

type userModel struct {
	ID           string `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()"`
	Email        string `gorm:"column:email;uniqueIndex;not null"`
	PasswordHash string `gorm:"column:password_hash;not null"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (userModel) TableName() string {
	return "users"
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	model := &userModel{
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	user.ID = model.ID
	return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var model userModel
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &domain.User{
		ID:           model.ID,
		Email:        model.Email,
		PasswordHash: model.PasswordHash,
	}, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var model userModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &domain.User{
		ID:           model.ID,
		Email:        model.Email,
		PasswordHash: model.PasswordHash,
	}, nil
}
