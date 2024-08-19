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
			UserApi.GET("/info", handlers.GetUserInfo)

			UserDocumentApi := UserApi.Group("/document/:documentId")
			{
				UserDocumentApi.POST("/save", handlers.SaveDocument)
				UserDocumentApi.GET("/get", handlers.GetDocument)
				UserDocumentApi.GET("/link", handlers.GetDocumentLink)
				ExcelApi := UserDocumentApi.Group("/excel")
				{
					ExcelApi.POST("/update", handlers.UpdateExcel)   //更新单元格
					ExcelApi.POST("/refresh", handlers.RefreshExcel) //获取更新的单元格
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
