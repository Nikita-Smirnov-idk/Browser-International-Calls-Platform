package calls

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/domain"
)

type mockCallRepository struct {
	createError error
	createdCall *domain.Call
}

func (m *mockCallRepository) Create(ctx context.Context, call *domain.Call) error {
	if m.createError != nil {
		return m.createError
	}
	call.ID = "test-call-id"
	call.CreatedAt = time.Now()
	m.createdCall = call
	return nil
}

func (m *mockCallRepository) Update(ctx context.Context, call *domain.Call) error {
	return nil
}

func (m *mockCallRepository) GetByID(ctx context.Context, id string) (*domain.Call, error) {
	return nil, nil
}

func (m *mockCallRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*domain.Call, error) {
	return nil, nil
}

type mockVoIPService struct {
	initiateError error
	session       *domain.CallSession
}

func (m *mockVoIPService) InitiateCall(ctx context.Context, phoneNumber string) (*domain.CallSession, error) {
	if m.initiateError != nil {
		return nil, m.initiateError
	}
	return m.session, nil
}

func (m *mockVoIPService) TerminateCall(ctx context.Context, sessionID string) error {
	return nil
}

func (m *mockVoIPService) GetSessionStatus(ctx context.Context, sessionID string) (domain.SessionStatus, error) {
	return domain.SessionStatusActive, nil
}

func TestInitiateCallUseCase_Execute_Success(t *testing.T) {
	mockRepo := &mockCallRepository{}
	mockVoIP := &mockVoIPService{
		session: &domain.CallSession{
			SessionID:   "test-session-id",
			PhoneNumber: "+491512345678",
			SDPOffer:    "test-sdp-offer",
			Status:      domain.SessionStatusInitialized,
			CreatedAt:   time.Now(),
			ExpiresAt:   time.Now().Add(5 * time.Minute),
		},
	}

	uc := NewInitiateCallUseCase(mockRepo, mockVoIP)

	input := InitiateCallInput{
		UserID:      "test-user-id",
		PhoneNumber: "+491512345678",
	}

	output, err := uc.Execute(context.Background(), input)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output == nil {
		t.Fatal("expected output, got nil")
	}

	if output.CallID != "test-call-id" {
		t.Errorf("expected call_id 'test-call-id', got '%s'", output.CallID)
	}

	if output.SessionID != "test-session-id" {
		t.Errorf("expected session_id 'test-session-id', got '%s'", output.SessionID)
	}

	if output.SDPOffer != "test-sdp-offer" {
		t.Errorf("expected sdp_offer 'test-sdp-offer', got '%s'", output.SDPOffer)
	}

	if output.Status != "connecting" {
		t.Errorf("expected status 'connecting', got '%s'", output.Status)
	}

	if mockRepo.createdCall == nil {
		t.Fatal("expected call to be created in repository")
	}

	if mockRepo.createdCall.SessionID != "test-session-id" {
		t.Errorf("expected created call session_id 'test-session-id', got '%s'", mockRepo.createdCall.SessionID)
	}
}

func TestInitiateCallUseCase_Execute_MissingUserID(t *testing.T) {
	mockRepo := &mockCallRepository{}
	mockVoIP := &mockVoIPService{}

	uc := NewInitiateCallUseCase(mockRepo, mockVoIP)

	input := InitiateCallInput{
		UserID:      "",
		PhoneNumber: "+491512345678",
	}

	output, err := uc.Execute(context.Background(), input)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if output != nil {
		t.Errorf("expected nil output, got %v", output)
	}

	if err.Error() != "user_id is required" {
		t.Errorf("expected error 'user_id is required', got '%s'", err.Error())
	}
}

func TestInitiateCallUseCase_Execute_MissingPhoneNumber(t *testing.T) {
	mockRepo := &mockCallRepository{}
	mockVoIP := &mockVoIPService{}

	uc := NewInitiateCallUseCase(mockRepo, mockVoIP)

	input := InitiateCallInput{
		UserID:      "test-user-id",
		PhoneNumber: "",
	}

	output, err := uc.Execute(context.Background(), input)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if output != nil {
		t.Errorf("expected nil output, got %v", output)
	}

	if err.Error() != "phone_number is required" {
		t.Errorf("expected error 'phone_number is required', got '%s'", err.Error())
	}
}

func TestInitiateCallUseCase_Execute_VoIPFailure(t *testing.T) {
	mockRepo := &mockCallRepository{}
	mockVoIP := &mockVoIPService{
		initiateError: errors.New("voip service unavailable"),
	}

	uc := NewInitiateCallUseCase(mockRepo, mockVoIP)

	input := InitiateCallInput{
		UserID:      "test-user-id",
		PhoneNumber: "+491512345678",
	}

	output, err := uc.Execute(context.Background(), input)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if output != nil {
		t.Errorf("expected nil output, got %v", output)
	}

	if err.Error() != "failed to initiate call" {
		t.Errorf("expected error 'failed to initiate call', got '%s'", err.Error())
	}
}

func TestInitiateCallUseCase_Execute_RepositoryFailure(t *testing.T) {
	mockRepo := &mockCallRepository{
		createError: errors.New("database error"),
	}
	mockVoIP := &mockVoIPService{
		session: &domain.CallSession{
			SessionID:   "test-session-id",
			PhoneNumber: "+491512345678",
			SDPOffer:    "test-sdp-offer",
			Status:      domain.SessionStatusInitialized,
			CreatedAt:   time.Now(),
			ExpiresAt:   time.Now().Add(5 * time.Minute),
		},
	}

	uc := NewInitiateCallUseCase(mockRepo, mockVoIP)

	input := InitiateCallInput{
		UserID:      "test-user-id",
		PhoneNumber: "+491512345678",
	}

	output, err := uc.Execute(context.Background(), input)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if output != nil {
		t.Errorf("expected nil output, got %v", output)
	}

	if err.Error() != "failed to create call record" {
		t.Errorf("expected error 'failed to create call record', got '%s'", err.Error())
	}
}

