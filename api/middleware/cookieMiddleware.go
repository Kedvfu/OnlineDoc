package middleware

import (
	"OnlineDoc/api/handlers"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CookieMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		cookie, err := context.Cookie("session_token")
		if errors.Is(err, http.ErrNoCookie) {
			context.Next()
			return
		}

		userId := handlers.RedisClient.Get(cookie)
		if userId.Err() != nil {
			context.Next()
			return
		} else {

			userId, err := userId.Result()
			// 将数据添加到上下文中
			if err == nil {
				context.Set("userId", userId)
				context.Set("session_token", cookie)
				// 继续处理请求
				context.Next()
			} else {
				context.Next()
				return
			}
		}
	}
}
