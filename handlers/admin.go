package handlers

import (
	"covid19kalteng/covid19"
	"net/http"

	"github.com/labstack/echo"
)

// AppInfo show covid19 configs
func AppInfo(c echo.Context) error {
	defer c.Request().Body.Close()

	type AppInfo struct {
		AppName string                 `json:"app_name"`
		Version string                 `json:"version"`
		ENV     string                 `json:"env"`
		Config  map[string]interface{} `json:"configs"`
	}

	var show AppInfo

	show.AppName = covid19.App.Name
	show.Version = covid19.App.Version
	show.ENV = covid19.App.ENV
	show.Config = covid19.App.Config.GetStringMap(covid19.App.ENV)

	return c.JSON(http.StatusOK, show)
}
