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
	flag.StringVar(&logger.Loglevel, "log", "INFO", "-log INFO|WARNING|ERROR|FATAL")
	flag.Parse()

	r := router.NewRouter()

	// home page
	r.GET("/", handlers.HomePageHandler)

	// test page
	// r.GET("/test/:id/deez/:a", handlers.TestHandler)
	// r.GET("/test/:id/deez/nuts", handlers.TestHandler)

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

	// refresh token
	r.POST("/refresh", handlers.RefreshTokenHandler)

	// dashboard
	r.GET("/dashboard", handlers.DashboardPageHandler)

	// new account page
	r.GET("/account/new", handlers.NewAccountPageHandler)
	r.POST("/account", handlers.NewAccountHandler)

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
