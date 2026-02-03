package calls

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/domain"
)

type mockCallRepositoryForTerminate struct {
	getError    error
	updateError error
	call        *domain.Call
	updatedCall *domain.Call
}

func (m *mockCallRepositoryForTerminate) Create(ctx context.Context, call *domain.Call) error {
	return nil
}

func (m *mockCallRepositoryForTerminate) Update(ctx context.Context, call *domain.Call) error {
	if m.updateError != nil {
		return m.updateError
	}
	m.updatedCall = call
	return nil
}

func (m *mockCallRepositoryForTerminate) GetByID(ctx context.Context, id string) (*domain.Call, error) {
	if m.getError != nil {
		return nil, m.getError
	}
	return m.call, nil
}

func (m *mockCallRepositoryForTerminate) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*domain.Call, error) {
	return nil, nil
}

type mockVoIPServiceForTerminate struct {
	terminateError error
}

func (m *mockVoIPServiceForTerminate) InitiateCall(ctx context.Context, phoneNumber string) (*domain.CallSession, error) {
	return nil, nil
}

func (m *mockVoIPServiceForTerminate) TerminateCall(ctx context.Context, sessionID string) error {
	return m.terminateError
}

func (m *mockVoIPServiceForTerminate) GetSessionStatus(ctx context.Context, sessionID string) (domain.SessionStatus, error) {
	return domain.SessionStatusActive, nil
}

func TestTerminateCallUseCase_Execute_Success(t *testing.T) {
	startTime := time.Now().Add(-30 * time.Second)
	mockRepo := &mockCallRepositoryForTerminate{
		call: &domain.Call{
			ID:          "test-call-id",
			UserID:      "test-user-id",
			PhoneNumber: "+491512345678",
			StartTime:   startTime,
			Status:      domain.CallStatusActive,
			SessionID:   "test-session-id",
		},
	}
	mockVoIP := &mockVoIPServiceForTerminate{}

	uc := NewTerminateCallUseCase(mockRepo, mockVoIP)

	input := TerminateCallInput{
		UserID: "test-user-id",
		CallID: "test-call-id",
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

	if output.Status != "completed" {
		t.Errorf("expected status 'completed', got '%s'", output.Status)
	}

	if output.Duration <= 0 {
		t.Errorf("expected positive duration, got %d", output.Duration)
	}

	if mockRepo.updatedCall == nil {
		t.Fatal("expected call to be updated in repository")
	}

	if mockRepo.updatedCall.Status != domain.CallStatusCompleted {
		t.Errorf("expected updated call status 'completed', got '%s'", mockRepo.updatedCall.Status)
	}
}

func TestTerminateCallUseCase_Execute_MissingCallID(t *testing.T) {
	mockRepo := &mockCallRepositoryForTerminate{}
	mockVoIP := &mockVoIPServiceForTerminate{}

	uc := NewTerminateCallUseCase(mockRepo, mockVoIP)

	input := TerminateCallInput{
		UserID: "test-user-id",
		CallID: "",
	}

	output, err := uc.Execute(context.Background(), input)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if output != nil {
		t.Errorf("expected nil output, got %v", output)
	}

	if err.Error() != "call_id is required" {
		t.Errorf("expected error 'call_id is required', got '%s'", err.Error())
	}
}

func TestTerminateCallUseCase_Execute_MissingUserID(t *testing.T) {
	mockRepo := &mockCallRepositoryForTerminate{}
	mockVoIP := &mockVoIPServiceForTerminate{}

	uc := NewTerminateCallUseCase(mockRepo, mockVoIP)

	input := TerminateCallInput{
		UserID: "",
		CallID: "test-call-id",
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

func TestTerminateCallUseCase_Execute_CallNotFound(t *testing.T) {
	mockRepo := &mockCallRepositoryForTerminate{
		call: nil,
	}
	mockVoIP := &mockVoIPServiceForTerminate{}

	uc := NewTerminateCallUseCase(mockRepo, mockVoIP)

	input := TerminateCallInput{
		UserID: "test-user-id",
		CallID: "test-call-id",
	}

	output, err := uc.Execute(context.Background(), input)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if output != nil {
		t.Errorf("expected nil output, got %v", output)
	}

	if err.Error() != "call not found" {
		t.Errorf("expected error 'call not found', got '%s'", err.Error())
	}
}

func TestTerminateCallUseCase_Execute_Unauthorized(t *testing.T) {
	mockRepo := &mockCallRepositoryForTerminate{
		call: &domain.Call{
			ID:          "test-call-id",
			UserID:      "other-user-id",
			PhoneNumber: "+491512345678",
			StartTime:   time.Now().Add(-30 * time.Second),
			Status:      domain.CallStatusActive,
			SessionID:   "test-session-id",
		},
	}
	mockVoIP := &mockVoIPServiceForTerminate{}

	uc := NewTerminateCallUseCase(mockRepo, mockVoIP)

	input := TerminateCallInput{
		UserID: "test-user-id",
		CallID: "test-call-id",
	}

	output, err := uc.Execute(context.Background(), input)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if output != nil {
		t.Errorf("expected nil output, got %v", output)
	}

	if err.Error() != "unauthorized" {
		t.Errorf("expected error 'unauthorized', got '%s'", err.Error())
	}
}

func TestTerminateCallUseCase_Execute_RepositoryFailure(t *testing.T) {
	mockRepo := &mockCallRepositoryForTerminate{
		call: &domain.Call{
			ID:          "test-call-id",
			UserID:      "test-user-id",
			PhoneNumber: "+491512345678",
			StartTime:   time.Now().Add(-30 * time.Second),
			Status:      domain.CallStatusActive,
			SessionID:   "test-session-id",
		},
		updateError: errors.New("database error"),
	}
	mockVoIP := &mockVoIPServiceForTerminate{}

	uc := NewTerminateCallUseCase(mockRepo, mockVoIP)

	input := TerminateCallInput{
		UserID: "test-user-id",
		CallID: "test-call-id",
	}

	output, err := uc.Execute(context.Background(), input)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if output != nil {
		t.Errorf("expected nil output, got %v", output)
	}

	if err.Error() != "failed to update call" {
		t.Errorf("expected error 'failed to update call', got '%s'", err.Error())
	}
}

func TestTerminateCallUseCase_Execute_VoIPFailure_ContinuesAnyway(t *testing.T) {
	startTime := time.Now().Add(-30 * time.Second)
	mockRepo := &mockCallRepositoryForTerminate{
		call: &domain.Call{
			ID:          "test-call-id",
			UserID:      "test-user-id",
			PhoneNumber: "+491512345678",
			StartTime:   startTime,
			Status:      domain.CallStatusActive,
			SessionID:   "test-session-id",
		},
	}
	mockVoIP := &mockVoIPServiceForTerminate{
		terminateError: errors.New("voip service unavailable"),
	}

	uc := NewTerminateCallUseCase(mockRepo, mockVoIP)

	input := TerminateCallInput{
		UserID: "test-user-id",
		CallID: "test-call-id",
	}

	output, err := uc.Execute(context.Background(), input)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output == nil {
		t.Fatal("expected output, got nil")
	}

	if mockRepo.updatedCall == nil {
		t.Fatal("expected call to be updated despite VoIP failure")
	}
}

