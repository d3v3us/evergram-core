package config

import (
	"flag"
	"os"
	"time"

	encryption "github.com/deveusss/evergram-core/encryption"

	"github.com/ilyakaznacheev/cleanenv"
)

type ConfigurationBase[Configuration any] struct {
	Config *Configuration
}
type DatabaseConfig struct {
	Host           string `yaml:"host"`
	Port           int    `yaml:"port"`
	User           string `yaml:"user"`
	Password       string `yaml:"password"`
	MigrationsPath string
	Name           string        `yaml:"name"`
	CacheTTL       time.Duration `yaml:"cache_ttl" env-default:"5m"`
	MaxRetries     int           `yaml:"max_retries" env-default:"3"`
	RetryWait      time.Duration `yaml:"retry_wait" env-default:"5s"`
	MaxIdleConns   int           `yaml:"max_idle_conns" env-default:"5"`
	MaxOpenConns   int           `yaml:"max_open_conns" env-default:"10"`
}

type AppConfig struct {
	Env        string         `yaml:"env" env-default:"local"`
	GRPC       GRPCConfig     `yaml:"grpc"`
	DbConfig   DatabaseConfig `yaml:"db"`
	AuthConfig AuthConfig     `yaml:"auth"`
}
type JwtConfig struct {
	TokenTTL time.Duration `yaml:"token_ttl" env-default:"1h"`
	Secret   string        `yaml:"secret"`
}

func (c *JwtConfig) GetSecret() encryption.ISecureString {
	return encryption.NewSecureString(c.Secret)
}

type ExternalAuthConfig struct {
	Google GoogleAuthConfig `yaml:"google"`
}
type GoogleAuthConfig struct {
	GoogleClientId     string `yaml:"google_client_id"`
	GoogleClientSecret string `yaml:"google_client_secret"`
}
type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}
type AuthConfig struct {
	ExternalAuthConfig ExternalAuthConfig `yaml:"external"`
	Jwt                JwtConfig          `yaml:"jwt"`
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
		Config: &cfg,
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
