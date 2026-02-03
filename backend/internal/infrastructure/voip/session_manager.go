package voip

import (
	"log/slog"
	"sync"
	"time"

	"github.com/Nikita-Smirnov-idk/Browser-International-Calls-Platform/backend/internal/domain"
)

type SessionManager struct {
	sessions map[string]*domain.CallSession
	mu       sync.RWMutex
	stopChan chan struct{}
}

func NewSessionManager() *SessionManager {
	sm := &SessionManager{
		sessions: make(map[string]*domain.CallSession),
		stopChan: make(chan struct{}),
	}

	go sm.cleanupExpiredSessions()

	return sm
}

func (sm *SessionManager) AddSession(session *domain.CallSession) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.sessions[session.SessionID] = session
	slog.Debug("session added", "session_id", session.SessionID)
}

func (sm *SessionManager) GetSession(sessionID string) *domain.CallSession {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	return sm.sessions[sessionID]
}

func (sm *SessionManager) RemoveSession(sessionID string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	delete(sm.sessions, sessionID)
	slog.Debug("session removed", "session_id", sessionID)
}

func (sm *SessionManager) Close() {
	close(sm.stopChan)
}

func (sm *SessionManager) cleanupExpiredSessions() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			sm.removeExpiredSessions()
		case <-sm.stopChan:
			return
		}
	}
}

func (sm *SessionManager) removeExpiredSessions() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	now := time.Now()
	expiredCount := 0

	for sessionID, session := range sm.sessions {
		if now.After(session.ExpiresAt) {
			delete(sm.sessions, sessionID)
			expiredCount++
		}
	}

	if expiredCount > 0 {
		slog.Info("expired sessions removed", "count", expiredCount)
	}
}

