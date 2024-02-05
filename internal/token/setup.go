package token

import (
	"errors"
	"fmt"
	"net/http"
	"pengoe/internal/db"
	"pengoe/internal/logger"
	"pengoe/internal/services"
	"pengoe/internal/utils"
	"sync"
	"time"
)

type Token struct {
	SessionID string
	Value     string
	Valid     time.Time
}

type tokenManagerInterface interface {
	Create(sessionId string) (*Token, error)
	Get(sessionId string) (*Token, error)
	Delete(sessionId string) error
	RenewToken(sessionId string) (*Token, error)
	VerifyOrRenewCSRFToken(sessionId string, tokenFromRequest string) (*Token, error)
}

type TokenManager struct {
	tokens map[string]Token // sessionID -> token
	mutex  sync.Mutex
}

func newManager() tokenManagerInterface {
	return &TokenManager{
		tokens: make(map[string]Token),
		mutex:  sync.Mutex{},
	}
}

var Manager tokenManagerInterface = newManager()

/*
Create creates a new token in memory. Server session and
tokens are linked by the sessionID, they should be in sync.
*/
func (m *TokenManager) Create(sessionID string) (*Token, error) {
	csrfToken, tokenErr := utils.GenerateCSRFToken()
	if tokenErr != nil {
		return &Token{}, tokenErr
	}

	locked := m.mutex.TryLock()
	if !locked {
		return &Token{}, errors.New("Could not lock token manager")
	}

	defer m.mutex.Unlock()

	token := Token{
		SessionID: sessionID,
		Value:     csrfToken,
		// TODO: change to 9 minutes
		Valid: time.Now().Add(10 * time.Second).UTC(),
	}

	m.tokens[sessionID] = token

	return &token, nil
}

/*
Get returns the token with the given sessionID.
*/
func (m *TokenManager) Get(sessionId string) (*Token, error) {
	locked := m.mutex.TryLock()
	if !locked {
		return &Token{}, errors.New("Could not lock token manager")
	}
	defer m.mutex.Unlock()

	token, ok := m.tokens[sessionId]
	if !ok {
		return &Token{}, errors.New(fmt.Sprintf("Session %s not found", sessionId))
	}

	return &token, nil
}

/*
Delete deletes the token with the given sessionID.
*/
func (m *TokenManager) Delete(sessionId string) error {
	locked := m.mutex.TryLock()
	if !locked {
		return errors.New("Could not lock token manager")
	}
	defer m.mutex.Unlock()

	delete(m.tokens, sessionId)

	return nil
}

/*
RenewCSRFToken generates a new token for the given sessionID.
*/
func (m *TokenManager) RenewToken(sessionId string) (*Token, error) {
	token, tokenErr := m.Get(sessionId)
	if tokenErr != nil {
		return nil, tokenErr
	}

	locked := m.mutex.TryLock()
	if !locked {
		return nil, errors.New("Could not lock token manager")
	}
	defer m.mutex.Unlock()

	newToken, newTokenErr := utils.GenerateCSRFToken()
	if newTokenErr != nil {
		return nil, newTokenErr
	}

	newCsrfToken := &Token{
		SessionID: sessionId,
		Value:     newToken,
		// TODO: change to 10 minutes
		Valid: time.Now().Add(10 * time.Second).UTC(),
	}

	token = newCsrfToken

	m.tokens[sessionId] = *token

	return newCsrfToken, nil
}

/*
GetOrRenewCSRFToken returns a new token or nothing if the token is valid.
*/
func (m *TokenManager) VerifyOrRenewCSRFToken(sessionId string, tokenFromRequest string) (*Token, error) {
	// check if there is a server session
	token, tokenErr := m.Get(sessionId)
	if tokenErr != nil {
		return nil, tokenErr
	}

	// check if the csrf token is valid
	if token.Value != tokenFromRequest {
		return nil, errors.New("CSRF token is invalid")
	}

	// check if the csrf token is expired
	if token.Valid.Before(time.Now().UTC()) {
		// csrf token is expired

		// renew csrf token
		newToken, tokenErr := Manager.RenewToken(sessionId)
		if tokenErr != nil {
			return nil, tokenErr
		}

		return newToken, nil
	}

	// csrf token is valid
	return nil, nil
}

func GetSessionFromCookie(r *http.Request) (*Token, error) {
	cookie, cookieErr := r.Cookie("session")
	if cookieErr != nil {
		return nil, cookieErr
	}

	token, tokenErr := Manager.Get(cookie.Value)
	if tokenErr != nil {
		return nil, tokenErr
	}

	return token, nil
}

func init() {
	log := logger.Get()

	// connect to the database
	dbManager := db.NewDBManager()
	db, dbErr := dbManager.GetDB()
	if dbErr != nil {
		log.Fatal(dbErr.Error())
	}
	defer db.Close()
	sessionService := services.NewSessionService(db)

	// get all active sessions from the database
	sessions, err := sessionService.GetActives()
	if err != nil {
		log.Fatal(err.Error())
	}

	// create a new token for each active session
	for _, session := range sessions {
		_, err = Manager.Create(session.Id)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}
