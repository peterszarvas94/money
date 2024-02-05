package services

import (
	"database/sql"
	"errors"
	"pengoe/internal/utils"
	"time"
)

type Event struct {
	Id          string
	Name        string
	Description string
	Income      int
	Reserved    int
	DeliveredAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	AccountId   string
}

type EventService interface {
	New(id, name, description string, income, reserved int, deliveredAt time.Time, accountId string) error
	GetById(id string) (*Event, error)
	GetByAccountId(accountId string) ([]*Event, error)
	Update(id, name, description string, income, reserved int, deliveredAt time.Time) error
	Delete(id string) error
}

type eventService struct {
	db *sql.DB
}

func NewEventService(db *sql.DB) EventService {
	return &eventService{db: db}
}

/*
New is a function that adds an event to the database.
*/
func (s *eventService) New(id, name, description string, income, reserved int, deliveredAt time.Time, accountId string) error {
	now := time.Now().UTC()

	_, err := s.db.Exec(
		`INSERT INTO event (
			id,
			name,
			description,
			income,
			reserved,
			delivered_at,
			created_at,
			updated_at,
			account_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);`,
		id,
		name,
		description,
		income,
		reserved,
		deliveredAt,
		now,
		now,
		accountId,
	)

	if err != nil {
		return err
	}

	return nil
}

/*
GetById is a function that returns an event by id.
*/
func (s *eventService) GetById(id string) (*Event, error) {
	row := s.db.QueryRow(
		`SELECT
			id,
			name,
			description,
			income,
			reserved,
			delivered_at,
			created_at,
			updated_at,
			account_id
		FROM event
		WHERE id = ?;`,
		id,
	)

	event := &Event{}

	var deliveredAtStr string
	var createdAtStr string
	var updatedAtStr string

	err := row.Scan(
		&event.Id,
		&event.Name,
		&event.Description,
		&event.Income,
		&event.Reserved,
		&deliveredAtStr,
		&createdAtStr,
		&updatedAtStr,
		&event.AccountId,
	)

	if err != nil {
		return nil, err
	}

	deliveredAt, err := utils.ConvertToTime(deliveredAtStr)
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

	event.DeliveredAt = deliveredAt
	event.CreatedAt = createdAt
	event.UpdatedAt = updatedAt

	return event, nil
}

/*
GetByAccountId is a function that returns all events for an account.
*/
func (s *eventService) GetByAccountId(accountId string) ([]*Event, error) {
	rows, err := s.db.Query(
		`SELECT
			id,
			name,
			description,
			income,
			reserved,
			delivered_at,
			created_at,
			updated_at,
			account_id
		FROM event
		WHERE account_id = ?;`,
		accountId,
	)

	if err != nil {
		return nil, err
	}

	events := []*Event{}

	for rows.Next() {
		event := &Event{}

		var deliveredAtStr string
		var createdAtStr string
		var updatedAtStr string

		err := rows.Scan(
			&event.Id,
			&event.Name,
			&event.Description,
			&event.Income,
			&event.Reserved,
			&deliveredAtStr,
			&createdAtStr,
			&updatedAtStr,
			&event.AccountId,
		)

		if err != nil {
			return nil, err
		}

		deliveredAt, err := utils.ConvertToTime(deliveredAtStr)
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

		event.DeliveredAt = deliveredAt
		event.CreatedAt = createdAt
		event.UpdatedAt = updatedAt

		events = append(events, event)
	}

	return events, nil
}

/*
Update is a function that updates an event in the database.
*/
func (s *eventService) Update(id, name, description string, income, reserved int, deliveredAt time.Time) error {
	mutation, err := s.db.Exec(
		`UPDATE event
		SET
			name = ?,
			description = ?,
			income = ?,
			reserved = ?,
			delivered_at = ?,
			updated_at = ?
		WHERE id = ?;`,
		name,
		description,
		income,
		reserved,
		deliveredAt,
		time.Now().UTC(),
		id,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := mutation.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("No rows affected")
	}

	return nil
}

/*
Delete is a function that deletes an event from the database.
*/
func (s *eventService) Delete(id string) error {
	mutation, err := s.db.Exec(
		`DELETE FROM event
		WHERE id = ?;`,
		id,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := mutation.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("No rows affected")
	}

	return nil
}
