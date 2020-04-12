package models

import (
	"covid19kalteng/components/basemodel"
)

// Edu main type
type Edu struct {
	basemodel.BaseModel
	Title       string `json:"title" gorm:"column:title"`
	Description string `json:"description" gorm:"column:description"`
}

// Create func
func (model *Edu) Create() error {
	return basemodel.Create(&model)
}

// Save func
func (model *Edu) Save() error {
	return basemodel.Save(&model)
}

// FirstOrCreate create if not exist, or skip if exist
func (model *Edu) FirstOrCreate() error {
	return basemodel.FirstOrCreate(&model)
}

// Delete func
func (model *Edu) Delete() error {
	return basemodel.Delete(&model)
}

// FindbyID func
func (model *Edu) FindbyID(id uint64) error {
	return basemodel.FindbyID(&model, id)
}

// SingleFindFilter func
func (model *Edu) FindSingle(filter interface{}) error {
	return basemodel.SingleFindFilter(&model, filter)
}

// PagedFindFilter func
func (base *Edu) FindPaged(filter interface{}) (basemodel.PagedFindResult, error) {
	edu := []Edu{}

	return base.PagedFindFilter(&edu, filter)
}
