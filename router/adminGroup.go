package router

import (
	"covid19kalteng/handlers"
	"covid19kalteng/handlersAdmin"
	"covid19kalteng/middlewares"

	"github.com/labstack/echo"
)

// AdminGroup func
func AdminGroup(e *echo.Echo) {
	g := e.Group("/admin")
	middlewares.SetClientJWTmiddlewares(g, "admin")

	// config info
	g.GET("/info", handlers.AppInfo)
	g.GET("/profile", handlersAdmin.AdminProfile)

	// Client Management
	g.POST("/client", handlersAdmin.CreateClient)

	// Role
	g.GET("/roles", handlersAdmin.RoleList)
	g.GET("/roles/:id", handlersAdmin.RoleDetails)
	g.POST("/roles", handlersAdmin.RoleNew)
	g.PATCH("/roles/:id", handlersAdmin.RolePatch)
	g.GET("/roles_all", handlersAdmin.RoleRange)

	// Permission
	g.GET("/permission", handlersAdmin.PermissionList)

	// User
	g.GET("/users", handlersAdmin.UserList)
	g.GET("/users/:id", handlersAdmin.UserDetails)
	g.POST("/users", handlersAdmin.UserNew)
	g.PATCH("/users/:id", handlersAdmin.UserPatch)

}
