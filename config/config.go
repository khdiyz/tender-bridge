package config

import (
	"os"
	"sync"

	"github.com/spf13/cast"
)

var (
	instance *Config
	once     sync.Once
)

type Config struct {
	HTTPHost string
	HTTPPort int

	Environment string
	Debug       bool

	PostgresHost     string
	PostgresPort     int
	PostgresDatabase string
	PostgresUser     string
	PostgresPassword string

	RedisHost     string
	RedisPort     int
	RedisPassword string
	RedisDB       int

	JWTSecret                string
	JWTAccessExpirationHours int
	JWTRefreshExpirationDays int

	HashKey string
}

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{
			HTTPHost:    cast.ToString(getOrReturnDefault("HOST", "localhost")),
			HTTPPort:    cast.ToInt(getOrReturnDefault("PORT", 8888)),
			Environment: cast.ToString(getOrReturnDefault("ENVIRONMENT", "development")),
			Debug:       cast.ToBool(getOrReturnDefault("DEBUG", true)),

			PostgresHost:     cast.ToString(getOrReturnDefault("POSTGRES_HOST", "db")),
			PostgresPort:     cast.ToInt(getOrReturnDefault("POSTGRES_PORT", 5432)),
			PostgresDatabase: cast.ToString(getOrReturnDefault("POSTGRES_DB", "tender_bridge_db")),
			PostgresUser:     cast.ToString(getOrReturnDefault("POSTGRES_USER", "postgres")),
			PostgresPassword: cast.ToString(getOrReturnDefault("POSTGRES_PASSWORD", "password")),

			RedisHost:     cast.ToString(getOrReturnDefault("REDIS_HOST", "redis")),
			RedisPort:     cast.ToInt(getOrReturnDefault("REDIS_PORT", 6379)),
			RedisPassword: cast.ToString(getOrReturnDefault("REDIS_PASSWORD", "")),
			RedisDB:       cast.ToInt(getOrReturnDefault("REDIS_DB", 0)),

			JWTSecret:                cast.ToString(getOrReturnDefault("JWT_SECRET", "tender-bridge-forever")),
			JWTAccessExpirationHours: cast.ToInt(getOrReturnDefault("JWT_ACCESS_EXPIRATION_HOURS", 12)),
			JWTRefreshExpirationDays: cast.ToInt(getOrReturnDefault("JWT_REFRESH_EXPIRATION_DAYS", 3)),

			HashKey: cast.ToString(getOrReturnDefault("HASH_KEY", "skd32r8wdahHSdqw")),
		}
	})

	return instance
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultValue
}
