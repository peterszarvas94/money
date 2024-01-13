package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"pengoe/internal/utils"
	"sync"
	"time"
)

type CSRFToken struct {
  Token string
  Valid time.Time
}

type Session struct {
  User *utils.User
  CSRFToken CSRFToken
}

type SessionManager struct {
  sessions map[int]Session // sessionID -> session
  mu sync.Mutex
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
    sessions: make(map[int]Session),
    mu: sync.Mutex{},
	}
}

func generateCSRFToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (sm *SessionManager) CreateSession(sessionID int, user *utils.User) (*Session, error) {
	csrfToken, tokenErr := generateCSRFToken()
	if tokenErr != nil {
    return &Session{}, tokenErr
	}

  locked := sm.mu.TryLock()
  if !locked {
    return &Session{}, errors.New("Could not lock session manager")
  }

	defer sm.mu.Unlock()

  session := Session{
    CSRFToken: CSRFToken{
      Token: csrfToken,
      Valid: time.Now().Add(1 * time.Hour),
    },
    User: user,
  }

	sm.sessions[sessionID] = session
  return &session, nil
}

func (sm *SessionManager) GetSession(sessionID int) (*Session, error) {
  locked := sm.mu.TryLock()
  if !locked {
    return &Session{}, errors.New("Could not lock session manager")
  }
	defer sm.mu.Unlock()

  session, ok := sm.sessions[sessionID]
  if !ok {
    return &Session{}, errors.New(fmt.Sprintf("Session %d not found", sessionID))
  }

	return &session, nil
}

func (sm *SessionManager) DeleteSession(sessionID int) error {
  locked := sm.mu.TryLock()
  if !locked {
    return errors.New("Could not lock session manager")
  }
  defer sm.mu.Unlock()

  delete(sm.sessions, sessionID)
  return nil
}

func (sm *SessionManager) RenewCSRFToken(sessionID int) (string, error) {
  session, sessionErr := sm.GetSession(sessionID)
  if sessionErr != nil {
    return "", sessionErr
  }

  locked := sm.mu.TryLock()
  if !locked {
    return "", errors.New("Could not lock session manager")
  }
  defer sm.mu.Unlock()

  csrfToken, tokenErr := generateCSRFToken()
  if tokenErr != nil {
    return "", tokenErr
  }

  session.CSRFToken = CSRFToken{
    Token: csrfToken,
    Valid: time.Now().Add(1 * time.Hour),
  }

  sm.sessions[sessionID] = *session
  return csrfToken, nil
}
