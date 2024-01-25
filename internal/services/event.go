package services

import (
	"database/sql"
	"errors"
	"time"
)

type Event struct {
	Id          int
	AccountId   int
	Name        string
	Description string
	Income      int
	Reserved    int
	DeliveredAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type EventService interface {
	New(user *Event) (*Event, error)
	GetById(id int) (*Event, error)
	GetByAccountId(accountId int) ([]*Event, error)
	Update(event *Event) (*Event, error)
	Delete(id int) error
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
func (s *eventService) New(event *Event) (*Event, error) {
	now := time.Now().UTC()

	mutation, mutationErr := s.db.Exec(
		`INSERT INTO event (
			account_id,
			name,
			description,
			income,
			reserved,
			delivered_at,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?);`,
		event.AccountId,
		event.Name,
		event.Description,
		event.Income,
		event.Reserved,
		event.DeliveredAt,
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

	newEvent := &Event{
		Id:          int(id),
		AccountId:   event.AccountId,
		Name:        event.Name,
		Description: event.Description,
		Income:      event.Income,
		Reserved:    event.Reserved,
		DeliveredAt: event.DeliveredAt,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return newEvent, nil
}

/*
GetById is a function that returns an event by id.
*/
func (s *eventService) GetById(id int) (*Event, error) {
	row := s.db.QueryRow(
		`SELECT
			id,
			account_id,
			name,
			description,
			income,
			reserved,
			delivered_at,
			created_at,
			updated_at
		FROM event
		WHERE id = ?;`,
		id,
	)

	event := &Event{}

	err := row.Scan(
		&event.Id,
		&event.AccountId,
		&event.Name,
		&event.Description,
		&event.Income,
		&event.Reserved,
		&event.DeliveredAt,
		&event.CreatedAt,
		&event.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return event, nil
}

/*
GetByAccountId is a function that returns all events for an account.
*/
func (s *eventService) GetByAccountId(accountId int) ([]*Event, error) {
	rows, err := s.db.Query(
		`SELECT
			id,
			account_id,
			name,
			description,
			income,
			reserved,
			delivered_at,
			created_at,
			updated_at
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

		err := rows.Scan(
			&event.Id,
			&event.AccountId,
			&event.Name,
			&event.Description,
			&event.Income,
			&event.Reserved,
			&event.DeliveredAt,
			&event.CreatedAt,
			&event.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

/*
Update is a function that updates an event in the database.
*/
func (s *eventService) Update(event *Event) (*Event, error) {
	mutation, mutationErr := s.db.Exec(
		`UPDATE event
		SET
			name = ?,
			description = ?,
			income = ?,
			reserved = ?,
			delivered_at = ?,
			updated_at = ?
		WHERE id = ?;`,
		event.Name,
		event.Description,
		event.Income,
		event.Reserved,
		event.DeliveredAt,
		time.Now().UTC(),
		event.Id,
	)

	if mutationErr != nil {
		return nil, mutationErr
	}

	rowsAffected, rowsAffectedErr := mutation.RowsAffected()
	if rowsAffectedErr != nil {
		return nil, rowsAffectedErr
	}

	if rowsAffected == 0 {
		return nil, errors.New("No rows affected")
	}

	return event, nil
}

/*
Delete is a function that deletes an event from the database.
*/
func (s *eventService) Delete(id int) error {
	mutation, mutationErr := s.db.Exec(
		`DELETE FROM event
		WHERE id = ?;`,
		id,
	)

	if mutationErr != nil {
		return mutationErr
	}

	rowsAffected, rowsAffectedErr := mutation.RowsAffected()
	if rowsAffectedErr != nil {
		return rowsAffectedErr
	}

	if rowsAffected == 0 {
		return errors.New("No rows affected")
	}

	return nil
}
