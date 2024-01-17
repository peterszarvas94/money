package serversession

import (
	"errors"
	"fmt"
	"os"
	"pengoe/internal/db"
	"pengoe/internal/logger"
	"pengoe/internal/services"
	"pengoe/internal/utils"
	"sync"
	"time"
)

type CSRFToken struct {
	Token string
	Valid time.Time
}

type serverSession struct {
	SessionID int
	CSRFToken CSRFToken
}

type serverSessionInterface interface {
	Create(sessionID int) (*serverSession, error)
	Get(sessionID int) (*serverSession, error)
	Delete(sessionID int) error
	RenewCSRFToken(sessionID int) (*CSRFToken, error)
}

type serverSessionManager struct {
	sessions map[int]serverSession // sessionID -> session
	mutex    sync.Mutex
}

func newServerSessionManager() serverSessionInterface {
	return &serverSessionManager{
		sessions: make(map[int]serverSession),
		mutex:    sync.Mutex{},
	}
}

var Manager serverSessionInterface = newServerSessionManager()

/*
Create creates a new server session in memory. Not to be confused with
the database session, it is only for CSRF protection. Server session and
database session are linked by the sessionID, they should be in sync.
*/
func (sm *serverSessionManager) Create(sessionID int) (*serverSession, error) {
	csrfToken, tokenErr := utils.GenerateCSRFToken()
	if tokenErr != nil {
		return &serverSession{}, tokenErr
	}

	locked := sm.mutex.TryLock()
	if !locked {
		return &serverSession{}, errors.New("Could not lock session manager")
	}

	defer sm.mutex.Unlock()

	session := serverSession{
		CSRFToken: CSRFToken{
			Token: csrfToken,
			Valid: time.Now().Add(time.Minute).UTC(),
		},
	}

	sm.sessions[sessionID] = session

	return &session, nil
}

/*
Get returns the session with the given sessionID from memory.
*/
func (sm *serverSessionManager) Get(sessionID int) (*serverSession, error) {
	locked := sm.mutex.TryLock()
	if !locked {
		return &serverSession{}, errors.New("Could not lock session manager")
	}
	defer sm.mutex.Unlock()

	session, ok := sm.sessions[sessionID]
	if !ok {
		return &serverSession{}, errors.New(fmt.Sprintf("Session %d not found", sessionID))
	}

	return &session, nil
}

/*
Delete deletes the session with the given sessionID from memory.
*/
func (sm *serverSessionManager) Delete(sessionID int) error {
	locked := sm.mutex.TryLock()
	if !locked {
		return errors.New("Could not lock session manager")
	}
	defer sm.mutex.Unlock()

	delete(sm.sessions, sessionID)

	return nil
}

/*
RenewCSRFToken generates a new CSRF token for the given sessionID.
*/
func (sm *serverSessionManager) RenewCSRFToken(sessionID int) (*CSRFToken, error) {
	session, sessionErr := sm.Get(sessionID)
	if sessionErr != nil {
		return nil, sessionErr
	}

	locked := sm.mutex.TryLock()
	if !locked {
		return nil, errors.New("Could not lock session manager")
	}
	defer sm.mutex.Unlock()

	newToken, tokenErr := utils.GenerateCSRFToken()
	if tokenErr != nil {
		return nil, tokenErr
	}

	newCsrfToken := &CSRFToken{
		Token: newToken,
		Valid: time.Now().Add(1 * time.Hour).UTC(),
	}

	session.CSRFToken = *newCsrfToken

	sm.sessions[sessionID] = *session

	return newCsrfToken, nil
}

func init() {
	// connect to the database
	dbManager := db.NewDBManager()
	db, dbErr := dbManager.GetDB()
	if dbErr != nil {
		logger.Log(logger.FATAL, "serversession/init/db", dbErr.Error())
		os.Exit(1)
	}
	defer db.Close()
	sessionService := services.NewSessionService(db)

	// get all active sessions from the database
	sessions, sessionsErr := sessionService.GetActiveSessions()
	if sessionsErr != nil {
		logger.Log(logger.FATAL, "serversession/init/active", sessionsErr.Error())
		os.Exit(1)
	}

	// create a new server session for each active session
	for _, session := range sessions {
		_, createErr := Manager.Create(session.Id)
		if createErr != nil {
			logger.Log(logger.FATAL, "serversession/init/create", createErr.Error())
			os.Exit(1)
		}
	}

	logger.Log(logger.INFO, "serversession/init", "Server session manager initialized")
}
