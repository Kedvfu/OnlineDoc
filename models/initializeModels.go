package models

import (
	"OnlineDoc/database"
	"log"
)

func InitializeModels() {
	db := database.GetDB()
	err := db.AutoMigrate(&User{})
	if err != nil {
		return
	}
	err = db.AutoMigrate(&DocumentInfo{})
	if err != nil {
		return
	}
	err = db.AutoMigrate(&DocumentPermission{})
	if err != nil {
		return
	}
	err = db.AutoMigrate(&DocumentContent{})
	if err != nil {
		return
	}
	rows, _ := db.Raw("SELECT table_name FROM information_schema.tables").Rows()
	log.Printf("Models initialized successfully %v", rows)
}
