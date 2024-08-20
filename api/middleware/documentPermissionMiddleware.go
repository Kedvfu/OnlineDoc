package middleware

import (
	"OnlineDoc/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

func DocumentPermissionMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		userId, _ := context.Get("userId")
		userIdNum, _ := userId.(int)
		documentId := context.Param("documentId")
		documentIdNum, err := strconv.Atoi(documentId)
		if err != nil {
			context.JSON(200, gin.H{
				"message": "unable to parse documentId to int",
			})
			context.Abort()
			return
		}
		permissionType, err := models.GetPermissionTypeByDocumentIdAndUserId(documentIdNum, userIdNum)
		if err != nil {
			context.JSON(200, gin.H{
				"message": "No permission found for this document",
			})
			context.Abort()
			return
		}
		if permissionType == 0 {
			context.JSON(400, gin.H{
				"permissionType": 0,
				"message":        "read only permission for this document",
			})
			context.Abort()
			return
		}
		if permissionType == 1 {
			context.Set("documentId", documentIdNum)
			context.Next()

			return
		}

	}

}
