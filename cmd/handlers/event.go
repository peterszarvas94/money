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
	"pengoe/internal/utils"
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
	if formToken == "" {
		router.BadRequest(w, r, p)
		return errors.New("CSRF token is required")
	}

	accountId := html.EscapeString(form.Get("account_id"))
	if accountId == "" {
		router.BadRequest(w, r, p)
		return errors.New("Account ID is required")
	}

	name := html.EscapeString(form.Get("name"))
	if name == "" {
		router.BadRequest(w, r, p)
		return errors.New("Name is required")
	}

	description := html.EscapeString(form.Get("description"))

	incomeStr := html.EscapeString(form.Get("income"))
	if incomeStr == "" {
		router.BadRequest(w, r, p)
		return errors.New("Income is required")
	}

	income, err := strconv.Atoi(incomeStr)
	if err != nil {
		router.BadRequest(w, r, p)
		return err
	}

	reservedStr := html.EscapeString(form.Get("reserved"))
	if reservedStr == "" {
		router.BadRequest(w, r, p)
		return errors.New("Reserved is required")
	}

	reserved, err := strconv.Atoi(reservedStr)
	if err != nil {
		router.BadRequest(w, r, p)
		return err
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

	eventService := services.NewEventService(db)
	accessService := services.NewAccessService(db)
	accountService := services.NewAccountService(db)

	// check if user has access to account
	ok := accessService.Check(session.UserId, accountId)
	if !ok {
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
	id := utils.NewUUID("evt")

	err = eventService.New(id, name, description, income, reserved, deliveredAt, accountId)
	if err != nil {
		router.InternalError(w, r, p)
		return err
	}

	account, err := accountService.GetById(accountId)
	if err != nil {
		router.InternalError(w, r, p)
		return err
	}

	component := c.NewEventCard(c.NewEventCardProps{
		EventCardProps: c.EventCardProps{
			Currency:    account.Currency,
			EventId:     id,
			Name:        name,
			Description: description,
			Income:      income,
			Reserved:    reserved,
			DeliveredAt: deliveredAt,
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

	eventId, found := p["id"]
	if !found {
		router.NotFound(w, r, p)
		return errors.New("Path variable \"id\" not found")
	}

	err := r.ParseForm()
	if err != nil {
		router.InternalError(w, r, p)
		return err
	}

	form := r.Form
	formToken := html.EscapeString(form.Get("csrf"))
	if formToken == "" {
		router.BadRequest(w, r, p)
		return errors.New("CSRF token is required")
	}

	accountId := html.EscapeString(form.Get("account_id"))
	if accountId == "" {
		router.BadRequest(w, r, p)
		return errors.New("Account id is required")
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

	incomeStr := html.EscapeString(form.Get("income"))
	if incomeStr == "" {
		router.BadRequest(w, r, p)
		return errors.New("Income is required")
	}

	income, err := strconv.Atoi(incomeStr)
	if err != nil {
		router.BadRequest(w, r, p)
		return err
	}

	reservedStr := html.EscapeString(form.Get("reserved"))
	if reservedStr == "" {
		router.BadRequest(w, r, p)
		return errors.New("Reserved is required")
	}

	reserved, err := strconv.Atoi(reservedStr)
	if err != nil {
		router.BadRequest(w, r, p)
		return err
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

	eventService := services.NewEventService(db)
	accessService := services.NewAccessService(db)

	// check if user has access to account
	ok := accessService.Check(session.UserId, accountId)
	if !ok {
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
		w.Header().Set("HX-Trigger", fmt.Sprintf("edit-event-%s", eventId))

		data := c.CsrfProps{
			Token: newToken,
		}

		component := c.Csrf(data)
		handler := templ.Handler(component)
		handler.ServeHTTP(w, r)
	}

	// csrf token is not expired

	err = eventService.Update(eventId, name, description, income, reserved, deliveredAt)
	if err != nil {
		router.InternalError(w, r, p)
		return err
	}

	data := c.EventCardProps{
		Currency:    account.Currency,
		EventId:     eventId,
		Name:        name,
		Description: description,
		Income:      income,
		Reserved:    reserved,
		DeliveredAt: deliveredAt,
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

	eventId, found := p["id"]
	if !found {
		router.NotFound(w, r, p)
		return errors.New("Path variable \"id\" not found")
	}

	err := r.ParseForm()
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

	accountId := html.EscapeString(formValues.Get("account_id"))
	if accountId == "" {
		router.BadRequest(w, r, p)
		return errors.New("Account id is required")
	}

	eventService := services.NewEventService(db)
	accessService := services.NewAccessService(db)

	// check if user has access to account
	ok := accessService.Check(session.UserId, accountId)
	if !ok {
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
		w.Header().Set("HX-Trigger", fmt.Sprintf("delete-event-%s", eventId))

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
