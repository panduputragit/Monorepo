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
		Name:        sharedconfig.String("SERVICE_NAME", "auth-service"),
		Port:        sharedconfig.String("AUTH_SERVICE_PORT", sharedconfig.String("PORT", "5001")),
		GinMode:     sharedconfig.String("GIN_MODE", "debug"),
		DatabaseURL: sharedconfig.String("AUTH_DATABASE_URL", ""),
	}
}
