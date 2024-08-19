package database

import (
	"gorm.io/gorm"
)

var db *gorm.DB

func InitialDatabase(database *gorm.DB) {
	db = database

}

func GetDB() *gorm.DB {
	return db
}
