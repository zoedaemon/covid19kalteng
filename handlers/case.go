package handlers

import (
	"covid19kalteng/covid19"
	"covid19kalteng/models"
	"covid19kalteng/modules/nlogs"
	"net/http"
	"strconv"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// CaseList get all bank list
func CaseList(c echo.Context) error {

	const LogTag = "CaseList"

	type Filter struct {
		Location string `json:"title" condition:"LIKE"`
	}

	defer c.Request().Body.Close()

	// pagination parameters
	rows, err := strconv.Atoi(c.QueryParam("rows"))
	page, err := strconv.Atoi(c.QueryParam("page"))
	orderby := strings.Split(c.QueryParam("orderby"), ",")
	sort := strings.Split(c.QueryParam("sort"), ",")

	// filters
	loc := c.QueryParam("location")

	//find rows by filter
	cases := models.Case{}
	result, err := cases.FindPaged(page, rows, orderby, sort, &Filter{
		Location: loc,
	})
	if err != nil {
		nlogs.NLog("warning", LogTag, map[string]interface{}{
			NLOGMSG:   "error get Cases list",
			NLOGERR:   err,
			NLOGQUERY: covid19.App.DB.QueryExpr()}, c.Get("user").(*jwt.Token), "", true)

		return returnInvalidResponse(http.StatusNoContent, err, "data kasus kosong")
	}

	return c.JSON(http.StatusOK, result)
}
