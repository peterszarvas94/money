package services

import (
	"database/sql"
	"time"
)

type Role string

const (
	Admin  Role = "admin"
	Viewer Role = "viewer"
)

type Access struct {
	Id        string
	Role      Role
	CreatedAt time.Time
	UpdatedAt time.Time
	UserId    string
	AccountId string
}

type AccessService interface {
	New(id string, role Role, userId string, accountId string) error
	Check(userId, accountId string) bool
}

type accessService struct {
	db *sql.DB
}

func NewAccessService(db *sql.DB) AccessService {
	return &accessService{db: db}
}

/*
New is a function that adds an access to the database.
*/
func (s *accessService) New(id string, role Role, userId, accountId string) error {
	now := time.Now().UTC()

	_, err := s.db.Exec(
		`INSERT INTO access (
			id,
			role,
			created_at,
			updated_at,
			user_id,
			account_id
		) VALUES (?, ?, ?, ?, ?, ?)`,
		id,
		role,
		now,
		now,
		userId,
		accountId,
	)

	if err != nil {
		return err
	}

	return nil
}

/*
Check is a function that checks if a user has access to an account.
*/
func (s *accessService) Check(userId string, accountId string) bool {
	row := s.db.QueryRow(
		`SELECT COUNT(*) FROM access WHERE user_id = ? AND account_id = ?`,
		userId,
		accountId,
	)

	var count int

	err := row.Scan(&count)
	if err != nil {
		return false
	}

	if count == 0 {
		return false
	}

	return true
}
