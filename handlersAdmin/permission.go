package handlersAdmin

import (
	"covid19kalteng/covid19"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

// PermissionList get all defined permissions
func PermissionList(c echo.Context) error {
	defer c.Request().Body.Close()
	err := validatePermission(c, "core_permission_list")
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	}

	permissions := covid19.App.Permission.GetStringMap(fmt.Sprintf("%s", "permissions"))

	return c.JSON(http.StatusOK, permissions)
}
