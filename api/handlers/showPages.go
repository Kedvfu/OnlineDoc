package handlers

import (
	"OnlineDoc/api/sessions"
	"OnlineDoc/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func ShowHomepage(context *gin.Context) {
	userId, exists := context.Get("userId")
	//var documentWithPermission []byte
	if !exists {
		userId = ""
	} else {
		userId, _ = strconv.Atoi(userId.(string))
	}
	sessionToken, exists := context.Get("session_token")
	if !exists {
		sessionToken = ""
	} else {
		sessionToken = sessionToken.(string)
		//var err error
		//documentWithPermission, err = models.GetPermissionTypeAndDocumentIdByUserId(userId.(int))
		//if err != nil {
		//	log.Printf("Unable to get permission type and document id for user %d", userId)
		//}
		//
		//context.Set("document_map", documentWithPermission)
	}

	context.HTML(http.StatusOK, "home.html", gin.H{
		"session_token": sessionToken,
		"user_id":       userId,
		//"document_with_permission": documentWithPermission,
	})

}

func ShowDefaultPage(context *gin.Context) {
	context.Redirect(http.StatusMovedPermanently, "/home")
}

func ShowLoginPage(context *gin.Context) {
	context.HTML(http.StatusOK, "login.html", gin.H{})
}
func ShowLogoutPage(context *gin.Context) {
	context.HTML(http.StatusOK, "logout.html", gin.H{})
	sessionToken, exists := context.Get("session_token")
	if exists {
		delResult := RedisClient.Del(sessionToken.(string))
		if delResult.Err() != nil {
			log.Printf("Unable to delete session %s token from redis", sessionToken)
		} else {
			log.Printf("Session %s token deleted from redis", sessionToken)
		}
	}

}

func ShowDocumentPage(context *gin.Context) {
	documentId, _ := context.Get("documentId")
	permissionType, _ := context.Get("permissionType")
	title, _ := context.Get("title")
	documentType, _ := context.Get("documentType")
	userId, _ := context.Get("userId")
	authorId, _ := context.Get("authorId")
	authorIdNum, _ := authorId.(int)

	if documentType == 2 {
		documentIdNum := documentId.(int)
		_, exists := (*sessions.ExcelSessions)[documentIdNum]
		if !exists {
			documentContent, err := models.GetLatestDocumentContent(documentIdNum)
			if err != nil {
				context.JSON(200, gin.H{
					"message": "Unable to get document content",
				})
				context.Abort()
				return
			}
			var excelData models.ExcelData
			if documentContent.Content == "" {
				excelData = *models.GetEmptyExcelData()
				excelDataJson, err := json.Marshal(excelData)
				documentContent.Content = string(excelDataJson)

				if err != nil {
					context.JSON(200, gin.H{
						"message": "Unable to get document content",
					})
					context.Abort()
					return
				}
			}

			err = json.Unmarshal([]byte(documentContent.Content), &excelData)

			if err != nil {
				context.JSON(200, gin.H{
					"message": "Unable to parse document content",
				})
				context.Abort()
				return
			}

			currentExcelData := (*sessions.ExcelSessions)[documentIdNum]
			if currentExcelData == nil {
				newCellHistory := make([]models.CellHistory, 0)
				excelData.OnlineUsers = &newCellHistory
				(*sessions.ExcelSessions)[documentIdNum] = &excelData
			}
		}
	}

	context.HTML(200, "document.html", gin.H{
		"author_id":      authorIdNum,
		"user_id":        userId,
		"document_id":    documentId,
		"permissionType": permissionType,
		"title":          title,
		"documentType":   documentType,
	})
}
func ShowRegisterPage(context *gin.Context) {
	context.HTML(http.StatusOK, "register.html", gin.H{})

}
