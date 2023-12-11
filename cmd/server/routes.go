package main

import (
	"bytes"
	"context"
	"net/http"
	"os"
	"time"

	"github.com/Francesco99975/shorehamex/internal/controllers"
	"github.com/Francesco99975/shorehamex/internal/middlewares"
	"github.com/Francesco99975/shorehamex/internal/models"
	"github.com/Francesco99975/shorehamex/views"

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

	e.GET("/patient", controllers.Patient(), middlewares.ExaminationReverseAuthMiddleware())

	e.POST("/patient", controllers.PatientLogin())

	e.GET("/success", controllers.Success())

	patientGroup := e.Group("/examination", middlewares.ExaminationAuthMiddleware())

	patientGroup.GET("", controllers.Examination())

	patientGroup.GET("/asq", controllers.Asq(false))

	patientGroup.POST("/asq", controllers.AsqCalc(false))

	patientGroup.GET("/bai", controllers.Bai(false))

	patientGroup.POST("/bai", controllers.BaiCalc(false))

	patientGroup.GET("/bdi", controllers.Bdi(false))

	patientGroup.POST("/bdi", controllers.BdiCalc(false))

	patientGroup.GET("/p3", controllers.P3(false))

	patientGroup.POST("/p3", controllers.P3Calc(false))

	patientGroup.GET("/mmpi", controllers.MMPI(false))

	patientGroup.POST("/mmpi", controllers.MMPICalc(false))

	adminGroup := e.Group("/admin", middlewares.AuthMiddleware())

	adminGroup.GET("", controllers.Admin())

	adminGroup.POST("", controllers.GenerateCodes())

	adminGroup.GET("/asq", controllers.Asq(true))

	adminGroup.POST("/asq", controllers.AsqCalc(true))

	adminGroup.GET("/bai", controllers.Bai(true))

	adminGroup.POST("/bai", controllers.BaiCalc(true))

	adminGroup.GET("/bdi", controllers.Bdi(true))

	adminGroup.POST("/bdi", controllers.BdiCalc(true))

	adminGroup.GET("/p3", controllers.P3(true))

	adminGroup.POST("/p3", controllers.P3Calc(true))

	adminGroup.GET("/mmpi", controllers.MMPI(true))

	adminGroup.POST("/mmpi", controllers.MMPICalc(true))

	e.POST("/login", controllers.Login())

	e.POST("/logout", controllers.Logout())

	e.HTTPErrorHandler = serverErrorHandler

	return e
}

func serverErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	data := models.Site{
		AppName:  "Shoreham Examination",
		Title:    "Client Error",
		Metatags: models.SEO{Description: "Examination tool", Keywords: "tools,exam"},
		Year:     time.Now().Year(),
	}

	buf := bytes.NewBuffer(nil)
	if code < 500 {
		views.ClientError(data, err).Render(context.Background(), buf)

	} else {
		views.ServerError(data, err).Render(context.Background(), buf)
	}

	c.Blob(200, "text/html; charset=utf-8", buf.Bytes())

}
