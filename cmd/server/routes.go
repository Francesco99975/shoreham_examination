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

	e.GET("/success", controllers.Success())

	adminGroup := e.Group("/admin", middlewares.AuthMiddleware())

	adminGroup.GET("", controllers.Admin())

	adminGroup.GET("/asq", controllers.Asq(true))

	adminGroup.POST("/asq", controllers.AsqCalc())

	adminGroup.GET("/bai", controllers.Bai(true))

	adminGroup.POST("/bai", controllers.BaiCalc())

	adminGroup.GET("/bdi", controllers.Bdi(true))

	adminGroup.POST("/bdi", controllers.BdiCalc())

	adminGroup.GET("/p3", controllers.P3(true))

	adminGroup.POST("/p3", controllers.P3Calc())

	adminGroup.GET("/mmpi", controllers.MMPI(true))

	adminGroup.POST("/mmpi", controllers.MMPICalc())

	e.POST("/login", controllers.Login())

	e.POST("/logout", controllers.Logout())

	return e
}
