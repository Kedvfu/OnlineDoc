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
type DocumentPermissionJson struct {
	UserId         int  `json:"user_id"`
	PermissionType bool `json:"permission_type"`
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

func (documentPermission *DocumentPermission) Add() (bool, error) {
	db := database.GetDB()
	err := db.Where("document_id =? AND user_id =?", documentPermission.DocumentId, documentPermission.UserId).First(&DocumentPermission{}).Error
	if err == nil {
		return true, nil
	}
	return false, db.Create(documentPermission).Error
}
func UpdateDocumentPermissionTypeByDocumentIdAndUserId(documentId int, userId int, permissionType bool) error {
	db := database.GetDB()
	err := db.Model(&DocumentPermission{}).Where("document_id =? AND user_id =?", documentId, userId).Update("permission_type", permissionType).Error
	if err != nil {
		return err
	}
	return nil
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

/*
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
*/

func GetPermissionTypeAndUserIdByDocumentId(documentId int) ([]byte, error) {
	db := database.GetDB()
	var documentPermissions []DocumentPermission
	db.Where("document_id = ? ", documentId).Find(&documentPermissions)

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
func DeleteDocumentPermissionByDocumentIdAndUserId(documentId int, userId int) error {
	db := database.GetDB()
	var documentPermissions []DocumentPermission
	err := db.Table("t_document_permission").
		Joins("left join t_document_info on t_document_info.document_id = t_document_permission.document_id").
		Find(&documentPermissions).Where("document_id = ? AND user_id = ?", documentId, userId).Delete(&documentPermissions).Error
	if err != nil {
		return err
	}

	return nil
}
