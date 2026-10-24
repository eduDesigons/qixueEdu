package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Addr        string
	DBHost      string
	DBPort      int
	DBName      string
	DBUser      string
	DBPassword  string
	RedisAddr   string
	RedisPass   string
	RedisDB     int
	AdminJWT    string
	UserJWT     string
	JWTExpire   time.Duration
	AssetDomain string
	LogDir      string
	LogKeepDays int
}

func Load() (Config, error) {
	loadEnvFile()

	cfg := Config{
		Addr:        env("GO_API_ADDR", ":8787"),
		DBHost:      env("DB_HOST", "127.0.0.1"),
		DBPort:      envInt("DB_PORT", 3306),
		DBName:      os.Getenv("DB_DATABASE"),
		DBUser:      os.Getenv("DB_USERNAME"),
		DBPassword:  os.Getenv("DB_PASSWORD"),
		RedisAddr:   fmt.Sprintf("%s:%d", env("REDIS_HOST", "127.0.0.1"), envInt("REDIS_PORT", 6379)),
		RedisPass:   os.Getenv("REDIS_PASSWORD"),
		RedisDB:     envInt("REDIS_DATABASE", 0),
		AdminJWT:    os.Getenv("JWT_KEY"),
		UserJWT:     os.Getenv("JWT_WEAPP_KEY"),
		JWTExpire:   time.Duration(envInt("JWT_EXPIRE", 5*3600)) * time.Second,
		AssetDomain: os.Getenv("ASSET_DOMAIN"),
		LogDir:      env("LOG_DIR", "logs"),
		LogKeepDays: envInt("LOG_KEEP_DAYS", 7),
	}
	if cfg.LogKeepDays < 1 {
		cfg.LogKeepDays = 7
	}
	if cfg.DBName == "" || cfg.DBUser == "" {
		return Config{}, fmt.Errorf("缺少 DB_DATABASE 或 DB_USERNAME")
	}
	if cfg.AdminJWT == "" || cfg.UserJWT == "" {
		return Config{}, fmt.Errorf("缺少 JWT_KEY 或 JWT_WEAPP_KEY")
	}
	if cfg.JWTExpire <= 0 {
		return Config{}, fmt.Errorf("JWT_EXPIRE 无效")
	}
	return cfg, nil
}

func loadEnvFile() {
	_ = godotenv.Load(".env")
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}
