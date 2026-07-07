package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type MissingEnvError struct {
	Key string
}

func (e MissingEnvError) Error() string {
	return fmt.Sprintf("required environment variable %q is not set", e.Key)
}

func Bool(key string, fallback bool) bool {
	value := strings.TrimSpace(os.Getenv(key))

	if value == "" {
		return fallback
	}

	v, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}

	return v
}

func Int(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))

	if value == "" {
		return fallback
	}

	v, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return v
}

func Duration(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))

	if value == "" {
		return fallback
	}

	d, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}

	return d
}

func Required(key string) (string, error) {
	value := strings.TrimSpace(os.Getenv(key))

	if value == "" {
		return "", MissingEnvError{
			Key: key,
		}
	}

	return value, nil
}

func MustString(key string) string {
	value, err := Required(key)
	if err != nil {
		panic(err)
	}

	return value
}

func MustInt(key string) int {
	value := MustString(key)

	v, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}

	return v
}

func MustBool(key string) bool {
	value := MustString(key)

	v, err := strconv.ParseBool(value)
	if err != nil {
		panic(err)
	}

	return v
}

func String(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	return value
}
