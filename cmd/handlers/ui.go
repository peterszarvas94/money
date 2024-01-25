package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"html"
	"net/http"
	"pengoe/internal/router"
	"pengoe/internal/services"
	c "pengoe/web/templates/components"
	"strconv"
	"time"

	"github.com/a-h/templ"
)

/*
/ui/new-event-form
*/
func NewEventForm(w http.ResponseWriter, r *http.Request, p map[string]string) error {
	db, found := r.Context().Value("db").(*sql.DB)
	if !found {
		router.InternalError(w, r, p)
		return errors.New("Should use db middleware")
	}

	err := r.ParseForm()
	if err != nil {
		router.BadRequest(w, r, p)
		return err
	}

	form := r.Form

	accountIdStr := html.EscapeString(form.Get("account_id"))
	accountId, err := strconv.Atoi(accountIdStr)
	if err != nil {
		router.BadRequest(w, r, p)
		return err
	}

	accountService := services.NewAccountService(db)
	account, err := accountService.GetById(accountId)
	if err != nil {
		router.InternalError(w, r, p)
		return err
	}

	eventFormData := c.EventFormProps{
		New:         true,
		Currency:    account.Currency,
		DeliveredAt: time.Now().UTC(),
		HxTarget:    "closest li",
	}

	eventForm := c.EventForm(eventFormData)

	popupData := c.PopupProps{
		CloseUrl: "/ui/new-event-form-button",
		Child:    eventForm,
	}

	popup := c.Popup(popupData)
	handler := templ.Handler(popup)
	handler.ServeHTTP(w, r)

	return nil
}

/*
/ui/new-event-form-button
*/
func NewEventFormButton(w http.ResponseWriter, r *http.Request, p map[string]string) error {
	element := c.NewEventFormButton()
	handler := templ.Handler(element)
	handler.ServeHTTP(w, r)

	return nil
}

/*
/ui/edit-event-form
*/
func EditEventForm(w http.ResponseWriter, r *http.Request, p map[string]string) error {
	db, found := r.Context().Value("db").(*sql.DB)
	if !found {
		router.InternalError(w, r, p)
		return errors.New("Should use db middleware")
	}

	idStr := p["id"]
	eventId, err := strconv.Atoi(idStr)
	if err != nil {
		router.BadRequest(w, r, p)
		return err
	}

	eventService := services.NewEventService(db)
	accountService := services.NewAccountService(db)

	event, err := eventService.GetById(eventId)
	if err != nil {
		router.InternalError(w, r, p)
		return err
	}

	account, err := accountService.GetById(event.AccountId)
	if err != nil {
		router.InternalError(w, r, p)
		return err
	}

	data := c.EventFormProps{
		Currency:    account.Currency,
		EventId:     event.Id,
		Name:        event.Name,
		Description: event.Description,
		Income:      event.Income,
		Reserved:    event.Reserved,
		DeliveredAt: event.DeliveredAt,
		HxTarget:    "closest div",
	}

	eventForm := c.EventForm(data)

	popupData := c.PopupProps{
		CloseUrl: fmt.Sprintf(
			"/ui/event-card/%d",
			eventId,
		),
		Child: eventForm,
	}

	popup := c.Popup(popupData)
	handler := templ.Handler(popup)
	handler.ServeHTTP(w, r)

	return nil
}

/*
/event-card
*/
func EventCard(w http.ResponseWriter, r *http.Request, p map[string]string) error {
	db, found := r.Context().Value("db").(*sql.DB)
	if !found {
		router.InternalError(w, r, p)
		return errors.New("Should use db middleware")
	}

	idStr := p["id"]
	eventId, err := strconv.Atoi(idStr)
	if err != nil {
		router.BadRequest(w, r, p)
		return err
	}

	eventService := services.NewEventService(db)
	accountService := services.NewAccountService(db)

	event, err := eventService.GetById(eventId)
	if err != nil {
		router.InternalError(w, r, p)
		return err
	}

	account, err := accountService.GetById(event.AccountId)
	if err != nil {
		router.InternalError(w, r, p)
		return err
	}

	data := c.EventCardProps{
		Currency:    account.Currency,
		EventId:     event.Id,
		Name:        event.Name,
		Description: event.Description,
		Income:      event.Income,
		Reserved:    event.Reserved,
		DeliveredAt: event.DeliveredAt,
	}

	eventCard := c.EventCard(data)
	handler := templ.Handler(eventCard)
	handler.ServeHTTP(w, r)

	return nil
}
