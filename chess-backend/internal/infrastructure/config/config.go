package config

import (
	"os"
	"strconv"

	"chess-backend/internal/constants/appconst"
)

type Config struct {
	Port          string
	DBHost        string
	DBPort        int
	DBUser        string
	DBPassword    string
	DBName        string
	RedisAddr     string
	RedisPassword string
	JWTSecret     string
}

func Load() *Config {
	port, _ := strconv.Atoi(getEnv("DB_PORT", strconv.Itoa(appconst.DefaultDBPort)))
	return &Config{
		Port:          getEnv("PORT", appconst.DefaultPort),
		DBHost:        getEnv("DB_HOST", appconst.DefaultDBHost),
		DBPort:        port,
		DBUser:        getEnv("DB_USER", appconst.DefaultDBUser),
		DBPassword:    getEnv("DB_PASSWORD", appconst.DefaultDBPassword),
		DBName:        getEnv("DB_NAME", appconst.DefaultDBName),
		RedisAddr:     getEnv("REDIS_ADDR", appconst.DefaultRedisAddr),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		JWTSecret:     getEnv("JWT_SECRET", "super-secret-key"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
