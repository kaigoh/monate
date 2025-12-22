package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	defaultPort                 = "8080"
	defaultDatabaseDriver       = "sqlite"
	defaultDatabaseDSN          = "monate.db"
	defaultManualCheckInterval  = time.Minute
	defaultManualCheckBatchSize = 25
)

// Config represents all runtime configuration Monate depends on.
type Config struct {
	Port                string
	DatabaseDriver      string
	DatabaseDSN         string
	PublicBaseURL       string
	ManualCheckInterval time.Duration
	ManualCheckBatch    int
}

// Load reads environment variables and produces a Config with sane defaults.
func Load() (*Config, error) {
	cfg := &Config{
		Port:                getEnv("PORT", defaultPort),
		DatabaseDriver:      getEnv("MONATE_DB_DRIVER", defaultDatabaseDriver),
		DatabaseDSN:         getEnv("MONATE_DB_DSN", defaultDatabaseDSN),
		PublicBaseURL:       os.Getenv("MONATE_PUBLIC_BASE_URL"),
		ManualCheckInterval: durationEnv("MONATE_MANUAL_CHECK_INTERVAL", defaultManualCheckInterval),
		ManualCheckBatch:    intEnv("MONATE_MANUAL_CHECK_BATCH", defaultManualCheckBatchSize),
	}

	if cfg.PublicBaseURL == "" {
		cfg.PublicBaseURL = fmt.Sprintf("http://localhost:%s", cfg.Port)
	}

	if cfg.ManualCheckBatch <= 0 {
		cfg.ManualCheckBatch = defaultManualCheckBatchSize
	}

	return cfg, nil
}

func getEnv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}

func durationEnv(key string, def time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		parsed, err := time.ParseDuration(val)
		if err == nil {
			return parsed
		}
	}
	return def
}

func intEnv(key string, def int) int {
	if val := os.Getenv(key); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			return parsed
		}
	}
	return def
}
