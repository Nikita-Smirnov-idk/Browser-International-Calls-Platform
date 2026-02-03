package config

import (
	"fmt"
	"os"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	VoIP     VoIPConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	Secret string
}

type VoIPConfig struct {
	Provider   string
	AccountSID string
	AuthToken  string
	APIKey     string
	FromNumber string
}

func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD", "postgres"),
			DBName:   getEnv("POSTGRES_DB", "postgres"),
			SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		},
		VoIP: VoIPConfig{
			Provider:   getEnv("VOIP_PROVIDER", "twilio"),
			AccountSID: getEnv("VOIP_ACCOUNT_SID", ""),
			AuthToken:  getEnv("VOIP_AUTH_TOKEN", ""),
			APIKey:     getEnv("VOIP_API_KEY", ""),
			FromNumber: getEnv("VOIP_FROM_NUMBER", ""),
		},
	}

	if cfg.JWT.Secret == "your-secret-key-change-in-production" {
		fmt.Println("WARNING: Using default JWT secret. Set JWT_SECRET environment variable in production.")
	}

	if cfg.VoIP.AccountSID == "" || cfg.VoIP.AuthToken == "" {
		fmt.Println("WARNING: VoIP credentials not set. WebRTC calls will not work. Set VOIP_ACCOUNT_SID and VOIP_AUTH_TOKEN.")
	}

	return cfg, nil
}

func (c *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

