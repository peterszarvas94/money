package services

import (
	"database/sql"
	"pengoe/internal/utils"
	"time"
)

type AccessService interface {
	New(user *utils.Access) (*utils.Access, error)
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
	now := time.Now()

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
		CreatedAt: now.String(),
		UpdatedAt: now.String(),
	}

	return newAccess, nil
}
