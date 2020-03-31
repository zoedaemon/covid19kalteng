package middlewares

import (
	"covid19kalteng/covid19"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// SetClientJWTmiddlewares func
func SetClientJWTmiddlewares(g *echo.Group, role string) {
	jwtConfig := covid19.App.Config.GetStringMap(fmt.Sprintf("%s.jwt", covid19.App.ENV))

	g.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS512",
		SigningKey:    []byte(jwtConfig["jwt_secret"].(string)),
	}))

	switch role {
	case "client":
		g.Use(validateJWTclient)
		break
	case "reporter":
		g.Use(validateJWTreporter)
		break
	case "admin":
		g.Use(validateJWTadmin)
		break
	}
}

func validateJWTadmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user")
		token := user.(*jwt.Token)

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if claims["group"] == "admin" {
				return next(c)
			}
			return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("%s", "invalid token"))
		}

		return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("%s", "invalid token"))
	}
}

func validateJWTclient(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user")
		token := user.(*jwt.Token)

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if claims["group"] == "client" {
				return next(c)
			}
			return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("%s", "invalid token"))
		}

		return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("%s", "invalid token"))
	}
}

func validateJWTreporter(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user")
		token := user.(*jwt.Token)

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if claims["group"] == "reporter" {
				return next(c)
			}
			return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("%s", "invalid token"))
		}

		return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("%s", "invalid token"))
	}
}
