package services

import (
	"database/sql"
	"pengoe/internal/utils"
	"time"
)

type Account struct {
	Id          string
	Name        string
	Description string
	Currency    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type AccountServiceInterface interface {
	New(id, name, description, currency string) error
	GetByUserId(userId string) ([]*Account, error)
	GetById(id string) (*Account, error)
	Delete(id string) error
}

type accountService struct {
	db *sql.DB
}

func NewAccountService(db *sql.DB) AccountServiceInterface {
	return &accountService{db: db}
}

/*
New is a function that adds an account to the database.
Gives back the id of the new account.
*/
func (s *accountService) New(id, name, description, currency string) error {
	now := time.Now().UTC()

	_, err := s.db.Exec(
		`INSERT INTO account (
			id,
			name,
			description,
			currency,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?)`,
		id,
		name,
		description,
		currency,
		now,
		now,
	)
	if err != nil {
		return err
	}

	return nil
}

/*
GetByUserId is a function that returns all accounts for a given user.
*/
func (s *accountService) GetByUserId(userId string) ([]*Account, error) {
	rows, err := s.db.Query(
		`SELECT 
			account.id,
			account.name,
			account.description,
			account.currency,
			account.created_at,
			account.updated_at
		FROM account
		INNER JOIN access ON account.id = access.account_id
		WHERE access.user_id = ?`,
		userId,
	)
	if err != nil {
		return nil, err
	}

	accounts := []*Account{}

	for rows.Next() {
		account := &Account{}

		var createdAtStr string
		var updatedAtStr string

		err := rows.Scan(
			&account.Id,
			&account.Name,
			&account.Description,
			&account.Currency,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			return nil, err
		}

		createdAt, err := utils.ConvertToTime(createdAtStr)
		if err != nil {
			return nil, err
		}

		updated, err := utils.ConvertToTime(updatedAtStr)
		if err != nil {
			return nil, err
		}

		account.CreatedAt = createdAt
		account.UpdatedAt = updated

		accounts = append(accounts, account)
	}

	return accounts, nil
}

/*
GetById is a function that returns an account for a given id.
*/
func (s *accountService) GetById(id string) (*Account, error) {
	row := s.db.QueryRow(
		`SELECT
			id,
			name,
			description,
			currency,
			created_at,
			updated_at
		FROM account
		WHERE account.id = ?`,
		id,
	)

	var createdAtStr string
	var updatedAtStr string

	account := &Account{}

	err := row.Scan(
		&account.Id,
		&account.Name,
		&account.Description,
		&account.Currency,
		&createdAtStr,
		&updatedAtStr,
	)
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

	account.CreatedAt = createdAt
	account.UpdatedAt = updatedAt

	return account, nil
}

/*
Delete is a function that deletes an account from the database.
*/
func (s *accountService) Delete(id string) error {
	_, err := s.db.Exec(
		`DELETE FROM account WHERE id = ?`,
		id,
	)
	if err != nil {
		return err
	}

	return nil
}
