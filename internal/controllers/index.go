package controllers

import (
	"bytes"
	"context"

	"github.com/Francesco99975/shorehamex/views"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func Index() echo.HandlerFunc {
	buf := bytes.NewBuffer(nil)

	// data := models.Site{
	// 	AppName:  "HTMX + GO",
	// 	Title:    "Home",
	// 	CSRF:     csrf.Token(r),
	// 	Metatags: models.SEO{Description: "Basic boilerplate for go web apps", Keywords: "go,htmx,web"},
	// 	Year:     time.Now().Year(),
	// }

	err := views.HomePage("Demo - Home").Render(context.Background(), buf)

	if err != nil {
		log.Warn("TODO: you need to implement this properly")
		log.Errorf("rendering index: %s", err)
	}

	return func(c echo.Context) error {
		return c.Blob(200, "text/html; charset=utf-8", buf.Bytes())
	}
}
