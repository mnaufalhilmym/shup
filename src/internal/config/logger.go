package config

import (
	"time"

	"github.com/rs/zerolog"
)

func ConfigureLogger(level string) {
	zerolog.TimeFieldFormat = time.RFC3339

	if level != "" {
		logLevel, err := zerolog.ParseLevel(level)
		if err == nil {
			zerolog.SetGlobalLevel(logLevel)
		}
	}
}
