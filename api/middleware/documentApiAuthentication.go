package middleware

import (
	"OnlineDoc/api/handlers"
	"github.com/gin-gonic/gin"
	"strconv"
)

func UserAuthentication() gin.HandlerFunc {
	return func(context *gin.Context) {
		//userId := context.Param("userId")
		//if userId == "" || userId == "undefined" {
		//	var err error
		//	userId, err = context.Cookie("user_id")
		//	if err != nil {
		//		context.JSON(200, gin.H{
		//			"message": "No user id, please login again",
		//		})
		//	}
		//}
		userId, err := context.Cookie("user_id")
		if err != nil {
			context.JSON(200, gin.H{
				"message": "No user id, please login again",
			})
		}
		//documentId := context.Param("documentId")
		sessionToken, err := context.Cookie("session_token")
		if err != nil {
			context.JSON(200, gin.H{
				"message": "No session token, please login again",
			})
			return
		}
		trueUserId, err := handlers.RedisClient.Get(sessionToken).Result()
		if err != nil {
			context.JSON(200, gin.H{
				"message": "Invalid session token, please login again",
			})
			return
		}
		userIdNum, err := strconv.Atoi(userId)
		trueUserIdNum, err := strconv.Atoi(trueUserId)
		if err != nil {
			context.JSON(200, gin.H{
				"message": "Invalid user id or session token",
			})
			return
		}
		if userIdNum != trueUserIdNum {
			context.JSON(200, gin.H{
				"message": "Invalid user id or session token",
			})
			return
		}
		context.Set("userId", userIdNum)
		context.Set("session_token", sessionToken)
		context.Next()
	}
}
