package config

import (
	"github.com/magiconair/properties"
)

type ServerConfig struct {
	Port string
}

func InitServerConfig() *ServerConfig {
	configProperties := properties.MustLoadFile("config/config.properties", properties.UTF8)
	port := configProperties.GetString("server.port", "8080")
	return &ServerConfig{
		Port: port,
	}
}
