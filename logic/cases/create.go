package cases

import (
	"covid19kalteng/models"

	. "covid19kalteng/modules"
	"covid19kalteng/modules/date"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

//New covid case
func New(c echo.Context, data interface{}) (interface{}, error) {

	cases := models.Case{}
	casePayload := CasePayload{}

	payloadRules := govalidator.MapData{
		"date":        []string{},
		"location":    []string{},
		"data_detail": []string{},
	}

	//validate data from body payload
	validate := ValidateRequestPayload(c, payloadRules, &casePayload)
	if validate != nil {
		return nil, ReturnInvalidResponse(http.StatusUnprocessableEntity, validate, "Kesalahan validasi")
	}

	//copy data from payload
	marshal, _ := json.Marshal(casePayload)
	json.Unmarshal(marshal, &cases)

	//parse date and update created_at & updated_at
	date, err := date.ParseSimple(casePayload.CreatedAt)
	if err != nil && err.Error() != "nil" {
		return nil, ReturnInvalidResponse(http.StatusUnprocessableEntity, err, "format end_date salah")

	}
	cases.CreatedAt = date
	cases.UpdatedAt = date

	//process jsonb data
	ProcessCase(&cases)

	//store new row
	err = cases.Create()

	//BUG: poor performance because kafka use new connection ? performance drop to 500ms (latency)
	//log change
	//nlogs.NAudittrail(models.Case{}, cases, c.Get("user").(*jwt.Token), "cases", fmt.Sprint(cases.ID), "create")

	return "ok", c.JSON(http.StatusCreated, cases)
}
