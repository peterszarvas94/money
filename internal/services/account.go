package services

import (
	"database/sql"
	"pengoe/internal/utils"
	"time"
)

type AccountService interface {
	New(user *utils.Account) (*utils.Account, error)
	GetByUserId(userId int) ([]*utils.Account, error)
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

/*
GetByUserId is a function that returns all accounts for a given user.
*/
func (s *accountService) GetByUserId(userId int) ([]*utils.Account, error) {
	rows, rowsErr := s.db.Query(
		`SELECT 
			account.id, account.name, account.description, account.currency, account.created_at, account.updated_at
		FROM account
		INNER JOIN access ON account.id = access.account_id
		WHERE access.user_id = ?`,
		userId,
	)
	if rowsErr != nil {
		return nil, rowsErr
	}

	accounts := []*utils.Account{}

	for rows.Next() {
		var id int
		var name string
		var description string
		var currency string
		var createdAt string
		var updatedAt string

		scanErr := rows.Scan(&id, &name, &description, &currency, &createdAt, &updatedAt)
		if scanErr != nil {
			return nil, scanErr
		}

		account := &utils.Account{
			Id:          id,
			Name:        name,
			Description: description,
			Currency:    currency,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}
