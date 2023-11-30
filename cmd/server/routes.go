package main

import (
	"os"

	"github.com/Francesco99975/shorehamex/internal/controllers"
	"github.com/Francesco99975/shorehamex/internal/middlewares"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type Reloader struct {
}

func (r *Reloader) Reload(data string) {
	log.Debugf("Reloading trigger: %s", data)
}

func createRouter() *echo.Echo {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))))

	e.Static("/assets", "./static")

	e.GET("/", controllers.Index())

	e.GET("/staff", controllers.Staff())

	e.GET("/patient", controllers.Patient())

	adminGroup := e.Group("/admin", middlewares.AuthMiddleware())

	adminGroup.GET("/", controllers.Admin())

	adminGroup.GET("/asq", controllers.Asq())

	e.POST("/login", controllers.Login())

	e.POST("/logout", controllers.Logout())

	return e
}
