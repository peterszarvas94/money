package services

import (
	"database/sql"
	"errors"
	"pengoe/internal/utils"
	"time"
)

type AccessService interface {
	New(user *utils.Access) (*utils.Access, error)
	Check(userId int, accountId int) error
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
func (s *accessService) New(access *utils.Access) (*utils.Access, error) {
	now := time.Now().UTC()

	mutation, mutationErr := s.db.Exec(
		`INSERT INTO access (
			role, user_id, account_id, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?)`,
		access.Role, access.UserId, access.AccountId, now, now,
	)

	if mutationErr != nil {
		return nil, mutationErr
	}

	id, idErr := mutation.LastInsertId()
	if idErr != nil {
		return nil, idErr
	}

	newAccess := &utils.Access{
		Id:        int(id),
		Role:      access.Role,
		AccountId: access.AccountId,
		UserId:    access.UserId,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return newAccess, nil
}

/*
Check is a function that checks if a user has access to an account.
*/
func (s *accessService) Check(userId int, accountId int) error {
	var count int

	row := s.db.QueryRow(
		`SELECT COUNT(*) FROM access WHERE user_id = ? AND account_id = ?`,
		userId,
		accountId,
	)

	err := row.Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("No access")
	}

	return nil
}
