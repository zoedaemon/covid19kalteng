package cases

import (
	"covid19kalteng/covid19"
	"covid19kalteng/models"
	. "covid19kalteng/modules"
	"fmt"

	"encoding/json"
	"log"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
)

//CasePayload main payload for creation
type CasePayload struct {
	CreatedAt  string         `json:"date"`
	Location   postgres.Jsonb `json:"location"`
	DataDetail postgres.Jsonb `json:"data_detail"`
}

//DataDetailDefs map custom json data data_detail
type DataDetailDefs map[string]interface{}

//LocationDefs mapping custom json data location
type LocationDefs struct {
	ProvinsiMain  string `json:"provinsi_main,omitempty"`
	Provinsi      string `json:"provinsi,omitempty"`
	KotaKabupaten string `json:"kota_kabupaten,omitempty"`
}

//Response from latest row to diff operation for getting correct per day data
type Response struct {
	//GOTCHAS: must defined gorm column type as jsonb, if not gorm never get the data
	DataDetail postgres.Jsonb `json:"data_detail" gorm:"column:data_detail;type:jsonb"`
}

//ProcessCase custom processing jsonb data
func ProcessCase(cases *models.Case, date time.Time) error {
	var dataDetails []DataDetailDefs
	var newData []DataDetailDefs
	var location LocationDefs
	var dataDetailFoundJSON Response
	var dataDetailFound []DataDetailDefs

	SampleKeyMapped := map[string]string{
		"positif": "Positif",
		"pdp":     "PDP",
		"odp":     "ODP",
		"sembuh":  "Sembuh",
	}

	json.Unmarshal(cases.Location.RawMessage, &location)
	json.Unmarshal(cases.DataDetail.RawMessage, &dataDetails)

	//get latest inserted detail data by the location
	db := covid19.App.DB
	db = db.Table("cases").
		Select("*") //data_detail#>'{}'
	if len(location.ProvinsiMain) > 0 {
		db = db.Where("location->>'provinsi_main' LIKE ?", location.ProvinsiMain)
	}
	if len(location.Provinsi) > 0 {
		db = db.Where("location->>'provinsi' LIKE ?", location.Provinsi)
	}
	if len(location.KotaKabupaten) > 0 {
		db = db.Where("location->>'kota_kabupaten' LIKE ?", location.KotaKabupaten)
	}
	//if not empty date do case processing before created_at = date
	if !date.IsZero() {
		db = db.Where("created_at <= ?", date)
	}

	//order by DESC for getting the latest row
	db = db.Order("created_at DESC LIMIT 1")

	//exec
	//BUGS "Record Not Found"
	err := db.Find(&dataDetailFoundJSON).Error
	if err != nil {
		log.Println("error : ", err)
		//NOTE: don't worry not return error here, coz if first row (no latest row exist)
		//		there's could always return error forever
	}

	//translate to object
	json.Unmarshal(dataDetailFoundJSON.DataDetail.RawMessage, &dataDetailFound)

	//dataDetailFoundJSON.DataDetail.RawMessage, dataDetailFound)
	if len(dataDetailFound) > 0 {

		for i, dat := range dataDetails {
			//get key/caption
			MainCaption := SampleKeyMapped[dat["key"].(string)]
			HarianCaption := MainCaption + " Harian"

			//latest total
			//must valid float (must handle with switch type, if not you would throws panic)
			Total, ok := dat[MainCaption].(float64)
			if !ok {
				return &ParseError{Info: fmt.Sprintf("field %s harus angka", MainCaption)}
			}
			log.Printf(">>>>> %s : %s\n", MainCaption, HarianCaption)

			//get Last total from db
			//CONCERN: this [i] might be a problem if from payload not sequenced as latest row
			if _, ok := dataDetailFound[i][MainCaption]; !ok {
				return &ParseError{Info: "fields data_detail tidak terurut sesuai standar"}
			}
			LastTotal := dataDetailFound[i][MainCaption].(float64)

			//calc new day changes (increment or decrement from this diff)
			TodayChange := Total - LastTotal

			dat[HarianCaption] = TodayChange
			log.Printf(">>>>> ProcessCase(%d) Last Total = %f; New Total = %f; Today = %f\n", i, Total,
				LastTotal, TodayChange)

			//append new updated object
			newData = append(newData, dat)
		}

		//raw []byte
		converted, err := json.Marshal(newData)
		if err != nil {
			return err
		}

		//replace old value
		cases.DataDetail.RawMessage = converted

	}

	return nil
}
