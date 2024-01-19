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
	r.GET("/signup", h.SignupPageHandler, m.WithRedirect)
	r.POST("/signup", h.SignupHandler, m.WithRedirect, m.WithDB)

	// signin
	r.GET("/signin", h.SigninPageHandler, m.WithRedirect)
	r.POST("/signin", h.SigninHandler, m.WithRedirect, m.WithDB)

	// signout
	r.POST("/signout", h.SignoutHandler, m.WithToken, m.WithDB, m.WithSession)

	// check useraname and email
	r.GET("/check", h.CheckUserHandler, m.WithDB)

	// dashboard
	r.GET("/dashboard", h.DashboardPageHandler, m.WithToken, m.WithDB, m.WithSession)

	// new account page
	r.GET("/account/new", h.NewAccountPageHandler, m.WithToken, m.WithDB, m.WithSession)
	r.POST("/account", h.NewAccountHandler, m.WithToken, m.WithDB, m.WithSession)
	r.DELETE("/account/:id", h.DeleteAccountHandler, m.WithToken, m.WithDB, m.WithSession)

	// account page
	r.GET("/account/:id", h.AccountPageHandler, m.WithToken, m.WithDB, m.WithSession)

	// static files
	r.SetStaticPath("/static", "./web/static")

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		logger.Log(logger.FATAL, "main/listen", err.Error())
		os.Exit(1)
	}
}
