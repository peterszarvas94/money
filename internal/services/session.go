package services

import (
	"database/sql"
	"net/http"
	"pengoe/internal/utils"
	"time"
)

type Session struct {
	Id         string
	ValidUntil time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
	UserId     string
}

type SessionServiceInterface interface {
	New(id, userId string) (*Session, error)
	GetActives() ([]*Session, error)
	GetById(id string) (*Session, error)
	GetByUserID(usedId string) (*Session, error)
	Delete(id string) error
	CheckFromCookie(*http.Request) (*Session, error)
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
func (s *sessionService) New(id, userId string) (*Session, error) {
	now := time.Now().UTC()

	validUntil := now.Add(time.Hour * 24 * 7)

	_, err := s.db.Exec(
		`INSERT INTO session (
			id,
			valid_until,
			created_at,
			updated_at,
			user_id
		) VALUES (?, ?, ?, ?, ?)`,
		id,
		validUntil,
		now,
		now,
		userId,
	)

	if err != nil {
		return nil, err
	}

	newSession := &Session{
		Id:         id,
		ValidUntil: validUntil,
		CreatedAt:  now,
		UpdatedAt:  now,
		UserId:     userId,
	}

	return newSession, nil
}

/*
GetActiveSessions returns all active sessions from the database.
*/
func (s *sessionService) GetActives() ([]*Session, error) {
	rows, err := s.db.Query(
		`SELECT
      id,
      valid_until,
      created_at,
      updated_at,
			user_id
    FROM session
    WHERE valid_until > ?`,
		time.Now().UTC(),
	)

	if err != nil {
		return nil, err
	}

	sessions := []*Session{}

	for rows.Next() {
		session := &Session{}
		var validUntilStr string
		var createdAtStr string
		var updatedAtStr string

		err := rows.Scan(
			&session.Id,
			&validUntilStr,
			&createdAtStr,
			&updatedAtStr,
			&session.UserId,
		)

		if err != nil {
			return nil, err
		}

		validUntil, err := utils.ConvertToTime(validUntilStr)
		if err != nil {
			return nil, err
		}

		createdAt, err := utils.ConvertToTime(createdAtStr)
		if err != nil {
			return nil, err
		}

		updatedAt, err := utils.ConvertToTime(updatedAtStr)
		if err != nil {
			return nil, err
		}

		session.ValidUntil = validUntil
		session.CreatedAt = createdAt
		session.UpdatedAt = updatedAt

		sessions = append(sessions, session)
	}

	return sessions, nil
}

/*
GetById returns the session with the given sessionID from the database.
*/
func (s *sessionService) GetById(id string) (*Session, error) {
	row := s.db.QueryRow(
		`SELECT
			id,
			valid_until,
			created_at,
			updated_at,
			user_id
		FROM session
		WHERE id = ?`,
		id,
	)

	session := &Session{}

	var validUntilStr string
	var createdAtStr string
	var updatedAtStr string

	err := row.Scan(
		&session.Id,
		&validUntilStr,
		&createdAtStr,
		&updatedAtStr,
		&session.UserId,
	)

	if err != nil {
		return nil, err
	}

	validUntil, err := utils.ConvertToTime(validUntilStr)
	if err != nil {
		return nil, err
	}

	createdAt, err := utils.ConvertToTime(createdAtStr)
	if err != nil {
		return nil, err
	}

	updatedAt, err := utils.ConvertToTime(updatedAtStr)
	if err != nil {
		return nil, err
	}

	session.ValidUntil = validUntil
	session.CreatedAt = createdAt
	session.UpdatedAt = updatedAt

	return session, nil
}

/*
GetByUserID returns the session with the given userID from the database.
*/
func (s *sessionService) GetByUserID(userId string) (*Session, error) {
	row := s.db.QueryRow(
		`SELECT
			id,
			valid_until,
			created_at,
			updated_at,
			user_id
		FROM session
		WHERE user_id = ?`,
		userId,
	)

	session := &Session{}

	var validUntilStr string
	var createdAtStr string
	var updatedAtStr string

	err := row.Scan(
		&session.Id,
		&validUntilStr,
		&createdAtStr,
		&updatedAtStr,
		&session.UserId,
	)

	if err != nil {
		return nil, err
	}

	validUntil, err := utils.ConvertToTime(validUntilStr)
	if err != nil {
		return nil, err
	}

	createdAt, err := utils.ConvertToTime(createdAtStr)
	if err != nil {
		return nil, err
	}

	updatedAt, err := utils.ConvertToTime(updatedAtStr)
	if err != nil {
		return nil, err
	}

	session.ValidUntil = validUntil
	session.CreatedAt = createdAt
	session.UpdatedAt = updatedAt

	return session, nil
}

/*
Delete deletes the session with the given sessionID from the database.
*/
func (s *sessionService) Delete(id string) error {
	_, err := s.db.Exec(
		`DELETE FROM session
		WHERE id = ?`,
		id,
	)

	if err != nil {
		return err
	}

	return nil
}

/*
CheckCookie returns the session from the cookie in the request.
*/
func (s *sessionService) CheckFromCookie(r *http.Request) (*Session, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return nil, err
	}

	session, err := s.GetById(cookie.Value)
	if err != nil {
		return nil, err
	}

	return session, nil
}
