package handlers

import (
	"covid19kalteng/models"
	. "covid19kalteng/modules"
	"net/http"

	"github.com/labstack/echo"
)

// EduList get all bank list
func EduList(c echo.Context) error {

	defer c.Request().Body.Close()
	const LogTag = "EduList"

	type Filter struct {
		Title string `json:"title" condition:"LIKE"`
	}

	// pagination parameters
	edu := models.Edu{}
	SetPaginationFilter(&edu.BaseModel, c)
	// custom filters
	title := c.QueryParam("title")

	//find rows by filter
	result, err := edu.FindPaged(&Filter{
		Title: title,
	})
	if err != nil {
		// nlogs.NLog("warning", LogTag, map[string]interface{}{
		// 	NLOGMSG:   "error get Education list",
		// 	NLOGERR:   err,
		// 	NLOGQUERY: covid19.App.DB.QueryExpr()}, c.Get("user").(*jwt.Token), "", true)

		return ReturnInvalidResponse(http.StatusNoContent, err, "data edukasi kosong")
	}

	return c.JSON(http.StatusOK, result)
}

// // BankNew create new bank
// func BankNew(c echo.Context) error {
// 	defer c.Request().Body.Close()
// 	err := validatePermission(c, "core_bank_new")
// 	if err != nil {
// 		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
// 	}

// 	bank := models.Bank{}
// 	bankPayload := BankPayload{}

// 	payloadRules := govalidator.MapData{
// 		"name":     []string{"required"},
// 		"image":    []string{},
// 		"type":     []string{"required", "valid_id:bank_types"},
// 		"address":  []string{"required"},
// 		"province": []string{"required"},
// 		"city":     []string{"required"},
// 		"services": []string{"required", "valid_id:services"},
// 		"products": []string{"required", "valid_id:products"},
// 		"pic":      []string{"required"},
// 		"phone":    []string{"required"},
// 	}

// 	validate := validateRequestPayload(c, payloadRules, &bankPayload)
// 	if validate != nil {
// 		NLog("warning", "BankNew", map[string]interface{}{"message": "error validate new bank", "error": validate}, c.Get("user").(*jwt.Token), "", false)

// 		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "Kesalahan Validasi")
// 	}

// 	marshal, _ := json.Marshal(bankPayload)
// 	json.Unmarshal(marshal, &bank)

// 	if len(bankPayload.Image) > 0 {
// 		unbased, _ := base64.StdEncoding.DecodeString(bankPayload.Image)
// 		filename := "agt" + strconv.FormatInt(time.Now().Unix(), 10)
// 		url, err := covid19.App.S3.UploadJPEG(unbased, filename)
// 		if err != nil {
// 			NLog("error", "BankNew", map[string]interface{}{"message": fmt.Sprintf("error upload image bank %v", bank.ID), "error": err}, c.Get("user").(*jwt.Token), "", false)

// 			return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat bank baru")
// 		}

// 		bank.Image = url
// 	}

// 	err = bank.Create()
// 	middlewares.SubmitKafkaPayload(bank, "bank_create")
// 	if err != nil {
// 		NLog("error", "BankNew", map[string]interface{}{"message": fmt.Sprintf("error submitting kafka bank %v", bank.ID), "error": err, "bank": bank}, c.Get("user").(*jwt.Token), "", false)

// 		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat bank baru")
// 	}

// 	NAudittrail(models.Bank{}, bank, c.Get("user").(*jwt.Token), "bank", fmt.Sprint(bank.ID), "create")

// 	return c.JSON(http.StatusCreated, bank)
// }

// // BankDetail get bank detail by id
// func BankDetail(c echo.Context) error {
// 	defer c.Request().Body.Close()
// 	err := validatePermission(c, "core_bank_detail")
// 	if err != nil {
// 		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
// 	}

// 	db := covid19.App.DB

// 	bankID, _ := strconv.Atoi(c.Param("bank_id"))

// 	db = db.Table("banks").
// 		Select("banks.*, bt.name as bank_type_name").
// 		Joins("INNER JOIN bank_types bt ON banks.type = bt.id").
// 		Where("banks.id = ?", bankID)

// 	bank := BankSelect{}
// 	err = db.Find(&bank).Error
// 	if err != nil {
// 		NLog("warning", "BankDetail", map[string]interface{}{"message": fmt.Sprintf("error finding bank %v", bankID), "error": err}, c.Get("user").(*jwt.Token), "", false)

// 		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("Bank type %v tidak ditemukan", bankID))
// 	}

// 	return c.JSON(http.StatusOK, bank)
// }

// // BankPatch edit bank by id
// func BankPatch(c echo.Context) error {
// 	defer c.Request().Body.Close()
// 	err := validatePermission(c, "core_bank_patch")
// 	if err != nil {
// 		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
// 	}

// 	bankID, _ := strconv.ParseUint(c.Param("bank_id"), 10, 64)

