package handlers

import (
	"OnlineDoc/api/authenticate"
	"OnlineDoc/api/sessions"
	"OnlineDoc/models"
	"crypto/md5"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func GetUserDocuments(context *gin.Context) {
	userId := context.Param("userId")
	userIdNum, err := strconv.Atoi(userId)
	if err != nil {
		context.JSON(200, gin.H{
			"message": "Invalid user id",
		})
		return
	}
	Documents, err := models.GetPermissionTypeAndDocumentIdByUserId(userIdNum)
	if err != nil {
		context.JSON(200, gin.H{
			"message": "Unable to get documents for user",
		})
	}
	if len(Documents) == 0 {
		context.JSON(200, gin.H{
			"message": "No documents found for user",
		})
	}
	context.Data(200, "application/json", Documents)
}

func SaveDocument(context *gin.Context) {
	//userId := context.Param("userId")
	userId, _ := context.Get("userId")
	userIdNum := userId.(int)
	documentId := context.Param("documentId")
	documentIdNum, err := strconv.Atoi(documentId)
	if err != nil {
		context.JSON(200, gin.H{
			"message": "Invalid document id",
		})
		return
	}

	var newTitle string
	var documentContent models.DocumentContent

	documentInfo, err := models.GetDocumentInfoById(documentIdNum)
	if err != nil {
		context.JSON(200, gin.H{
			"message": "Unable to get document info",
		})
		return
	}
	documentType := documentInfo.DocumentType

	//
	var jsonData map[string]interface{}
	err = context.ShouldBindJSON(&jsonData)
	if err != nil {
		context.JSON(200, gin.H{
			"message": "Invalid document content",
		})
		return
	}
	if documentType == 1 {
		if value, ok := jsonData["content"]; ok {
			documentContent.Content = value.(string)
		} else {
			context.JSON(200, gin.H{
				"message": "Invalid document content",
			})
			return
		}
	} else if documentType == 2 {
		excelContent, err := json.Marshal((*sessions.ExcelSessions)[documentIdNum])
		if err != nil {
			context.JSON(200, gin.H{
				"message": "Unable to parse document content",
			})
			return
		}
		documentContent.Content = string(excelContent)
	}

	if value, ok := jsonData["title"]; ok {
		newTitle = value.(string)
	}

	if documentContent.Content == "" {
		context.JSON(200, gin.H{
			"message": "Document content cannot be empty",
		})
		return
	}
	//

	ContentHash := md5.Sum([]byte(documentContent.Content))
	currentContentHash := RedisClient.Get(documentId).Val()
	if currentContentHash == string(ContentHash[:]) {
		context.JSON(200, gin.H{
			"message": "No changes detected in document content",
		})
		return
	}
	documentContent.DocumentId, err = strconv.Atoi(documentId)
	if err != nil {
		context.JSON(200, gin.H{
			"message": "Invalid document id",
		})
		return
	}
	documentContent.UserId = userIdNum
	documentContent.Updated = time.Now()
	documentContent.Add()
	if err != nil {
		context.JSON(200, gin.H{
			"message": "Unable to save document content",
		})
		return
	}
	RedisClient.Set(documentId, string(ContentHash[:]), time.Second*3600)

	if newTitle != "" {
		err := models.UpdateTitleByDocumentId(documentIdNum, newTitle)
		if err != nil {
			context.JSON(200, gin.H{
				"message": "Unable to update document title",
			})
			return
		}

	}

}
func ShowDocumentFromSharePage(context *gin.Context) {
	sessionToken, err := context.Cookie("session_token")
	if err != nil {
		context.Redirect(200, "/login")
		return
	}
	trueUserId, err := RedisClient.Get(sessionToken).Result()
	if err != nil {
		context.Redirect(200, "/login")
		return
	}
	trueUserIdNum, err := strconv.Atoi(trueUserId)
	shareUrl := context.Param("shareUrl")
	if shareUrl == "" {
		context.Redirect(200, "/home")
		return
	}
	documentId, err := models.GetDocumentIdByShareUrl(shareUrl)
	if err != nil {
		context.Redirect(200, "/home")
		return
	}
	documentPermission := models.DocumentPermission{
		UserId:         trueUserIdNum,
		DocumentId:     documentId,
		PermissionType: false,
	}
	err = documentPermission.Add()
	if err != nil {
		context.JSON(200, gin.H{
			"message": "unable to add permission for user",
		})
		return
	}
	context.Redirect(301, "/document/"+strconv.Itoa(documentId))
}
func GetDocument(context *gin.Context) {

	documentId := context.Param("documentId")

	documentIdNum, err := strconv.Atoi(documentId)
	if err != nil {
		context.JSON(200, gin.H{
			"message": "Invalid document id",
		})
		return
	}
	//documentInfo, err := models.GetDocumentInfoById(documentIdNum)
	//if err != nil {
	//	context.JSON(200, gin.H{
	//		"message": "Unable to get document info",
	//	})
	//	return
	//}
	documentUsersJson, _ := models.GetPermissionTypeAndUserIdByDocumentId(documentIdNum)
	var jsonData []interface{}
	err = json.Unmarshal(documentUsersJson, &jsonData)
	if err != nil {
		context.JSON(200, gin.H{
			"message": "Unable to parse document users",
		})
	}
	var contentData string
	if value, ok := (*sessions.ExcelSessions)[documentIdNum]; ok {
		contentDataJson, err := json.Marshal(value)
		if err != nil {
			context.JSON(200, gin.H{
				"message": "Unable to parse document content",
			})
			return
		}
		contentData = string(contentDataJson)
	} else {
		documentContent, err := models.GetLatestDocumentContent(documentIdNum)
		if err != nil {
			context.JSON(200, gin.H{
				"message": "Unable to get document content",
			})
			return
		}
		if documentContent == nil {
			context.JSON(200, gin.H{
				"message": "No document content found",
			})
			return
		}
		contentData = documentContent.Content
	}
	context.JSON(200, gin.H{
		"content": contentData,
		//"updated":       documentContent.Updated,
		"documentUsers": jsonData,
		//"documentType": documentInfo.DocumentType,
	})

}
func GetDocumentLink(context *gin.Context) {
	documentId := context.Param("documentId")
	documentIdNum, err := strconv.Atoi(documentId)
	if err != nil {
		context.JSON(200, gin.H{
			"message": "Invalid document id",
		})
		return
	}
	userId, _ := context.Get("userId")
	userIdNum := userId.(int)
	PermissionType, err := models.GetPermissionTypeByDocumentIdAndUserId(documentIdNum, userIdNum)
	if err != nil {
		context.JSON(200, gin.H{
			"message": "Unable to get permission type for user",
		})
		return
	}
	if PermissionType == 0 {
		context.JSON(200, gin.H{
			"message": "No permission to create link",
		})
		return
	} else if PermissionType == 1 {
		documentInfo, _ := models.GetDocumentInfoById(documentIdNum)
		link := documentInfo.ShareUrl
		if link == "" {
			link, _ = authenticate.GenerateSessionToken()
			err := models.UpdateShareUrlByDocumentId(documentIdNum, link)
			if err != nil {
				context.JSON(200, gin.H{
					"message": "Unable to generate link",
				})
				return
			}

		}

		context.JSON(200, gin.H{
			"link": link,
		})
	}
}
