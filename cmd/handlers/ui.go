package handlers

import (
	"errors"
	"fmt"
	"html"
	"net/http"
	"pengoe/internal/router"
	t "pengoe/internal/token"
	"pengoe/web/templates/components"
	"strconv"

	"github.com/a-h/templ"
)

func NewEventFormHandler(w http.ResponseWriter, r *http.Request, p map[string]string) error {
	token, found := r.Context().Value("token").(*t.Token)
	if !found {
		router.InternalError(w, r, p)
		return errors.New("Should use token middleware")
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

	data := components.EventFormProps{
		EmptyFormAccountId: accountId,
		Token:              token,
		Currency:           "USD",
	}

	eventForm := components.EventForm(data)

	popupData := components.PopupProps{
		CloseUrl: fmt.Sprintf(
			"/ui/new-event-form-button?account_id=%s",
			accountIdStr,
		),
		Child: eventForm,
	}

	popup := components.Popup(popupData)
	handler := templ.Handler(popup)
	handler.ServeHTTP(w, r)

	return nil
}

func NewEventFormButtonHandler(w http.ResponseWriter, r *http.Request, p map[string]string) error {
	err := r.ParseForm()
	if err != nil {
		router.BadRequest(w, r, p)
		return err
	}

	form := r.Form

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

	element := components.NewEventFormButton(components.NewEventFormButtonProps{
		AccountId: accountId,
	})

	handler := templ.Handler(element)
	handler.ServeHTTP(w, r)

	return nil
}
