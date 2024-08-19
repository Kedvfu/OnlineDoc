package handlers

import "github.com/gin-gonic/gin"

func GetUserInfo(context *gin.Context) {
	RedisClient.Get("user_")

}
