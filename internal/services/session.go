package services

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"
)

type Session struct {
	Id         int
	UserId     int
	ValidUntil time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type SessionServiceInterface interface {
	New(user *User) (*Session, error)
	GetActiveSessions() ([]*Session, error)
	GetById(sessionId int) (*Session, error)
	GetByUserID(userId int) (*Session, error)
	Delete(sessionId int) error
	CheckCookie(r *http.Request) (*Session, error)
}

type sessionService struct {
	db *sql.DB
}

func NewSessionService(db *sql.DB) SessionServiceInterface {
	return &sessionService{db: db}
}

/*
New creates a new session for the given user in the database.
*/
func (s *sessionService) New(user *User) (*Session, error) {
	now := time.Now().UTC()

	validUntil := now.Add(time.Hour * 24 * 7)

	existingSession, existingSessionErr := s.GetByUserID(user.Id)
	if existingSessionErr == nil && existingSession != nil {
		// session already exists, update it
		_, mutationErr := s.db.Exec(
			`UPDATE session
      SET valid_until = ?
      WHERE id = ?`,
			validUntil,
			existingSession.Id,
		)

		if mutationErr != nil {
			return nil, mutationErr
		}

		newSession := &Session{
			Id:         existingSession.Id,
			UserId:     user.Id,
			ValidUntil: validUntil,
			UpdatedAt:  now,
			CreatedAt:  existingSession.CreatedAt,
		}

		return newSession, nil
	}

	mutation, mutationErr := s.db.Exec(
		`INSERT INTO session (
      user_id,
      valid_until,
      created_at,
      updated_at
    ) VALUES (?, ?, ?, ?)`,
		user.Id,
		validUntil,
		now,
		now,
	)

	if mutationErr != nil {
		return nil, mutationErr
	}

	id, idErr := mutation.LastInsertId()
	if idErr != nil {
		return nil, idErr
	}

	newSession := &Session{
		Id:         int(id),
		UserId:     user.Id,
		ValidUntil: validUntil,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	return newSession, nil
}

/*
GetActiveSessions returns all active sessions from the database.
*/
func (s *sessionService) GetActiveSessions() ([]*Session, error) {
	rows, queryErr := s.db.Query(
		`SELECT
      id,
      user_id,
      valid_until,
      created_at,
      updated_at
    FROM session
    WHERE valid_until > ?`,
		time.Now().UTC(),
	)

	if queryErr != nil {
		return nil, queryErr
	}

	sessions := []*Session{}

	for rows.Next() {

		session := &Session{}

		scanErr := rows.Scan(
			&session.Id,
			&session.UserId,
			&session.ValidUntil,
			&session.CreatedAt,
			&session.UpdatedAt,
		)

		if scanErr != nil {
			return nil, scanErr
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}

/*
GetById returns the session with the given sessionID from the database.
*/
func (s *sessionService) GetById(sessionId int) (*Session, error) {
	row := s.db.QueryRow(
		`SELECT
      id,
      user_id,
      valid_until,
      created_at,
      updated_at
    FROM session
    WHERE id = ?`,
		sessionId,
	)

	session := &Session{}

	scanErr := row.Scan(
		&session.Id,
		&session.UserId,
		&session.ValidUntil,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if scanErr != nil {
		return nil, scanErr
	}

	return session, nil
}

/*
GetByUserID returns the session with the given userID from the database.
*/
func (s *sessionService) GetByUserID(userId int) (*Session, error) {
	row := s.db.QueryRow(
		`SELECT
      id,
      user_id,
      valid_until,
      created_at,
      updated_at
    FROM session
    WHERE user_id = ?`,
		userId,
	)

	session := &Session{}

	scanErr := row.Scan(
		&session.Id,
		&session.UserId,
		&session.ValidUntil,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if scanErr != nil {
		return nil, scanErr
	}

	return session, nil
}

/*
Delete deletes the session with the given sessionID from the database.
*/
func (s *sessionService) Delete(sessionId int) error {
	_, mutationErr := s.db.Exec(
		`DELETE FROM session
    WHERE id = ?`,
		sessionId,
	)

	if mutationErr != nil {
		return mutationErr
	}

	return nil
}

/*
CheckCookie returns the session from the cookie in the request.
*/
func (s *sessionService) CheckCookie(r *http.Request) (*Session, error) {
	cookie, cookieErr := r.Cookie("session")
	if cookieErr != nil {
		return nil, cookieErr
	}

	sessionId, idErr := strconv.Atoi(cookie.Value)
	if idErr != nil {
		return nil, idErr
	}

	session, sessionErr := s.GetById(sessionId)
	if sessionErr != nil {
		return nil, sessionErr
	}

	return session, nil
}