package main

import (
	"os"
	"strconv"
)

type Config struct {
	Port                   string
	ShutdownTimeoutSeconds int
	RequestTimeoutSeconds  int
}

func LoadConfig() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	shutdown, _ := strconv.Atoi(os.Getenv("SHUTDOWN_TIMEOUT_SECONDS"))
	if shutdown < 10 {
		shutdown = 10
	}
	requestTimeout, _ := strconv.Atoi(os.Getenv("REQUEST_TIMEOUT"))
	if requestTimeout <= 0 {
		requestTimeout = 5
	}
	return Config{Port: port, ShutdownTimeoutSeconds: shutdown, RequestTimeoutSeconds: requestTimeout}
}
