package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"html"
	"net/http"
	"pengoe/internal/router"
	"pengoe/internal/services"
	t "pengoe/internal/token"
	"pengoe/web/templates/components"
	"strconv"
	"time"

	"github.com/a-h/templ"
)

func NewEventHandler(w http.ResponseWriter, r *http.Request, p map[string]string) error {
	token, found := r.Context().Value("token").(*t.Token)
	if !found {
		router.InternalError(w, r, p)
		return errors.New("Should use token middleware")
	}
	db, found := r.Context().Value("db").(*sql.DB)
	if !found {
		router.InternalError(w, r, p)
		return errors.New("Should use db middleware")
	}
	session, found := r.Context().Value("session").(*services.Session)
	if !found {
		router.InternalError(w, r, p)
		return errors.New("Should use session middleware")
	}

	err := r.ParseForm()
	if err != nil {
		router.InternalError(w, r, p)
		return err
	}

	form := r.Form

	formToken := html.EscapeString(form.Get("csrf"))

	accountIdStr := html.EscapeString(form.Get("account_id"))
	if accountIdStr == "" {
		router.BadRequest(w, r, p)
		return errors.New("Account id is required")
	}

	accountId, err := strconv.Atoi(accountIdStr)
	if err != nil {
		router.BadRequest(w, r, p)
		return err
	}

	name := html.EscapeString(form.Get("name"))
	if name == "" {
		router.BadRequest(w, r, p)
		return errors.New("Name is required")
	}

	description := html.EscapeString(form.Get("description"))

	income := html.EscapeString(form.Get("income"))
	if income == "" {
		router.BadRequest(w, r, p)
		return errors.New("Income is required")
	}

	reserved := html.EscapeString(form.Get("reserved"))
	if reserved == "" {
		router.BadRequest(w, r, p)
		return errors.New("Reserved is required")
	}

	deliveredAtStr := html.EscapeString(form.Get("delivered_at"))
	if deliveredAtStr == "" {
		router.BadRequest(w, r, p)
		return errors.New("Delivered at is required")
	}

	deliveredAt, err := time.Parse("2006-01-02", deliveredAtStr)
	if err != nil {
		router.BadRequest(w, r, p)
		return err
	}

	incomeInt, err := strconv.Atoi(income)
	if err != nil {
		router.BadRequest(w, r, p)
		return err
	}

	reservedInt, err := strconv.Atoi(reserved)
	if err != nil {
		router.BadRequest(w, r, p)
		return err
	}

	newEvent := &services.Event{
		AccountId:   accountId,
		Name:        name,
		Description: description,
		Income:      incomeInt,
		Reserved:    reservedInt,
		DeliveredAt: deliveredAt,
	}

	eventService := services.NewEventService(db)
	accessService := services.NewAccessService(db)

	// check if user has access to account
	err = accessService.Check(session.UserId, accountId)
	if err != nil {
		router.Unauthorized(w, r, p)
		return err
	}

	// check if the tokens match
	if token.Value != formToken {
		return router.Unauthorized(w, r, p)
	}

	// token is expired
	if token.Valid.Before(time.Now().UTC()) {
		newCsrfToken, tokenErr := t.Manager.RenewToken(session.Id)
		if tokenErr != nil {
			return router.InternalError(w, r, p)
		}

		data := components.EventFormProps{
			EmptyFormAccountId: accountId,
			Event:              newEvent,
			Token:              newCsrfToken,
			Refetch:            true,
		}

		component := components.EventForm(data)
		handler := templ.Handler(component)
		handler.ServeHTTP(w, r)

		return nil
	}

	// csrf token is not expired

	_, newEventErr := eventService.New(newEvent)
	if newEventErr != nil {
		router.InternalError(w, r, p)
		return newEventErr
	}

	w.Header().Set("HX-Redirect", fmt.Sprintf("/account/%d", accountId))
	return nil
}
