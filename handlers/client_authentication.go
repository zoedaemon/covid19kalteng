package handlers

import (
	"covid19kalteng/covid19"
	"covid19kalteng/handlersAdmin"
	"covid19kalteng/models"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// ClientLogin func
func ClientLogin(c echo.Context) error {
	defer c.Request().Body.Close()

	data, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Basic "))
	if err != nil {
		return returnInvalidResponse(http.StatusUnauthorized, "", "Login tidak valid")
	}

	auth := strings.Split(string(data), ":")
	if len(auth) < 2 {
		return returnInvalidResponse(http.StatusUnauthorized, "", "Login tidak valid")
	}
	type Login struct {
		Key    string `json:"key"`
		Secret string `json:"secret"`
	}

	client := models.Client{}
	err = client.SingleFindFilter(&Login{
		Key:    auth[0],
		Secret: auth[1],
	})
	if err != nil {
		handlersAdmin.NLog("warning", "ClientLogin", map[string]interface{}{"message": "client login failed"}, c.Get("user").(*jwt.Token), "", true)

		return returnInvalidResponse(http.StatusUnauthorized, "", "Login tidak valid")
	}

	token, err := createJwtToken(strconv.FormatUint(client.ID, 10), "client")
	if err != nil {
		handlersAdmin.NLog("warning", "ClientLogin", map[string]interface{}{"message": "fail creating client token"}, c.Get("user").(*jwt.Token), "", false)

		return returnInvalidResponse(http.StatusInternalServerError, err, "Terjadi kesalahan")
	}

	jwtConf := covid19.App.Config.GetStringMap(fmt.Sprintf("%s.jwt", covid19.App.ENV))
	expiration := time.Duration(jwtConf["duration"].(int)) * time.Minute

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token":      token,
		"expires_in": expiration.Seconds(),
	})
}
