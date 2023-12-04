package main

import (
	"os"

	"github.com/Francesco99975/shorehamex/internal/controllers"
	"github.com/Francesco99975/shorehamex/internal/middlewares"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func createRouter() *echo.Echo {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))))

	e.Static("/assets", "./static")

	e.GET("/", controllers.Index(), middlewares.ReverseAuthMiddleware())

	e.GET("/staff", controllers.Staff(), middlewares.ReverseAuthMiddleware())

	e.GET("/patient", controllers.Patient())

	adminGroup := e.Group("/admin", middlewares.AuthMiddleware())

	adminGroup.GET("", controllers.Admin())

	adminGroup.GET("/asq", controllers.Asq(true))

	adminGroup.GET("/bai", controllers.Bai(true))

	adminGroup.GET("/bdi", controllers.Bdi(true))

	adminGroup.GET("/p3", controllers.P3(true))

	adminGroup.GET("/mmpi", controllers.MMPI(true))

	adminGroup.POST("/mmpi", controllers.MMPICalc())

	e.POST("/login", controllers.Login())

	e.POST("/logout", controllers.Logout())

	return e
}
