package handlers

import (
	"covid19kalteng/covid19"
	"covid19kalteng/models"
	"covid19kalteng/modules/nlogs"

	"fmt"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
	"golang.org/x/crypto/bcrypt"
)

type (
	// ReporterLoginCreds type
	ReporterLoginCreds struct {
		Key      string `json:"key"`
		Password string `json:"password"`
	}
)

// ReporterLogin lender can choose either login with email / phone
func ReporterLogin(c echo.Context) error {
	defer c.Request().Body.Close()

	var (
		credentials ReporterLoginCreds
		lender      models.User
		validKey    bool
		token       string
		err         error
	)

	rules := govalidator.MapData{
		"key":      []string{"required"},
		"password": []string{"required"},
	}

	validate := validateRequestPayload(c, rules, &credentials)
	if validate != nil {
		nlogs.NLog("warning", "ReporterLogin", map[string]interface{}{"message": "validation error", "error": validate}, c.Get("user").(*jwt.Token), "", true)

		return returnInvalidResponse(http.StatusBadRequest, validate, "Login tidak valid")
	}

	// check if theres record
	validKey = covid19.App.DB.
		Where("username = ?", credentials.Key).
		Where("status = ?", "active").
		Find(&lender).RecordNotFound()

	if !validKey { // check the password
		err = bcrypt.CompareHashAndPassword([]byte(lender.Password), []byte(credentials.Password))
		if err != nil {
			nlogs.NLog("warning", "ReporterLogin", map[string]interface{}{"message": fmt.Sprintf("password error on user %v", credentials.Key), "error": err}, c.Get("user").(*jwt.Token), "", true)

			return returnInvalidResponse(http.StatusUnauthorized, err, "Login tidak valid")
		}

		token, err = createJwtToken(strconv.FormatUint(lender.ID, 10), "reporter")
		if err != nil {
			nlogs.NLog("warning", "ReporterLogin", map[string]interface{}{"message": "error generating token", "error": err}, c.Get("user").(*jwt.Token), "", true)

			return returnInvalidResponse(http.StatusInternalServerError, err, "Terjadi kesalahan")
		}
	} else {
		nlogs.NLog("warning", "ReporterLogin", map[string]interface{}{"message": fmt.Sprintf("user not found %v", credentials.Key)}, c.Get("user").(*jwt.Token), "", true)

		return returnInvalidResponse(http.StatusUnauthorized, "username not found", "Login tidak valid")
	}

	jwtConf := covid19.App.Config.GetStringMap(fmt.Sprintf("%s.jwt", covid19.App.ENV))
	expiration := time.Duration(jwtConf["duration"].(int)) * time.Minute

	nlogs.NLog("info", "ReporterLogin", map[string]interface{}{"message": fmt.Sprintf("%v login", credentials.Key)}, c.Get("user").(*jwt.Token), "", true)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token":      token,
		"expires_in": expiration.Seconds(),
	})
}
