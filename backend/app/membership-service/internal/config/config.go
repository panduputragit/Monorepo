package config

import sharedconfig "github.com/panduputragit/gym/backend/packages/config"

type Config struct {
	Name        string
	Port        string
	GinMode     string
	DatabaseURL string
}

func Load() Config {
	_ = sharedconfig.LoadEnv(".env", "../../.env", "../../../.env")

	return Config{
		Name:        sharedconfig.String("SERVICE_NAME", "membership-service"),
		Port:        sharedconfig.String("MEMBERSHIP_SERVICE_PORT", sharedconfig.String("PORT", "5005")),
		GinMode:     sharedconfig.String("GIN_MODE", "debug"),
		DatabaseURL: sharedconfig.String("MEMBERSHIP_DATABASE_URL", ""),
	}
}
