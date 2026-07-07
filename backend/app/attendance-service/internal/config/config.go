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
		Name:        sharedconfig.String("SERVICE_NAME", "attendance-service"),
		Port:        sharedconfig.String("ATTENDANCE_SERVICE_PORT", sharedconfig.String("PORT", "5006")),
		GinMode:     sharedconfig.String("GIN_MODE", "debug"),
		DatabaseURL: sharedconfig.String("ATTENDANCE_DATABASE_URL", ""),
	}
}
