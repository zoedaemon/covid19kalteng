package cases

import (
	"covid19kalteng/models"
	"encoding/json"
	"log"

	"github.com/jinzhu/gorm/dialects/postgres"
)

//CasePayload main payload for creation
type CasePayload struct {
	CreatedAt  string         `json:"date"`
	Location   postgres.Jsonb `json:"location"`
	DataDetail postgres.Jsonb `json:"data_detail"`
}

type DataDetailDefs map[string]interface{}

type LocationDefs struct {
	ProvinsiMain  string `json:"provinsi_main,omitempty"`
	Provinsi      string `json:"provinsi,omitempty"`
	KotaKabupaten string `json:"kota_kabupaten,omitempty"`
}

//ProcessCase custom processing jsonb data
func ProcessCase(cases *models.Case) error {
	var dataDetails []DataDetailDefs
	var newData []DataDetailDefs
	var location LocationDefs

	SampleKeyMapped := map[string]string{
		"positif": "Positif",
		"pdp":     "PDP",
		"odp":     "ODP",
		"sembuh":  "Sembuh",
	}

	json.Unmarshal(cases.Location.RawMessage, &location)
	json.Unmarshal(cases.DataDetail.RawMessage, &dataDetails)

	log.Println("XXXXXXXXXX location = ", location)
	for i, dat := range dataDetails {
		MainCaption := SampleKeyMapped[dat["key"].(string)]
		HarianKey := MainCaption + " Harian"

		Total := dat[MainCaption].(float64)
		Total -= 10
		dat[HarianKey] = Total
		log.Println(i, "XXXXXXXXXX dataDetail = ", dat)
		newData = append(newData, dat)
	}
	converted, err := json.Marshal(newData)
	if err != nil {
		return err
	}

	cases.DataDetail.RawMessage = converted

	return nil
}
