package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Port        string `env:"PORT"`
	DSN         string `env:"DSN"`
	Secret      string `env:"JWT_SECRET"`
	AccessTime  string `env:"JWT_ACCESS_TIME"`
	RefreshTime string `env:"JWT_REFRESH_TIME"`
}

var cfg *Config

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	cfg = &Config{
		Port:        os.Getenv("PORT"),
		DSN:         os.Getenv("DSN"),
		Secret:      os.Getenv("JWT_SECRET"),
		AccessTime:  os.Getenv("JWT_ACCESS_TIME"),
		RefreshTime: os.Getenv("JWT_REFRESH_TIME"),
	}
}

func GetEnvConfig() *Config {
	return cfg
}
