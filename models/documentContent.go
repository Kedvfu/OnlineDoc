package models

import (
	"OnlineDoc/database"
	"time"
)

type DocumentContent struct {
	ContentId int `json:"contentId" gorm:"primary_key;auto_increment"`

	DocumentInfo DocumentInfo `gorm:"foreignKey:document_id"`
	DocumentId   int          `json:"document_id" gorm:"not null"`

	Updated time.Time `json:"updated" gorm:"column:updated"`

	User   User `gorm:"foreignKey:user_id"`
	UserId int  `json:"user_id" gorm:"not null"`

	Content string `json:"content" gorm:"type:text"`
}

//func InitializeDocumentContent() {
//	db := database.GetDB()
//	err := db.AutoMigrate(&DocumentContent{})
//	if err != nil {
//		return
//	}
//}

func (documentContent *DocumentContent) TableName() string {
	return "t_document_content"
}

func (documentContent *DocumentContent) Add() int {
	db := database.GetDB()
	newDocumentInfo := db.Create(documentContent)
	if newDocumentInfo.Error != nil {
		return -1
	} else {
		return documentContent.ContentId
	}
}

func GetLatestDocumentContent(documentId int) (*DocumentContent, error) {
	db := database.GetDB()
	var documentContent DocumentContent
	err := db.Model(&documentContent).Where("document_id =?", documentId).Order("updated desc").First(&documentContent).Error
	if err != nil {
		return nil, err
	}
	return &documentContent, nil
}
