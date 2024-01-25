package services

import (
	"database/sql"
	"pengoe/internal/utils"
	"time"
)

type Account struct {
	Id          int
	Name        string
	Description string
	Currency    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type AccountServiceInterface interface {
	New(user *Account) (*Account, error)
	GetByUserId(userId int) ([]*Account, error)
	GetById(accountId int) (*Account, error)
	Delete(accountId int) error
}

type accountService struct {
	db *sql.DB
}

func NewAccountService(db *sql.DB) AccountServiceInterface {
	return &accountService{db: db}
}

/*
New is a function that adds an account to the database.
*/
func (s *accountService) New(account *Account) (*Account, error) {

	now := time.Now().UTC()

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

	newAccount := &Account{
		Id:          int(id),
		Name:        account.Name,
		Description: account.Description,
		Currency:    account.Currency,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return newAccount, nil
}

/*
GetByUserId is a function that returns all accounts for a given user.
*/
func (s *accountService) GetByUserId(userId int) ([]*Account, error) {
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

	accounts := []*Account{}

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

		created, createdErr := utils.ConvertToTime(createdAt)
		if createdErr != nil {
			return nil, createdErr
		}

		updated, updatedErr := utils.ConvertToTime(updatedAt)
		if updatedErr != nil {
			return nil, updatedErr
		}

		account := &Account{
			Id:          id,
			Name:        name,
			Description: description,
			Currency:    currency,
			CreatedAt:   created,
			UpdatedAt:   updated,
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

/*
GetById is a function that returns an account for a given id.
*/
func (s *accountService) GetById(accountId int) (*Account, error) {
	row := s.db.QueryRow(
		`SELECT
			account.id, account.name, account.description, account.currency, account.created_at, account.updated_at
		FROM account
		WHERE account.id = ?`,
		accountId,
	)

	var id int
	var name string
	var description string
	var currency string
	var createdAt string
	var updatedAt string

	scanErr := row.Scan(&id, &name, &description, &currency, &createdAt, &updatedAt)
	if scanErr != nil {
		return nil, scanErr
	}

	created, createdErr := utils.ConvertToTime(createdAt)
	if createdErr != nil {
		return nil, createdErr
	}

	updated, updatedErr := utils.ConvertToTime(updatedAt)
	if updatedErr != nil {
		return nil, updatedErr
	}

	account := &Account{
		Id:          id,
		Name:        name,
		Description: description,
		Currency:    currency,
		CreatedAt:   created,
		UpdatedAt:   updated,
	}

	return account, nil
}

/*
Delete is a function that deletes an account from the database.
*/
func (s *accountService) Delete(accountId int) error {
	_, deleteErr := s.db.Exec(
		`DELETE FROM account WHERE id = ?`,
		accountId,
	)
	if deleteErr != nil {
		return deleteErr
	}

	return nil
}
