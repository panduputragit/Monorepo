package config

import sharedconfig "github.com/panduputragit/gym/backend/packages/config"

type Config struct {
	Name        string
	Port        string
	GinMode     string
	DatabaseURL string
}

func Load() Config {
	return Config{
		Name:        sharedconfig.String("SERVICE_NAME", "employee-service"),
		Port:        sharedconfig.String("EMPLOYEE_SERVICE_PORT", sharedconfig.String("PORT", "5002")),
		GinMode:     sharedconfig.String("GIN_MODE", "debug"),
		DatabaseURL: sharedconfig.String("EMPLOYEE_DATABASE_URL", ""),
	}
}
