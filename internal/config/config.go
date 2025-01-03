package config

import "os"

type Config struct {
	Port string
}

func Read() *Config {

	port, exists := os.LookupEnv("SERVICE_PORT")
	if !exists {
		port = "8080"
	}

	return &Config{
		Port: port,
	}
}
