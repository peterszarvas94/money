package main

import (
	"flag"
	"net/http"
	"os"
	h "pengoe/cmd/handlers"
	m "pengoe/cmd/middlewares"
	"pengoe/internal/logger"
	"pengoe/internal/router"
)

func main() {
	// get log level from command line flag -log
	flag.StringVar(&logger.Loglevel, "log", "INFO", "-log INFO|WARNING|ERROR|FATAL")
	flag.Parse()

	// create router
	r := router.NewRouter()

	// home page
	r.GET("/", h.HomePageHandler)

	// signup
	r.GET("/signup", h.SignupPage, m.AuthPage)
	r.POST("/signup", h.Signup, m.AuthPage, m.DB)

	// signin
	r.GET("/signin", h.SigninPage, m.AuthPage)
	r.POST("/signin", h.Signin, m.AuthPage, m.DB)

	// signout
	r.POST("/signout", h.Signout, m.Token, m.DB, m.Session)

	// dashboard
	r.GET("/dashboard", h.DashboardPage, m.Token, m.DB, m.Session)

	// account
	r.GET("/account/new", h.NewAccountPage, m.Token, m.DB, m.Session)
	r.POST("/account", h.NewAccount, m.Token, m.DB, m.Session)
	r.GET("/account/:id", h.AccountPage, m.Token, m.DB, m.Session)
	r.DELETE("/account/:id", h.DeleteAccount, m.Token, m.DB, m.Session)

	// event
	r.POST("/event", h.NewEvent, m.Token, m.DB, m.Session)
	r.PATCH("/event/:id", h.EditEvent, m.Token, m.DB, m.Session)
	r.DELETE("/event/:id", h.DeleteEvent, m.Token, m.DB, m.Session)

	// ui
	r.GET("/ui/check", h.CheckUser, m.DB)
	r.GET("/ui/new-event-form", h.NewEventForm, m.DB)
	r.GET("/ui/new-event-form-button", h.NewEventFormButton)
	r.GET("/ui/edit-event-form/:id", h.EditEventForm, m.DB)
	r.GET("/ui/event-card/:id", h.EventCard, m.DB)

	// static files
	r.SetStaticPath("/static", "./web/static")

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		logger.Log(logger.FATAL, "main/listen", err.Error())
		os.Exit(1)
	}
}
