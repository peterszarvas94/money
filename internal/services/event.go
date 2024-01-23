package services

import (
	"database/sql"
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
	GetByAccountId(accountId int) ([]*Event, error)
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
