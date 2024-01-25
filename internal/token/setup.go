package token

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"pengoe/internal/db"
	"pengoe/internal/logger"
	"pengoe/internal/services"
	"pengoe/internal/utils"
	"strconv"
	"sync"
	"time"
)

type Token struct {
	SessionID int
	Value     string
	Valid     time.Time
}

type tokenManagerInterface interface {
	Create(sessionID int) (*Token, error)
	Get(sessionID int) (*Token, error)
	Delete(sessionID int) error
	RenewToken(sessionID int) (*Token, error)
	VerifyOrRenewCSRFToken(sessionId int, tokenFromRequest string) (*Token, error)
}

type TokenManager struct {
	tokens map[int]Token // sessionID -> token
	mutex  sync.Mutex
}

func newManager() tokenManagerInterface {
	return &TokenManager{
		tokens: make(map[int]Token),
		mutex:  sync.Mutex{},
	}
}

var Manager tokenManagerInterface = newManager()

/*
Create creates a new token in memory. Server session and
tokens are linked by the sessionID, they should be in sync.
*/
func (m *TokenManager) Create(sessionID int) (*Token, error) {
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
<<<<<<< HEAD
=======
		// TODO: change to 9 minutes
>>>>>>> 7aab448 (add: events to account)
		Valid:     time.Now().Add(10 * time.Second).UTC(),
	}

	m.tokens[sessionID] = token

	return &token, nil
}

/*
Get returns the token with the given sessionID.
*/
func (m *TokenManager) Get(sessionID int) (*Token, error) {
	locked := m.mutex.TryLock()
	if !locked {
		return &Token{}, errors.New("Could not lock token manager")
	}
	defer m.mutex.Unlock()

	token, ok := m.tokens[sessionID]
	if !ok {
		return &Token{}, errors.New(fmt.Sprintf("Session %d not found", sessionID))
	}

	return &token, nil
}

/*
Delete deletes the token with the given sessionID.
*/
func (m *TokenManager) Delete(sessionID int) error {
	locked := m.mutex.TryLock()
	if !locked {
		return errors.New("Could not lock token manager")
	}
	defer m.mutex.Unlock()

	delete(m.tokens, sessionID)

	return nil
}

/*
RenewCSRFToken generates a new token for the given sessionID.
*/
func (m *TokenManager) RenewToken(sessionID int) (*Token, error) {
	token, tokenErr := m.Get(sessionID)
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
		SessionID: sessionID,
		Value:     newToken,
<<<<<<< HEAD
		// TODO
=======
		// TODO: change to 10 minutes
>>>>>>> 7aab448 (add: events to account)
		Valid:     time.Now().Add(10 * time.Second).UTC(),
	}

	token = newCsrfToken

	m.tokens[sessionID] = *token

	return newCsrfToken, nil
}

/*
GetOrRenewCSRFToken returns a new token or nothing if the token is valid.
*/
func (m *TokenManager) VerifyOrRenewCSRFToken(sessionId int, tokenFromRequest string) (*Token, error) {
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

	sessionId, idErr := strconv.Atoi(cookie.Value)
	if idErr != nil {
		return nil, idErr
	}

	token, tokenErr := Manager.Get(sessionId)
	if tokenErr != nil {
		return nil, tokenErr
	}

	return token, nil
}

func init() {
	// connect to the database
	dbManager := db.NewDBManager()
	db, dbErr := dbManager.GetDB()
	if dbErr != nil {
		logger.Log(logger.FATAL, "csrftoken/init/db", dbErr.Error())
		os.Exit(1)
	}
	defer db.Close()
	sessionService := services.NewSessionService(db)

	// get all active sessions from the database
	sessions, sessionsErr := sessionService.GetActiveSessions()
	if sessionsErr != nil {
		logger.Log(logger.FATAL, "csrftoken/init/active", sessionsErr.Error())
		os.Exit(1)
	}

	// create a new token for each active session
	for _, session := range sessions {
		_, createErr := Manager.Create(session.Id)
		if createErr != nil {
			logger.Log(logger.FATAL, "csrftoken/init/create", createErr.Error())
			os.Exit(1)
		}
	}
}
