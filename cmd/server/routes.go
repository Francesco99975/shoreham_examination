package main

import (
	"github.com/Francesco99975/shorehamex/internal/controllers"

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
	// e.Static("/assets", "")

	e.GET("/", controllers.Index())

	e.GET("/staff", controllers.Staff())

	e.GET("/patient", controllers.Patient())

	e.GET("/admin", controllers.Admin())

	e.GET("/asq", controllers.Asq())

	return e
}
