package config

import (
	"os"
)

type Config struct {
	Server   ServerConfig
	Casdoor  CasdoorConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port string
	Host string
}

type CasdoorConfig struct {
	Endpoint     string
	ClientID     string
	ClientSecret string
	Organization string
	Application  string
	RedirectURL  string
	Certificate  string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "localhost"),
		},
		Casdoor: CasdoorConfig{
			Endpoint:     getEnv("CASDOOR_ENDPOINT", "http://localhost:8000"),
			ClientID:     getEnv("CASDOOR_CLIENT_ID", ""),
			ClientSecret: getEnv("CASDOOR_CLIENT_SECRET", ""),
			Organization: getEnv("CASDOOR_ORGANIZATION", "built-in"),
			Application:  getEnv("CASDOOR_APPLICATION", "app-built-in"),
			RedirectURL:  getEnv("CASDOOR_REDIRECT_URL", "http://localhost:8080/api/auth/callback"),
			Certificate:  getEnv("CASDOOR_CERTIFICATE", ""),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "casbin"),
			Password: getEnv("DB_PASSWORD", "casbinpw"),
			DBName:   getEnv("DB_NAME", "casdoor"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

