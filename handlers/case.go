package handlers

import (
	"covid19kalteng/logic/cases"
	"covid19kalteng/models"
	. "covid19kalteng/modules"
	"covid19kalteng/modules/date"
	"covid19kalteng/modules/nlogs"
	"covid19kalteng/modules/utils"
	"fmt"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

//CaseFilter direct query %...% to data detail
type CaseFilter struct {
	ID         int64  `json:"id" condition:"optional"`
	DataDetail string `json:"data_detail" condition:"LIKE,optional"`
}

//CaseList get all bank list
func CaseList(c echo.Context) error {
	const LogTag = "CaseList"
	defer c.Request().Body.Close()

	//get token
	token := c.Get("user").(*jwt.Token)

	//filters locs
	const locProvMainField = "provinsi_main"
	const locProvField = "provinsi"
	const locKotKabField = "kota_kabupaten"
	locProvMainData := c.QueryParam(locProvMainField)
	locProvMain := c.QueryParam(locProvField)
	locKotKabData := c.QueryParam(locKotKabField)

	//filters date
	startDate, err := date.ParseSimple(c.QueryParam("start_date"))
	if err != nil && err.Error() != "nil" {
		return ReturnInvalidResponse(http.StatusUnprocessableEntity, err, "format start_date salah")
	}
	endDate, err := date.ParseSimple(c.QueryParam("end_date"))
	if err != nil && err.Error() != "nil" {
		return ReturnInvalidResponse(http.StatusUnprocessableEntity, err, "format end_date salah")
	}

	//init models for response
	var cases []models.Case

	//Pagination Custom Query
	QPaged := QueryPaged{}
	QPaged.Init(c)

	//custom query
	db := utils.GetApp(c).DB
	db = db.Table("cases").
		Select("*")
	// Joins("INNER JOIN banks b ON products.id IN (SELECT UNNEST(b.products)) ").

	//Filter locations
	LocQuery := "LOWER((location->'%s')::text) LIKE ?"
	if len(locProvMainData) != 0 {
		db = db.Where(fmt.Sprintf(LocQuery, locProvMainField), "%"+strings.ToLower(locProvMainData)+"%")
	}
	if len(locProvMain) != 0 {
		db = db.Where(fmt.Sprintf(LocQuery, locProvField), "%"+strings.ToLower(locProvMain)+"%")
	}
	if len(locKotKabData) != 0 {
		db = db.Where(fmt.Sprintf(LocQuery, locKotKabField), "%"+strings.ToLower(locKotKabData)+"%")
	}

	//default format for filter
	const beetweenCondition = "created_at BETWEEN ? AND ?"
	//filter by date
	if !startDate.IsZero() {
		if !endDate.IsZero() {
			//GOTCHAS : between startdate now() must now() + 1 day
			db = db.Where(beetweenCondition, startDate, endDate.AddDate(0, 0, 1))
		} else {
			db = db.Where(beetweenCondition, startDate, startDate.AddDate(0, 0, 1))
		}
	} else if !endDate.IsZero() {
		db = db.Where(beetweenCondition, endDate, endDate.AddDate(0, 0, 1))
	}

	//generate filter, return db and error
	db, err = QPaged.GenerateFilters(db, CaseFilter{}, "cases")
	if err != nil {
		nlogs.NLog("warning", LogTag, map[string]interface{}{
			"message": "error list cases",
			"error":   err}, token, "", false)

		return ReturnInvalidResponse(http.StatusInternalServerError, err, "kesalahan dalam mendapatkan data")
	}

	//execute anonymous function pass db and data pass by reference (services)
	err = QPaged.Exec(db, &cases, func(DB *gorm.DB, rows interface{}) error {
		//manual type casting :)
		err := DB.Find(rows.(*[]models.Case)).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		nlogs.NLog("warning", LogTag, map[string]interface{}{
			"message": "empty data cases",
			"error":   err}, token, "", false)
		return ReturnInvalidResponse(http.StatusInternalServerError, err, "Data Kasus Kosong")
	}

	//get result format
	result := QPaged.GetPage(cases)
	return c.JSON(http.StatusOK, result)
}

//CaseNew new covid case
func CaseNew(c echo.Context) error {
	defer c.Request().Body.Close()

	//TODO: fix this generals function validatePermission and others
	err := ValidatePermission(c, "cases_new")
	if err != nil {
		return ReturnInvalidResponse(http.StatusForbidden, err, err.Error())
	}
	_, info := cases.New(c, nil)
	return info
}
