package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	PORT string `env:"PORT"`
	DSN  string `env:"DSN"`
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func GetEnvConfig() *Config {
	return &Config{
		PORT: os.Getenv("PORT"),
		DSN:  os.Getenv("DSN"),
	}
}
