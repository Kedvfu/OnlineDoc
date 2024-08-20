package config

import (
	"github.com/magiconair/properties"
)

var AllowRegistration bool

func InitializeApplicationConfig() {
	applicationProperties := properties.MustLoadFile("config/config.properties", properties.UTF8)
	AllowRegistration = applicationProperties.GetBool("Application.registrationsEnabled", true)
}
