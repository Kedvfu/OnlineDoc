package handlers

import (
	"OnlineDoc/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

func GetUserInfo(context *gin.Context) {
	targetUserID := context.Param("targetUserId")
	targetUserIDList := strings.Split(targetUserID, ";")
	userList := make([]models.User, 0)
	for _, userID := range targetUserIDList {
		userIdNum, err := strconv.Atoi(userID)
		if err != nil {
			context.JSON(200, gin.H{
				"message": "Invalid user ID",
			})
			context.Abort()
			return
		}
		userJson, _ := RedisClient.Get("user_" + userID).Result()
		if userJson == "" {
			user, err := models.GetUserByUserId(userIdNum)
			if err != nil {
				context.JSON(200, gin.H{
					"message": "User not found",
				})
				context.Abort()
				return
			}
			userList = append(userList, user)
			userJson, err := json.Marshal(user)
			if err != nil {
				context.JSON(200, gin.H{
					"message": "Internal server error",
				})
				context.Abort()
				return
			}
			RedisClient.Set("user_"+userID, userJson, 0)
		} else {
			user := models.User{}
			err := json.Unmarshal([]byte(userJson), &user)
			if err != nil {
				context.JSON(200, gin.H{
					"message": "Internal server error",
				})
				context.Abort()
				return
			}
			userList = append(userList, user)
			RedisClient.Expire("user_"+userID, 0)
		}
	}
	context.JSON(200, userList)
}

func UpdateUserPermissionType(context *gin.Context) {
	targetUserId := context.Param("targetUserId")
	targetUserIdNum, err := strconv.Atoi(targetUserId)
	if err != nil {
		context.JSON(200, gin.H{
			"message": "Invalid target user ID",
		})
		context.Abort()
		return
	}

	documentId := context.Param("documentId")
	documentIdNum, err := strconv.Atoi(documentId)
	if err != nil {
		context.JSON(200, gin.H{
			"message": "Invalid document ID",
		})
		context.Abort()
		return
	}
	permissionType := context.Param("permissionType")
	var permissionTypeBool bool
	if permissionType == "true" {
		permissionTypeBool = true
	} else if permissionType == "false" {
		permissionTypeBool = false
	} else {
		context.JSON(200, gin.H{
			"message": "Invalid permission type",
		})
		context.Abort()
		return
	}
	err = models.UpdateDocumentPermissionTypeByDocumentIdAndUserId(documentIdNum, targetUserIdNum, permissionTypeBool)
	if err != nil {
		context.JSON(200, gin.H{
			"message": "Internal server error",
		})
		context.Abort()
		return

	}
	context.JSON(200, gin.H{
		"status":  1,
		"message": "Permission type updated successfully",
	})

}
