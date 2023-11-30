package controllers

import (
	"bytes"
	"context"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func GeneratePage(page templ.Component) echo.HandlerFunc {
	buf := bytes.NewBuffer(nil)

	err := page.Render(context.Background(), buf)

	if err != nil {
		log.Warn("TODO: you need to implement this properly")
		log.Errorf("rendering index: %s", err)
	}

	return func(c echo.Context) error {
		return c.Blob(200, "text/html; charset=utf-8", buf.Bytes())
	}
}
