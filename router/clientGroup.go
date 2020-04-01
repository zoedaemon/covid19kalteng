package router

import (
	"covid19kalteng/handlers"
	"covid19kalteng/handlersAdmin"
	"covid19kalteng/middlewares"

	"github.com/labstack/echo"
)

// ClientGroup group
func ClientGroup(e *echo.Echo) {
	g := e.Group("/client")
	middlewares.SetClientJWTmiddlewares(g, "client")
	g.POST("/reporter_login", handlers.ReporterLogin)
	g.POST("/admin_login", handlersAdmin.AdminLogin)
	g.POST("/forgotpassword", handlers.UserResetPasswordRequest)
	g.POST("/resetpassword", handlers.UserResetPasswordVerify)

	g.GET("/serviceinfo", handlers.ServiceInfo)

	//edukasi
	g.GET("/educations", handlers.EduList)
}
