package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/domain"
	"gorm.io/gorm"
)

type CallRepository struct {
	db *gorm.DB
}

func NewCallRepository(db *gorm.DB) *CallRepository {
	return &CallRepository{db: db}
}

type callModel struct {
	ID          string    `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID      string    `gorm:"column:user_id;not null;index"`
	PhoneNumber string    `gorm:"column:phone_number;not null"`
	StartTime   time.Time `gorm:"column:start_time;not null;index"`
	Duration    int       `gorm:"column:duration;default:0"`
	Status      string    `gorm:"column:status;not null;default:initiated"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (callModel) TableName() string {
	return "calls"
}

func (r *CallRepository) Create(ctx context.Context, call *domain.Call) error {
	model := &callModel{
		UserID:      call.UserID,
		PhoneNumber: call.PhoneNumber,
		StartTime:   call.StartTime,
		Duration:    call.Duration,
		Status:      string(call.Status),
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	call.ID = model.ID
	call.CreatedAt = model.CreatedAt
	return nil
}

func (r *CallRepository) Update(ctx context.Context, call *domain.Call) error {
	updates := map[string]interface{}{
		"duration": call.Duration,
		"status":   string(call.Status),
	}

	result := r.db.WithContext(ctx).Model(&callModel{}).Where("id = ?", call.ID).Updates(updates)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("call not found")
	}

	return nil
}

func (r *CallRepository) GetByID(ctx context.Context, id string) (*domain.Call, error) {
	var model callModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &domain.Call{
		ID:          model.ID,
		UserID:      model.UserID,
		PhoneNumber: model.PhoneNumber,
		StartTime:   model.StartTime,
		Duration:    model.Duration,
		Status:      domain.CallStatus(model.Status),
		CreatedAt:   model.CreatedAt,
	}, nil
}

func (r *CallRepository) ListByUserID(ctx context.Context, userID string) ([]*domain.Call, error) {
	var models []callModel
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("start_time DESC").
		Find(&models).Error

	if err != nil {
		return nil, err
	}

	calls := make([]*domain.Call, 0, len(models))
	for _, model := range models {
		calls = append(calls, &domain.Call{
			ID:          model.ID,
			UserID:      model.UserID,
			PhoneNumber: model.PhoneNumber,
			StartTime:   model.StartTime,
			Duration:    model.Duration,
			Status:      domain.CallStatus(model.Status),
			CreatedAt:   model.CreatedAt,
		})
	}

	return calls, nil
}
