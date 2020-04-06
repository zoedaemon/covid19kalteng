package handlers

import (
	"covid19kalteng/covid19"
	"covid19kalteng/models"
	"covid19kalteng/modules"
	"covid19kalteng/modules/nlogs"
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

// CaseList get all bank list
func CaseList(c echo.Context) error {
	const LogTag = "CaseList"
	defer c.Request().Body.Close()
	//get token
	token := c.Get("user").(*jwt.Token)

	// filters
	const locProvMainField = "provinsi_main"
	const locProvField = "provinsi"
	const locKotKabField = "kota_kabupaten"
	locProvMainData := c.QueryParam(locProvMainField)
	locProvMain := c.QueryParam(locProvField)
	locKotKabData := c.QueryParam(locKotKabField)

	//init models for response
	var cases []models.Case

	//Pagination Custom Query
	QPaged := modules.QueryPaged{}
	QPaged.Init(c)
	//custom query
	db := covid19.App.DB
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

	//generate filter, return db and error
	db, err = QPaged.GenerateFilters(db, CaseFilter{}, "cases")
	if err != nil {
		nlogs.NLog("warning", LogTag, map[string]interface{}{
			"message": "error list cases",
			"error":   err}, token, "", false)

		return returnInvalidResponse(http.StatusInternalServerError, err, "kesalahan dalam mendapatkan data")
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
		return returnInvalidResponse(http.StatusNoContent, err, "Data Kasus Kosong")
	}

	//get result format
	result := QPaged.GetPage(cases)
	return c.JSON(http.StatusOK, result)
}
