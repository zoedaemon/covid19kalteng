package utils

import (
	"covid19kalteng/covid19"

	"github.com/labstack/echo"
)

//GetApp just type casting
func GetApp(c echo.Context) *covid19.Application {
	return c.Get("app").(*covid19.Application)
}
