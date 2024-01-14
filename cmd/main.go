package main

import (
	"flag"
	"net/http"
	"os"
	"pengoe/cmd/handlers"
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
	r.GET("/", handlers.HomePageHandler)

	// test page
	// r.GET("/test/:id", handlers.TestPageHandler)
	// r.GET("/test/help", handlers.TestPageHandler)
	// r.GET("/test/:id/nuts", handlers.TestPageHandler)
	// r.GET("/test/deez/nuts", handlers.TestPageHandler)

	// signup
	r.GET("/signup", handlers.SignupPageHandler)
	r.POST("/signup", handlers.NewUserHandler)

	// signin
	r.GET("/signin", handlers.SigninPageHandler)
	r.POST("/signin", handlers.SigninHandler)

	// signout
	r.POST("/signout", handlers.SignoutHandler)

	// check useraname and email
	r.GET("/check", handlers.CheckUserHandler)

	// dashboard
	r.GET("/dashboard", handlers.DashboardPageHandler)

	// new account page
	r.GET("/account/new", handlers.NewAccountPageHandler)
	r.POST("/account", handlers.NewAccountHandler)
	r.DELETE("/account/:id", handlers.DeleteAccountHandler)

	// account page
	r.GET("/account/:id", handlers.AccountPageHandler)

	// static files
	r.SetStaticPath("/static", "./web/static")

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		logger.Log(logger.FATAL, "main/listen", err.Error())
		os.Exit(1)
	}
}
