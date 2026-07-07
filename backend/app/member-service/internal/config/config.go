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
		Name:        sharedconfig.String("SERVICE_NAME", "member-service"),
		Port:        sharedconfig.String("MEMBER_SERVICE_PORT", sharedconfig.String("PORT", "5004")),
		GinMode:     sharedconfig.String("GIN_MODE", "debug"),
		DatabaseURL: sharedconfig.String("MEMBER_DATABASE_URL", ""),
	}
}
