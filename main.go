package main

import (
	"pengoe/handlers"
	"pengoe/utils"
	"net/http"
	"os"
)

func main() {
	r := utils.NewRouter()

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

	r.SetStaticPath("/static")

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		utils.Log(utils.FATAL, "main/listen", err.Error())
		os.Exit(1)
	}
}
