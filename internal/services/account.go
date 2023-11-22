package services

import (
	"database/sql"
	"pengoe/internal/utils"
	"time"
)

type AccountService interface {
	New(user *utils.Account) (*utils.Account, error)
}

type accountService struct {
	db *sql.DB
}

func NewAccountService(db *sql.DB) AccountService {
	return &accountService{db: db}
}

/*
New is a function that adds an account to the database.
*/
func (s *accountService) New(account *utils.Account) (*utils.Account, error) {

	now := time.Now()

	mutation, mutationErr := s.db.Exec(
		`INSERT INTO account (
			name, description, currency, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?)`,
		account.Name, account.Description, account.Currency, now, now,
	)
	if mutationErr != nil {
		return nil, mutationErr
	}

	id, idErr := mutation.LastInsertId()
	if idErr != nil {
		return nil, idErr
	}

	newAccount := &utils.Account{
		Id:          int(id),
		Name:        account.Name,
		Description: account.Description,
		Currency:    account.Currency,
		CreatedAt:   now.String(),
		UpdatedAt:   now.String(),
	}

	return newAccount, nil
}
