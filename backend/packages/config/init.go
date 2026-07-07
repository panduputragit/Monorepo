package config

import "log"

func Init(paths ...string) error {
	if len(paths) == 0 {
		paths = []string{".env", "../../.env", "../../../.env"}
	}
	return LoadEnv(paths...)
}

// MustInit calls Init and panics on any error.
func MustInit(paths ...string) {
	if err := Init(paths...); err != nil {
		log.Fatalf("load config: %v", err)
	}
}
