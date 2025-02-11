package db

import (
	"gorm.io/gorm"
)

type SamGovData struct {
	gorm.Model
	// Define fields for SamGov data here
}

func PersistSamGovData(db *gorm.DB, data SamGovData) error {
	result := db.Create(&data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
