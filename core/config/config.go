package config

import (
	"effect-mobile/pkg/logger"
	"effect-mobile/pkg/utils"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ApiPort int
	DSN     string
}

var config *Config

func GetConfig() *Config {
	return config
}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		logger.GetLogger().Error(".env file not found!")
	}

	port := utils.ParseStringToIntOrDefault(os.Getenv("API_PORT"), 3016)

	config = &Config{
		ApiPort: port,
		DSN:     os.Getenv("POSTGRES_URL"),
	}

}
