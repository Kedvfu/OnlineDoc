package models

import (
	"OnlineDoc/database"
	"encoding/json"
)

type DocumentPermission struct {
	User   User `gorm:"foreignKey:user_id"`
	UserId int  `json:"user_id"`

	DocumentInfo DocumentInfo `gorm:"foreignKey:document_id"`
	DocumentId   int          `json:"document_id"`

	PermissionType bool `json:"permission_type" gorm:"type:boolean"` //true for write, false for read
}

//	func InitializeDocumentPermission() {
//		db := database.GetDB()
//		err := db.AutoMigrate(&DocumentPermission{})
//		if err != nil {
//			return
//		}
//	}
func (documentPermission *DocumentPermission) TableName() string {
	return "t_document_permission"
}

func (documentPermission *DocumentPermission) Add() error {
	db := database.GetDB()
	err := db.Where("document_id =? AND user_id =?", documentPermission.DocumentId, documentPermission.UserId).First(&DocumentPermission{}).Error
	if err == nil {
		return nil
	}
	return db.Create(documentPermission).Error
}

func GetPermissionTypeByDocumentIdAndUserId(documentId int, userId int) (int, error) {
	db := database.GetDB()
	var permissionType bool
	err := db.Model(&DocumentPermission{}).Where("document_id =? AND user_id =?", documentId, userId).Select("permission_type").First(&permissionType).Error
	if err != nil {
		return -1, err
	}
	if permissionType {
		return 1, nil
	} else {
		return 0, nil
	}
}
func GetPermissionTypeAndDocumentIdByUserId(userId int) ([]byte, error) {
	db := database.GetDB()
	var documentPermissions []DocumentPermission
	db.Where("user_id = ? ", userId).Find(&documentPermissions)
	type DocumentPermissionJson struct {
		DocumentId     int  `json:"document_id"`
		PermissionType bool `json:"permission_type"`
	}

	documentPermissionJsons := make([]DocumentPermissionJson, 0)
	for _, documentPermission := range documentPermissions {
		documentPermissionJsons = append(documentPermissionJsons, DocumentPermissionJson{
			DocumentId:     documentPermission.DocumentId,
			PermissionType: documentPermission.PermissionType,
		})
	}

	bytes, err := json.Marshal(documentPermissionJsons)
	if err != nil {
		return nil, err
	}
	return bytes, nil

}
func GetPermissionTypeAndUserIdByDocumentId(documentId int) ([]byte, error) {
	db := database.GetDB()
	var documentPermissions []DocumentPermission
	db.Where("document_id = ? ", documentId).Find(&documentPermissions)
	type DocumentPermissionJson struct {
		UserId         int  `json:"user_id"`
		PermissionType bool `json:"permission_type"`
	}
	documentPermissionJsons := make([]DocumentPermissionJson, 0)
	for _, documentPermission := range documentPermissions {
		documentPermissionJsons = append(documentPermissionJsons, DocumentPermissionJson{
			UserId:         documentPermission.UserId,
			PermissionType: documentPermission.PermissionType,
		})
	}

	bytes, err := json.Marshal(documentPermissionJsons)
	if err != nil {
		return nil, err
	}
	return bytes, nil

}
