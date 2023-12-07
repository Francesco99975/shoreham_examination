package middlewares

import (
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, err := session.Get("session", c)
			if err != nil {
				return c.Redirect(http.StatusSeeOther, "/")
			}

			if auth, ok := sess.Values["authenticated"].(bool); !ok || !auth {
				return c.Redirect(http.StatusSeeOther, "/")
			} else {
				return next(c)
			}
		}
	}
}

func ReverseAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, err := session.Get("session", c)
			if err != nil {
				return next(c)
			}

			if auth, ok := sess.Values["authenticated"].(bool); !ok || !auth {
				return next(c)
			} else {
				return c.Redirect(http.StatusSeeOther, "/admin")
			}
		}
	}
}

func ExaminationAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, err := session.Get("session", c)
			if err != nil {
				return c.Redirect(http.StatusSeeOther, "/")
			}

			if auth, ok := sess.Values["examauth"].(bool); !ok || !auth {
				return c.Redirect(http.StatusSeeOther, "/")
			} else {
				return next(c)
			}
		}
	}
}

func ExaminationReverseAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, err := session.Get("session", c)
			if err != nil {
				return next(c)
			}

			if auth, ok := sess.Values["examauth"].(bool); !ok || !auth {
				return next(c)
			} else {
				return c.Redirect(http.StatusSeeOther, "/")
			}
		}
	}
}
