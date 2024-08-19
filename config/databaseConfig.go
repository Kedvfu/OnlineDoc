package config

import (
	"fmt"
	"github.com/magiconair/properties"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

type DatabaseConfig struct {
	Database string
	Host     string
	Port     string
	Username string
	Password string
}

func InitDatabaseConfig() (string, *gorm.Config) {
	databaseProperties := properties.MustLoadFile("config/config.properties", properties.UTF8)
	database := databaseProperties.GetString("database.name", "OnlineDocDatabase")
	host := databaseProperties.GetString("database.host", "localhost")
	port := databaseProperties.GetString("database.port", "3306")
	username := databaseProperties.GetString("database.username", "root")
	password := databaseProperties.GetString("database.password", "password")

	databaseConfig := DatabaseConfig{
		Database: database,
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}

	newLogger := logger.New(
		log.New(log.Writer(), "\r[GORM] ", log.LstdFlags),
		logger.Config{

			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	return formatDatabaseUrl(&databaseConfig), &gorm.Config{
		Logger: newLogger,
	}
}

func formatDatabaseUrl(databaseConfig *DatabaseConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", databaseConfig.Username, databaseConfig.Password, databaseConfig.Host, databaseConfig.Port, databaseConfig.Database)
}
