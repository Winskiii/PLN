package config

import (
	"time"
)

type Config struct {
	AppName            string        `mapstructure:"APP_NAME"`
	Env                string        `mapstructure:"APP_ENV"`
	Host               string        `mapstructure:"APP_HOST"`
	Port               int           `mapstructure:"APP_PORT"`
	LogLevel           string        `mapstructure:"LOG_LEVEL"`
	AllowOrigins       []string      `mapstructure:"ALLOW_ORIGINS"`
	RateLimitRequests  int           `mapstructure:"RATE_LIMIT_REQUESTS"`
	RateLimitWindow    time.Duration `mapstructure:"RATE_LIMIT_WINDOW"`
	ReadTimeout        time.Duration `mapstructure:"READ_TIMEOUT"`
	WriteTimeout       time.Duration `mapstructure:"WRITE_TIMEOUT"`
	IdleTimeout        time.Duration `mapstructure:"IDLE_TIMEOUT"`
	MaxRequestBodySize int64         `mapstructure:"REQUEST_MAX_BODY_BYTES"`

	// Auth/JWT (for demo purposes; replace with real IDP in production)
	JWTSecret string        `mapstructure:"JWT_SECRET"`
	JWTTTL    time.Duration `mapstructure:"JWT_TTL"`
	DemoUser  string        `mapstructure:"AUTH_DEMO_USERNAME"`
	DemoPass  string        `mapstructure:"AUTH_DEMO_PASSWORD"`
}
