package middleware

import (
	"OnlineDoc/models"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func DocumentMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		userId, exists := context.Get("userId")
		userId, _ = strconv.Atoi(userId.(string))
		if !exists {
			context.AbortWithStatusJSON(401, gin.H{"message": "Unauthorized"})
			return
		} else {
			var documentIdResult int
			var permissionTypeResult bool
			var titleResult string
			var documentTypeResult int
			//var documentContentResult string

			documentId := context.Param("documentId")
			documentType := context.Param("documentType")

			documentTypeResult = models.GetDocumentTypeByTypeName(documentType)

			if documentId == "new" && documentType != "" {
				// create new document
				documentInfo := models.DocumentInfo{
					Title:        "新文档",
					Created:      time.Now(),
					Updated:      time.Now(),
					DocumentType: documentTypeResult,
					UserId:       userId.(int),
				}

				documentId := documentInfo.Add()
				if documentId == -1 {
					context.AbortWithStatusJSON(500, gin.H{"message": "Failed to create new document"})
					return
				}
				documentContent := models.DocumentContent{
					DocumentId: documentId,
					Updated:    time.Now(),
					UserId:     userId.(int),
					Content:    "",
				}
				ContentId := documentContent.Add()
				if ContentId == -1 {
					context.AbortWithStatusJSON(500, gin.H{"message": "Failed to create new document content"})
					return
				}
				documentPermission := models.DocumentPermission{
					UserId:         userId.(int),
					DocumentId:     documentId,
					PermissionType: true,
				}
				err := documentPermission.Add()
				if err != nil {
					context.AbortWithStatusJSON(500, gin.H{"message": "Failed to create new document permission"})
				}

				titleResult = documentInfo.Title
				documentIdResult = documentId
				permissionTypeResult = true
				//documentContentResult = ""
			} else {

				documentId, err := strconv.Atoi(documentId)
				if err != nil {
					context.AbortWithStatusJSON(400, gin.H{"message": "Invalid document id"})
					return
				}
				permissionType, err := models.GetPermissionTypeByDocumentIdAndUserId(documentId, userId.(int))
				if err != nil {
					context.AbortWithStatusJSON(500, gin.H{"message": "No permission for this document"})
					return
				}

				documentInfo, err := models.GetDocumentInfoById(documentId)
				if err != nil {
					context.AbortWithStatusJSON(500, gin.H{"message": "Failed to get document info"})
					return
				}
				titleResult = documentInfo.Title
				documentTypeResult = documentInfo.DocumentType

				//documentContent, err := models.GetLatestDocumentContent(documentId)
				//if err != nil {
				//	context.AbortWithStatusJSON(500, gin.H{"message": "Failed to get document content"})
				//	return
				//}
				//documentContentResult = documentContent.Content

				documentIdResult = documentId
				if permissionType == 1 {
					permissionTypeResult = true
				} else {
					permissionTypeResult = false
				}
			}
			userIdNum, _ := userId.(int)

			context.Set("userId", userIdNum)
			context.Set("documentId", documentIdResult)
			context.Set("permissionType", permissionTypeResult)
			context.Set("title", titleResult)
			context.Set("documentType", documentTypeResult)
			//context.Set("documentContent", documentContentResult)
		}

	}

}
