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
		Name:        sharedconfig.String("SERVICE_NAME", "notification-service"),
		Port:        sharedconfig.String("NOTIFICATION_SERVICE_PORT", sharedconfig.String("PORT", "5007")),
		GinMode:     sharedconfig.String("GIN_MODE", "debug"),
		DatabaseURL: sharedconfig.String("NOTIFICATION_DATABASE_URL", ""),
	}
}
