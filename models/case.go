package models

import (
	"covid19kalteng/components/basemodel"

	"github.com/jinzhu/gorm/dialects/postgres"
)

// Case main type
type Case struct {
	basemodel.BaseModel
	Location   postgres.Jsonb `json:"location" gorm:"column:location"`
	DataDetail postgres.Jsonb `json:"data_detail" gorm:"column:data_detail"`
}

// Create func
func (model *Case) Create() error {
	return basemodel.Create(&model)
}

// Save func
func (model *Case) Save() error {
	return basemodel.Save(&model)
}

// FirstOrCreate create if not exist, or skip if exist
func (model *Case) FirstOrCreate() error {
	return basemodel.FirstOrCreate(&model)
}

// Delete func
func (model *Case) Delete() error {
	return basemodel.Delete(&model)
}

// FindbyID func
func (model *Case) FindbyID(id uint64) error {
	return basemodel.FindbyID(&model, id)
}

// SingleFindFilter func
func (model *Case) FindSingle(filter interface{}) error {
	return basemodel.SingleFindFilter(&model, filter)
}

// PagedFindFilter func
func (model *Case) FindPaged(page int, rows int, orderby []string, sort []string, filter interface{}) (basemodel.PagedFindResult, error) {
	cases := []Case{}

	return basemodel.PagedFindFilter(&cases, page, rows, orderby, sort, filter)
}
