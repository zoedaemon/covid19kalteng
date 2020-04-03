package router

import (
	"covid19kalteng/middlewares"

	"github.com/labstack/echo"
)

// ReporterGroup group
func ReporterGroup(e *echo.Echo) {
	g := e.Group("/reporter")
	middlewares.SetClientJWTmiddlewares(g, "reporter")

	// // Profile endpoints
	// g.GET("/profile", handlers.LenderProfile)
	// g.PATCH("/profile", handlers.LenderProfileEdit)
	// g.POST("/first_login", handlers.UserFirstLoginChangePassword)

}
