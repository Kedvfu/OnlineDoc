package models

import (
	"OnlineDoc/database"
	"time"
)

type DocumentInfo struct {
	DocumentID int `json:"document_id" gorm:"primary_key;auto_increment"`

	User   User `gorm:"reference:user_id"`
	UserId int  `json:"user_id" gorm:"foreignKey:user_id;not null"`

	Title        string    `json:"title" gorm:"type:varchar(255)"`
	Created      time.Time `json:"created" gorm:"type:datetime"`
	Updated      time.Time `json:"updated" gorm:"type:datetime"`
	DocumentType int       `json:"document_type" gorm:"type:int"`
	ShareUrl     string    `gorm:"type:varchar(255)"`
}

//	func InitializeDocumentInfo() {
//		db := database.GetDB()
//		err := db.AutoMigrate(&DocumentInfo{})
//		if err != nil {
//			return
//		}
//	}
func (documentInfo *DocumentInfo) TableName() string {
	return "t_document_info"
}

func GetDocumentTypeByTypeName(typeName string) int {
	switch typeName {
	case "markdown":
		return 1
	case "excel":
		return 2
	default:
		return 0
	}
}

func (documentInfo *DocumentInfo) Add() int {
	db := database.GetDB()
	newDocumentInfo := db.Create(documentInfo)
	if newDocumentInfo.Error != nil {
		return -1
	} else {
		return documentInfo.DocumentID
	}

}

func GetDocumentInfoById(documentId int) (*DocumentInfo, error) {
	db := database.GetDB()
	var documentInfo DocumentInfo
	db.First(&documentInfo, documentId)
	return &documentInfo, nil
}

func UpdateTitleByDocumentId(documentId int, title string) error {
	db := database.GetDB()
	var documentInfo DocumentInfo
	err := db.Model(&documentInfo).Where("document_id = ?", documentId).Update("title", title).Error
	return err
}
func UpdateShareUrlByDocumentId(documentId int, shareUrl string) error {
	db := database.GetDB()
	var documentInfo DocumentInfo
	err := db.Model(&documentInfo).Where("document_id = ?", documentId).Update("share_url", shareUrl).Error
	return err
}

func GetDocumentIdByShareUrl(shareUrl string) (int, error) {
	db := database.GetDB()
	var documentInfo DocumentInfo
	err := db.First(&documentInfo, "share_url = ?", shareUrl).Error
	return documentInfo.DocumentID, err

}
