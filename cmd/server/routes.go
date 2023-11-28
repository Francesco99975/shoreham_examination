package main

import (
	"github.com/Francesco99975/shorehamex/cmd/boot"
	"github.com/Francesco99975/shorehamex/internal/controllers"

	client "github.com/jdudmesh/gomon-client"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type Reloader struct {
}

func (r *Reloader) Reload(data string) {
	log.Debugf("Reloading trigger: %s", data)
}

func createRouter(cfg *boot.Config) *echo.Echo {

	e := echo.New()
	e.Use(middleware.Logger())
	e.Static("/assets", cfg.StaticFileDir)

	rl := &Reloader{}
	t, err := client.New(rl, e.Logger)
	if err != nil {
		log.Fatalf("unable to start reloader: %v", err)
	}
	defer t.Close()
	if err := t.Run(); err != nil {
		panic(err)
	}

	e.GET("/", controllers.Index())

	// e.GET("/about", controllers.About)

	return e
}
