package config

import (
	"flag"
	"os"
	"time"

	encryption "github.com/deveusss/evergram-core/encryption"

	"github.com/ilyakaznacheev/cleanenv"
)

type ConfigurationBase[Configuration any] struct {
	Config Configuration
}

type AuthConfig struct {
	JwtExpiration      int    `env:"JWT_EXPIRATION"`
	GoogleClientId     string `env:"GOOGLE_AUTH_CLIENT_ID"`
	GoogleClientSecret string `env:"GOOGLE_AUTH_CLIENT_SECRET"`
}
type PostgresConfig struct {
	Host     string `env:"POSTGRES_HOST"`
	Port     int    `env:"POSTGRES_PORT"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	Name     string `env:"POSTGRES_DB"`
}

type AppConfig struct {
	Env            string     `yaml:"env" env-default:"local"`
	GRPC           GRPCConfig `yaml:"grpc"`
	MigrationsPath string
	TokenTTL       time.Duration `yaml:"token_ttl" env-default:"1h"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func Load[Configuration any]() *ConfigurationBase[Configuration] {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return LoadFromPath[Configuration](configPath)
}

func LoadFromPath[Configuration any](configPath string) *ConfigurationBase[Configuration] {
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Configuration

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &ConfigurationBase[Configuration]{
		Config: cfg,
	}
}
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}

func (auth *AuthConfig) JwtSecret() encryption.ISecureString {
	return encryption.NewSecureString(os.Getenv("JWT_SECRET"))
}
