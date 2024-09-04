package handlers

import (
	"OnlineDoc/api/authenticate"
	"OnlineDoc/api/sessions"
	"OnlineDoc/models"
	"crypto/md5"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
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
	documentInfos, err := models.GetDocumentInfoByPermissionTypeByUserId(userIdNum)
	if err != nil {
		context.JSON(200, gin.H{
			"message": "Unable to get documents for user",
		})
		context.Abort()
		return
	}
	if len(*documentInfos) == 0 {
		context.JSON(200, gin.H{
			"message": "No documents found for user",
		})
		context.Abort()
		return
	}
	context.JSON(200, *documentInfos)
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
		context.Abort()
		return
	}

	var newTitle string
	var documentContent models.DocumentContent

	documentInfo, err := models.GetDocumentInfoById(documentIdNum)
	if err != nil {
		context.JSON(200, gin.H{
			"message": "Unable to get document info",
		})
		context.Abort()
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
		context.Abort()
		return
	}
	if documentType == 1 {
		if value, ok := jsonData["content"]; ok {
			documentContent.Content = value.(string)
		} else {
			context.JSON(200, gin.H{
				"message": "Invalid document content",
			})
			context.Abort()
			return
		}
	} else if documentType == 2 {
		excelContent, err := json.Marshal((*sessions.ExcelSessions)[documentIdNum])
		if err != nil {
			context.JSON(200, gin.H{
				"message": "Unable to parse document content",
			})
			context.Abort()
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
		context.Abort()
		return
	}
	//

	ContentHash := md5.Sum([]byte(documentContent.Content))
	currentContentHash := RedisClient.Get("documentHash_" + documentId).Val()
	if currentContentHash == string(ContentHash[:]) {
		context.JSON(200, gin.H{
			"message": "No changes detected in document content",
		})
		context.Abort()
		return
	}
	documentContent.DocumentId, err = strconv.Atoi(documentId)
	if err != nil {
		context.JSON(200, gin.H{
			"message": "Invalid document id",
		})
		context.Abort()
		return
	}
	documentContent.UserId = userIdNum
	documentContent.Updated = time.Now()

	documentContent.Add()

	if err != nil {
		context.JSON(200, gin.H{
			"message": "Unable to save document content",
		})
		context.Abort()
		return
	}
	RedisClient.Set("documentHash_"+documentId, string(ContentHash[:]), 0)

	if newTitle != "" {
		err := models.UpdateTitleByDocumentId(documentIdNum, newTitle)
		if err != nil {
			context.JSON(200, gin.H{
				"message": "Unable to update document title",
			})
			context.Abort()
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

	exists, err := documentPermission.Add()
	if err != nil {
		context.JSON(200, gin.H{
			"message": "unable to add permission for user",
		})
		context.Abort()
		return
	}
	if !exists {
		RedisClient.Del("documentUsers_" + strconv.Itoa(documentId))
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
		context.Abort()
		return
	}
	//documentInfo, err := models.GetDocumentInfoById(documentIdNum)
	//if err != nil {
	//	context.JSON(200, gin.H{
	//		"message": "Unable to get document info",
	//	})
	//	return
	//}
	var documentUsersJson []byte

	documentUsersJsonString, err := RedisClient.Get("documentUsers_" + documentId).Result()
	if documentUsersJsonString == "No data" {
		context.JSON(200, gin.H{
			"message": "Unable to get document users",
		})
		context.Abort()
		return
	}

	if errors.Is(err, redis.Nil) {
		documentUsersJson, err = models.GetPermissionTypeAndUserIdByDocumentId(documentIdNum)
		if documentUsersJson == nil {
			context.JSON(200, gin.H{
				"message": "Unable to get document users",
			})
			RedisClient.Set("documentUsers_"+documentId, "No data", 0)
			context.Abort()
			return
		}
		RedisClient.Set("documentUsers_"+documentId, string(documentUsersJson), 0)

	} else if err != nil {
		context.JSON(200, gin.H{
			"message": "Unable to get document users",
		})
		context.Abort()
		return
	} else {
		documentUsersJson = []byte(documentUsersJsonString)
	}
	var jsonData []interface{}
	err = json.Unmarshal(documentUsersJson, &jsonData)
	if err != nil {
		context.JSON(200, gin.H{
			"message": "Unable to parse document users",
		})
		context.Abort()
		return
	}
	var contentData string
	if value, ok := (*sessions.ExcelSessions)[documentIdNum]; ok {
		value.RWMutex.RLock()
		defer value.RWMutex.RUnlock()
		contentDataJson, err := json.Marshal(value)
		if err != nil {
			context.JSON(200, gin.H{
				"message": "Unable to parse document content",
			})
			context.Abort()
			return
		}
		contentData = string(contentDataJson)
	} else {
		documentContent, err := models.GetLatestDocumentContent(documentIdNum)
		if err != nil {
			context.JSON(200, gin.H{
				"message": "Unable to get document content",
			})
			context.Abort()
			return
		}
		if documentContent == nil {
			context.JSON(200, gin.H{
				"message": "No document content found",
			})
			context.Abort()
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
		context.Abort()
		return
	}
	userId, _ := context.Get("userId")
	userIdNum := userId.(int)
	PermissionType, err := models.GetPermissionTypeByDocumentIdAndUserId(documentIdNum, userIdNum)
	if err != nil {
		context.JSON(200, gin.H{
			"message": "Unable to get permission type for user",
		})
		context.Abort()
		return
	}
	if PermissionType == 0 {
		context.JSON(200, gin.H{
			"message": "No permission to create link",
		})
		context.Abort()
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
				context.Abort()
				return
			}

		}

		context.JSON(200, gin.H{
			"link": link,
		})
	}
}

func DeleteDocument(context *gin.Context) {
	userId, _ := context.Get("userId")
	userIdNum := userId.(int)
	documentId := context.Param("documentId")
	documentIdNum, err := strconv.Atoi(documentId)
	if err != nil {
		context.JSON(200, gin.H{
			"message": "Invalid document id",
		})
		context.Abort()
		return
	}
	err = models.DeleteDocumentPermissionByDocumentIdAndUserId(documentIdNum, userIdNum)

	if err != nil {

		context.JSON(200, gin.H{

			"message": "Unable to delete document for user",
		})
		context.Abort()
		return
	}
	context.JSON(200, gin.H{
		"status":  1,
		"message": "Document deleted",
	})
}
