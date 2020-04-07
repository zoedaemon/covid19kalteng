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

type DataDetailDefs struct {
	Key string `json:"key"`
}

type LocationDefs struct {
	ProvinsiMain  string `json:"provinsi_main,omitempty"`
	Provinsi      string `json:"provinsi,omitempty"`
	KotaKabupaten string `json:"kota_kabupaten,omitempty"`
}

//ProcessCase custom processing jsonb data
func ProcessCase(cases *models.Case) {
	var dataDetails []DataDetailDefs
	var location LocationDefs

	json.Unmarshal(cases.Location.RawMessage, &location)
	json.Unmarshal(cases.DataDetail.RawMessage, &dataDetails)

	log.Println("XXXXXXXXXX location = ", location)
	for i, dat := range dataDetails {
		log.Println(i, "XXXXXXXXXX dataDetail = ", dat)
	}
}
