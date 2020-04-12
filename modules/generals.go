package modules

import (
	"covid19kalteng/covid19"
	"covid19kalteng/modules/nlogs"
	"encoding/json"
	"math/rand"

	"fmt"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

const (
	//NLOGMSG for message body
	NLOGMSG = "message"
	//NLOGERR for error info
	NLOGERR = "error"
	//NLOGQUERY for detailed query tracing
	NLOGQUERY = "query"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

type (
	// JWTclaims jwtclaims
	JWTclaims struct {
		Username    string   `json:"username"`
		Group       string   `json:"group"`
		Permissions []string `json:"permissions"`
		jwt.StandardClaims
	}
)

// general function to validate all kind of api request payload / body
func ValidateRequestPayload(c echo.Context, rules govalidator.MapData, data interface{}) (i interface{}) {
	opts := govalidator.Options{
		Request: c.Request(),
		Data:    data,
		Rules:   rules,
	}

	v := govalidator.New(opts)

	mappedError := v.ValidateJSON()

	if len(mappedError) > 0 {
		i = mappedError
	}

	return i
}

// general function to validate all kind of api request url query
func ValidateRequestQuery(c echo.Context, rules govalidator.MapData) (i interface{}) {
	opts := govalidator.Options{
		Request: c.Request(),
		Rules:   rules,
	}

	v := govalidator.New(opts)

	mappedError := v.Validate()

	if len(mappedError) > 0 {
		i = mappedError
	}

	return i
}

func ReturnInvalidResponse(httpcode int, details interface{}, message string) error {
	responseBody := map[string]interface{}{
		"message": message,
		"details": details,
	}
	return echo.NewHTTPError(httpcode, responseBody)
}

// self explanation
func createJwtToken(id string, group string) (string, error) {
	jwtConf := covid19.App.Config.GetStringMap(fmt.Sprintf("%s.jwt", covid19.App.ENV))

	type PermModel struct {
		Permission string `json:"permissions" gorm:"column:permissions"`
	}
	var permissions []string
	var permModel []PermModel
	var db = covid19.App.DB
	switch group {
	case "users":
		err := db.Table("roles").
			Select("DISTINCT TRIM(UNNEST(roles.permissions)) as permissions").
			Joins("INNER JOIN users u ON roles.id IN (SELECT UNNEST(u.roles))").
			Where("roles.status = ?", "active").
			Where("u.id = ?", id).Scan(&permModel).Error
		if err != nil {
			return "", err
		}
		for _, v := range permModel {
			permissions = append(permissions, v.Permission)
		}
		break
	}

	claim := JWTclaims{
		id,
		group,
		permissions,
		jwt.StandardClaims{
			Id:        id,
			ExpiresAt: time.Now().Add(time.Duration(jwtConf["duration"].(int)) * time.Minute).Unix(),
		},
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claim)
	token, err := rawToken.SignedString([]byte(jwtConf["jwt_secret"].(string)))
	if err != nil {
		return "", err
	}

	return token, nil
}

// RandString random string alphanumeric. parameter length
func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func customSplit(str string, separator string) []string {
	split := strings.Split(str, separator)
	if len(split) == 1 {
		if split[0] == "" {

			split = []string{}
		}
	}

	return split
}

func validatePermission(c echo.Context, permission string) error {
	user := c.Get("user")
	token := user.(*jwt.Token)

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if claimPermissions, ok := claims["permissions"]; ok {
			s := strings.Split(strings.Trim(fmt.Sprintf("%v", claimPermissions), "[]"), " ")
			for _, v := range s {
				if strings.ToLower(v) == strings.ToLower(permission) || strings.ToLower(v) == "all" {
					return nil
				}
			}
		}
		nlogs.NLog("warning", "validatePermission", map[string]interface{}{"message": fmt.Sprintf("user dont have permission %v", permission)}, user.(*jwt.Token), "", false)

		return fmt.Errorf("Permission Denied")
	}
	nlogs.NLog("warning", "validatePermission", map[string]interface{}{"message": fmt.Sprintf("user dont have permission %v", permission)}, user.(*jwt.Token), "", false)

	return fmt.Errorf("Permission Denied")
}

//ParseError custom errors detail
type ParseError struct {
	Info interface{}
}

//Error must implement this Error func
func (e *ParseError) Error() string {
	dat, _ := json.Marshal(e.Info)
	return string(dat)
}
