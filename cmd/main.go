package main

import (
	"OnlineDoc/api/handlers"
	"OnlineDoc/api/sessions"

	"OnlineDoc/api/routes"
	"OnlineDoc/config"
	"OnlineDoc/database"
	"OnlineDoc/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func main() {
	databaseConfig, databaseLogConfig := config.InitDatabaseConfig()
	db, err := gorm.Open(mysql.Open(databaseConfig), databaseLogConfig)
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err.Error())
	}
	database.InitialDatabase(db)

	Initialize()
	//authenticate.InitialDatabase()

	serverConfig := config.InitServerConfig()
	router := gin.Default()
	router.LoadHTMLGlob("web/templates/*")
	routes.RegisterRoutes(router)
	err = router.Run(":" + serverConfig.Port)
	if err != nil {
		log.Fatalf("failed to start server: %s", err.Error())
	}

}
func Initialize() {
	models.InitializeModels()
	handlers.InitializeRedis()
	sessions.InitialExcelSessions()
}
