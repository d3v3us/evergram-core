package database

import "time"

type DatabaseConfig struct {
	Host           int    `yaml:"host"`
	Port           int    `yaml:"port"`
	User           string `yaml:"user"`
	Password       string `yaml:"password"`
	MigrationsPath string
	Name           string        `yaml:"name"`
	CacheTTL       time.Duration `yaml:"cache_ttl" env-default:"5m"`
	MaxRetries     int           `yaml:"max_retries" env-default:"3"`
	RetryWait      time.Duration `yaml:"retry_wait" env-default:"5s"`
}
