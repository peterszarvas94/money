package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"pengoe/internal/router"
	"pengoe/internal/services"
	t "pengoe/internal/token"
	"pengoe/web/templates/components"
	c "pengoe/web/templates/components"
	"strconv"
	"time"

	"github.com/a-h/templ"
)

func NewEvent(w http.ResponseWriter, r *http.Request, p map[string]string) error {
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

	eventService := services.NewEventService(db)
	accessService := services.NewAccessService(db)
	accountService := services.NewAccountService(db)

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
		newToken, err := t.Manager.RenewToken(session.Id)
		if err != nil {
			return router.InternalError(w, r, p)
		}

		w.Header().Set("HX-Retarget", "#csrf")
		w.Header().Set("HX-Trigger", "new-event")

		data := c.CsrfProps{
			Token: newToken,
		}

		component := c.Csrf(data)
		handler := templ.Handler(component)
		handler.ServeHTTP(w, r)

		return nil
	}

	// csrf token is not expired

	newEvent := &services.Event{
		AccountId:   accountId,
		Name:        name,
		Description: description,
		Income:      incomeInt,
		Reserved:    reservedInt,
		DeliveredAt: deliveredAt,
	}

	event, newEventErr := eventService.New(newEvent)
	if newEventErr != nil {
		router.InternalError(w, r, p)
		return newEventErr
	}

	account, err := accountService.GetById(accountId)
	if err != nil {
		router.InternalError(w, r, p)
		return err
	}

	component := c.NewEventCard(c.NewEventCardProps{
		EventCardProps: c.EventCardProps{
			Currency:    account.Currency,
			EventId:     event.Id,
			Name:        event.Name,
			Description: event.Description,
			Income:      event.Income,
			Reserved:    event.Reserved,
			DeliveredAt: event.DeliveredAt,
		},
	})
	handler := templ.Handler(component)
	handler.ServeHTTP(w, r)

	return nil
}

func EditEvent(w http.ResponseWriter, r *http.Request, p map[string]string) error {
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

	id := p["id"]
	eventId, err := strconv.Atoi(id)
	if err != nil {
		router.BadRequest(w, r, p)
		return err
	}

	err = r.ParseForm()
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

	accountService := services.NewAccountService(db)
	account, err := accountService.GetById(accountId)
	if err != nil {
		router.InternalError(w, r, p)
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
		Id:          eventId,
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
		newToken, err := t.Manager.RenewToken(session.Id)
		if err != nil {
			return router.InternalError(w, r, p)
		}

		w.Header().Set("HX-Retarget", "#csrf")
		w.Header().Set("HX-Trigger", fmt.Sprintf("edit-event-%d", eventId))

		data := c.CsrfProps{
			Token: newToken,
		}

		component := c.Csrf(data)
		handler := templ.Handler(component)
		handler.ServeHTTP(w, r)
	}

	// csrf token is not expired

	_, err = eventService.Update(newEvent)
	if err != nil {
		router.InternalError(w, r, p)
		return err
	}

	data := c.EventCardProps{
		Currency:    account.Currency,
		EventId:     newEvent.Id,
		Name:        newEvent.Name,
		Description: newEvent.Description,
		Income:      newEvent.Income,
		Reserved:    newEvent.Reserved,
		DeliveredAt: newEvent.DeliveredAt,
	}

	component := c.EventCard(data)
	handler := templ.Handler(component)
	handler.ServeHTTP(w, r)

	return nil
}

func DeleteEvent(w http.ResponseWriter, r *http.Request, p map[string]string) error {
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

	id := p["id"]
	eventId, err := strconv.Atoi(id)
	if err != nil {
		router.BadRequest(w, r, p)
		return err
	}

	err = r.ParseForm()
	if err != nil {
		router.InternalError(w, r, p)
		return err
	}

	// manually parse body, (because DELETE request)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		router.InternalError(w, r, p)
		return err
	}

	formValues, err := url.ParseQuery(string(body))
	if err != nil {
		router.InternalError(w, r, p)
		return err
	}

	formToken := html.EscapeString(formValues.Get("csrf"))

	accountIdStr := html.EscapeString(formValues.Get("account_id"))
	if accountIdStr == "" {
		router.BadRequest(w, r, p)
		return errors.New("Account id is required")
	}

	accountId, err := strconv.Atoi(accountIdStr)
	if err != nil {
		router.BadRequest(w, r, p)
		return err
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
		router.Unauthorized(w, r, p)
		return errors.New("Tokens do not match")
	}

	// token is expired
	if token.Valid.Before(time.Now().UTC()) {
		newToken, err := t.Manager.RenewToken(session.Id)
		if err != nil {
			router.InternalError(w, r, p)
			return err
		}

		data := components.CsrfProps{
			Token: newToken,
		}

		w.Header().Set("HX-Retarget", fmt.Sprintf("#csrf"))
		w.Header().Set("HX-Trigger", fmt.Sprintf("delete-event-%d", eventId))

		component := components.Csrf(data)
		handler := templ.Handler(component)
		handler.ServeHTTP(w, r)

		return nil
	}

	// csrf token is not expired
	err = eventService.Delete(eventId)
	if err != nil {
		router.InternalError(w, r, p)
		return err
	}

	// no return because delete

	return nil
}