// 	bank := models.Bank{}
// 	bankPayload := BankPayload{}
// 	err = bank.FindbyID(bankID)
// 	if err != nil {
// 		NLog("warning", "BankPatch", map[string]interface{}{"message": fmt.Sprintf("error finding bank %v", bankID), "error": err}, c.Get("user").(*jwt.Token), "", false)

// 		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("Bank %v tidak ditemukan", bankID))
// 	}
// 	origin := bank

// 	payloadRules := govalidator.MapData{
// 		"name":     []string{},
// 		"image":    []string{},
// 		"type":     []string{"valid_id:bank_types"},
// 		"address":  []string{},
// 		"province": []string{},
// 		"city":     []string{},
// 		"services": []string{"valid_id:services"},
// 		"products": []string{"valid_id:products"},
// 		"pic":      []string{},
// 		"phone":    []string{},
// 	}

// 	validate := validateRequestPayload(c, payloadRules, &bankPayload)
// 	if validate != nil {
// 		NLog("warning", "BankPatch", map[string]interface{}{"message": "error validation", "error": validate}, c.Get("user").(*jwt.Token), "", false)

// 		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "Kesalahan Validasi")
// 	}

// 	if len(bankPayload.Name) > 0 {
// 		bank.Name = bankPayload.Name
// 	}
// 	if bankPayload.Type > 0 {
// 		bank.Type = bankPayload.Type
// 	}
// 	if len(bankPayload.Address) > 0 {
// 		bank.Address = bankPayload.Address
// 	}
// 	if len(bankPayload.Province) > 0 {
// 		bank.Province = bankPayload.Province
// 	}
// 	if len(bankPayload.City) > 0 {
// 		bank.City = bankPayload.City
// 	}
// 	if len(bankPayload.Services) > 0 {
// 		bank.Services = pq.Int64Array(bankPayload.Services)
// 	}
// 	if len(bankPayload.Products) > 0 {
// 		bank.Products = pq.Int64Array(bankPayload.Products)
// 	}
// 	if len(bankPayload.PIC) > 0 {
// 		bank.PIC = bankPayload.PIC
// 	}
// 	if len(bankPayload.Phone) > 0 {
// 		bank.Phone = bankPayload.Phone
// 	}
// 	if len(bankPayload.Image) > 0 {
// 		unbased, _ := base64.StdEncoding.DecodeString(bankPayload.Image)
// 		filename := "agt" + strconv.FormatInt(time.Now().Unix(), 10)
// 		url, err := covid19.App.S3.UploadJPEG(unbased, filename)
// 		if err != nil {
// 			return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat bank baru")
// 		}

// 		i := strings.Split(bank.Image, "/")
// 		delImage := i[len(i)-1]
// 		err = covid19.App.S3.DeleteObject(delImage)
// 		if err != nil {
// 			log.Printf("failed to delete image %v from s3 bucket", delImage)
// 		}

// 		bank.Image = url
// 	}

// 	err = middlewares.SubmitKafkaPayload(bank, "bank_update")
// 	if err != nil {
// 		NLog("error", "BankPatch", map[string]interface{}{"message": fmt.Sprintf("error submitting bank %v", bank.ID), "error": err, "bank": bank}, c.Get("user").(*jwt.Token), "", false)

// 		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update bank %v", bankID))
// 	}

// 	NAudittrail(origin, bank, c.Get("user").(*jwt.Token), "bank type", fmt.Sprint(bank.ID), "update")

// 	return c.JSON(http.StatusOK, bank)
// }

// // BankDelete delete bank
// func BankDelete(c echo.Context) error {
// 	defer c.Request().Body.Close()
// 	err := validatePermission(c, "core_bank_delete")
// 	if err != nil {
// 		return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
// 	}

// 	bankID, _ := strconv.ParseUint(c.Param("bank_id"), 10, 64)

// 	bank := models.Bank{}
// 	err = bank.FindbyID(bankID)
// 	if err != nil {
// 		NLog("warning", "BankDelete", map[string]interface{}{"message": fmt.Sprintf("error finding bank %v", bankID), "error": err}, c.Get("user").(*jwt.Token), "", false)

// 		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("Bank type %v tidak ditemukan", bankID))
// 	}

// 	err = middlewares.SubmitKafkaPayload(bank, "bank_delete")
// 	if err != nil {
// 		NLog("error", "BankDelete", map[string]interface{}{"message": fmt.Sprintf("error submitting kafka bank %v", bankID), "error": err, "bank": bank}, c.Get("user").(*jwt.Token), "", false)

// 		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update bank tipe %v", bankID))
// 	}

// 	NAudittrail(bank, models.Bank{}, c.Get("user").(*jwt.Token), "bank", fmt.Sprint(bank.ID), "delete")

// 	return c.JSON(http.StatusOK, "ok")
// }
