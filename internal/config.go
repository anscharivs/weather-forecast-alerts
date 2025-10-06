package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl  string
	APIKey string
}

func LoadConfig() Config {

	_ = godotenv.Load()

	// if err != nil {
	// 	panic("ENV error")
	// }

	return Config{
		DBUrl:  os.Getenv("DATABASE_URL"),
		APIKey: os.Getenv("OPENWEATHER_API_KEY"),
	}
}
