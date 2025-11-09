package config

// Re-export config from root config package
import (
	rootConfig "casdoor-casbin-openbao/config"
)

// Config re-exports the Config type from root config package
type Config = rootConfig.Config

var AppConfig *Config

func Init() {
	AppConfig = rootConfig.LoadConfig()
}

func GetConfig() *Config {
	return AppConfig
}

