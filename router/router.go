package router

import (
	"covid19kalteng/handlers"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// NewRouter func
func NewRouter() *echo.Echo {
	e := echo.New()

	// ignore /api-lender
	e.Pre(middleware.Rewrite(map[string]string{
		"/api-lender/*":       "/$1",
		"/api-lender-devel/*": "/$1",
	}))

	e.GET("/clientauth", handlers.ClientLogin)

	AdminGroup(e)
	ClientGroup(e)
	ReporterGroup(e)

	return e
}
