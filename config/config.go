package config

import (
	"os"

	"github.com/deveusss/chronio-core/encryption"

	"github.com/caarlos0/env/v10"
)

type AppConfig struct {
	Name    string `env:"APP_NAME"`
	Version string `env:"APP_VERSION"`
	Debug   bool   `env:"APP_DEBUG"`
	Env     string `env:"APP_ENV"`
}

type PostgresConfig struct {
	Host     string `env:"POSTGRES_HOST"`
	Port     int    `env:"POSTGRES_PORT"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	Name     string `env:"POSTGRES_DB"`
}

var config *Configuration

type ServerConfig struct {
	Host                string `env:"SERVER_HOST"`
	Port                int    `env:"SERVER_PORT"`
	RateLimit           int    `env:"SERVER_RATE_LIMIT"`
	RateLimitExpiration int    `env:"SERVER_RATE_LIMIT_EXPIRATION"`
}

type AuthConfig struct {
	JwtExpiration      int    `env:"JWT_EXPIRATION"`
	GoogleClientId     string `env:"GOOGLE_AUTH_CLIENT_ID"`
	GoogleClientSecret string `env:"GOOGLE_AUTH_CLIENT_SECRET"`
}

func (auth *AuthConfig) JwtSecret() encryption.ISecureString {
	return encryption.NewSecureString(os.Getenv("JWT_SECRET"))
}

type Configuration struct {
	App      AppConfig      `env:"-"`
	Database PostgresConfig `env:"-"`
	Server   ServerConfig   `env:"-"`
	Auth     AuthConfig     `env:"-"`
}

func load() *Configuration {
	config = &Configuration{}
	if err := env.Parse(config); err != nil {
		panic(err)
	}
	return config
}

func init() {
	load()
}

func Config() *Configuration {
	if config == nil {
		panic("config not loaded")
	}
	return config
}
