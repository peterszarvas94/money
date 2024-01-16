package services

import (
	"database/sql"
	"pengoe/internal/utils"
	"time"
)

type EventService interface {
	New(user *utils.Event) (*utils.Event, error)
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
func (s *eventService) New(event *utils.Event) (*utils.Event, error) {
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
			updated_at,
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
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

	newEvent := &utils.Event{
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

// TODO:
// 1. simplify some handler logic, reuse code pls (new account form, delete btn)
// 2. new event from on account page (copmonent, can be added more than one)
