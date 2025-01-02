package config

import "os"

type Config struct {
	Port string
}

func New() *Config {

	port, exists := os.LookupEnv("HTTP_ADDR")
	if !exists {
		port = "8080"
	}

	return &Config{
		Port: port,
	}
}
