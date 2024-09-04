package routes

import (
	"OnlineDoc/api/handlers"
	"OnlineDoc/api/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(engine *gin.Engine) {
	engine.Use(middleware.CookieMiddleware())
	api := engine.Group("/api")
	{
		api.GET("/ping", handlers.Ping)

		UserApi := api.Group("/user/:userId")
		{
			UserApi.Use(middleware.UserAuthentication())

			UserApi.GET("/documents", handlers.GetUserDocuments)
			UserApi.GET("/info/:targetUserId", handlers.GetUserInfo)

			UserDocumentApi := UserApi.Group("/document/:documentId")
			{
				UserDocumentApi.POST("/permission/:targetUserId/:permissionType", handlers.UpdateUserPermissionType)
				UserDocumentSaveApi := UserDocumentApi.Group("/save")
				{
					UserDocumentSaveApi.Use(middleware.DocumentPermissionMiddleware())
					UserDocumentSaveApi.POST("/", handlers.SaveDocument)
				}
				UserDocumentDeleteApi := UserDocumentApi.Group("/delete")
				{
					UserDocumentDeleteApi.Use(middleware.DocumentPermissionMiddleware())
					UserDocumentDeleteApi.POST("/", handlers.DeleteDocument)
				}
				UserDocumentApi.GET("/get", handlers.GetDocument)

				UserDocumentApi.GET("/link", handlers.GetDocumentLink)

				ExcelApi := UserDocumentApi.Group("/excel")
				{
					ExcelApiUpdateApi := ExcelApi.Group("/update")
					{
						ExcelApiUpdateApi.Use(middleware.DocumentPermissionMiddleware())
						ExcelApiUpdateApi.POST("/", handlers.UpdateExcel)
					}
					ExcelApiDownloadApi := ExcelApi.Group("/download")
					{
						ExcelApiDownloadApi.Use(middleware.DocumentPermissionMiddleware())
						ExcelApiDownloadApi.GET("/", handlers.DownloadExcel)
					}

					ExcelApi.POST("/refresh", handlers.RefreshExcel)
				}

			}
		}

	}
	homepage := engine.Group("/home")
	{
		homepage.GET("/", handlers.ShowHomepage)
	}
	defaultPage := engine.Group("/")
	{

		defaultPage.GET("/", handlers.ShowDefaultPage)
	}
	assets := engine.Group("/web/assets")
	{
		assets.Static("/css", "./web/assets/css")
		assets.Static("/js", "./web/assets/js")
		assets.Static("/img", "./web/assets/img")
	}
	loginPage := engine.Group("/login")
	{
		loginPage.GET("/", handlers.ShowLoginPage)
		loginPage.POST("/", handlers.HandleLogin)
	}
	registerPage := engine.Group("/register")
	{
		registerPage.GET("/", handlers.ShowRegisterPage)
		registerPage.POST("/", handlers.HandleRegister)
	}
	logoutPage := engine.Group("/logout")
	{
		logoutPage.GET("/", handlers.ShowLogoutPage)
	}
	documentPage := engine.Group("/document/:documentId")
	{
		documentPage.Use(middleware.DocumentMiddleware())
		documentPage.GET("/:documentType", handlers.ShowDocumentPage)
		documentPage.GET("/", handlers.ShowDocumentPage)
	}
	fromShare := engine.Group("/share/:shareUrl")
	{
		fromShare.GET("/", handlers.ShowDocumentFromSharePage)
	}

}
