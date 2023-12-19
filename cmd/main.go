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

	r.GET("/", handlers.HomePageHandler)
	
	r.GET("/signup", handlers.SignupPageHandler)
	r.POST("/signup", handlers.NewUserHandler)

	r.GET("/signin", handlers.SigninPageHandler)
	r.POST("/signin", handlers.SigninHandler)

	r.POST("/signout", handlers.SignoutHandler)

	r.GET("/check", handlers.CheckUserHandler)

	r.POST("/refresh", handlers.RefreshTokenHandler)

	r.GET("/dashboard", handlers.DashboardPageHandler)

	r.GET("/account/new", handlers.NewAccountPageHandler)
	r.POST("/account", handlers.NewAccountHandler)

	r.SetStaticPath("/static", "./web/static")

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		logger.Log(logger.FATAL, "main/listen", err.Error())
		os.Exit(1)
	}
}
