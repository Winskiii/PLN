package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

func Load() (*Config, error) {
	// Load .env into environment for local development (ignore if missing)
	_ = gotenv.Load()

	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Defaults
	v.SetDefault("APP_NAME", "secure-api")
	v.SetDefault("APP_ENV", "development")
	v.SetDefault("APP_HOST", "0.0.0.0")
	v.SetDefault("APP_PORT", 8080)
	v.SetDefault("LOG_LEVEL", "info")
	v.SetDefault("ALLOW_ORIGINS", []string{"http://localhost:3001"})
	v.SetDefault("RATE_LIMIT_REQUESTS", 100)
	v.SetDefault("RATE_LIMIT_WINDOW", "1m")
	v.SetDefault("READ_TIMEOUT", "10s")
	v.SetDefault("WRITE_TIMEOUT", "10s")
	v.SetDefault("IDLE_TIMEOUT", "60s")
	v.SetDefault("REQUEST_MAX_BODY_BYTES", 1<<20)
	v.SetDefault("JWT_SECRET", "change-me-please")
	v.SetDefault("JWT_TTL", "15m")
	v.SetDefault("AUTH_DEMO_USERNAME", "admin")
	v.SetDefault("AUTH_DEMO_PASSWORD", "password123")

	// Optionally read .env as config file if present (best effort)
	v.SetConfigFile(".env")
	_ = v.ReadInConfig()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	// Parse durations that may come as string
	if d, err := time.ParseDuration(v.GetString("RATE_LIMIT_WINDOW")); err == nil {
		cfg.RateLimitWindow = d
	}
	if d, err := time.ParseDuration(v.GetString("READ_TIMEOUT")); err == nil {
		cfg.ReadTimeout = d
	}
	if d, err := time.ParseDuration(v.GetString("WRITE_TIMEOUT")); err == nil {
		cfg.WriteTimeout = d
	}
	if d, err := time.ParseDuration(v.GetString("IDLE_TIMEOUT")); err == nil {
		cfg.IdleTimeout = d
	}
	if d, err := time.ParseDuration(v.GetString("JWT_TTL")); err == nil {
		cfg.JWTTTL = d
	}

	// Parse allow origins when provided as comma-separated string in env
	if len(cfg.AllowOrigins) == 0 {
		raw := v.GetString("ALLOW_ORIGINS")
		if raw != "" {
			parts := strings.Split(raw, ",")
			cleaned := make([]string, 0, len(parts))
			for _, p := range parts {
				s := strings.TrimSpace(p)
				if s != "" {
					cleaned = append(cleaned, s)
				}
			}
			if len(cleaned) > 0 {
				cfg.AllowOrigins = cleaned
			}
		}
	}

	return &cfg, nil
}
