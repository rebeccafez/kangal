package config

import (
	"fmt"
	"log"
	"os"
)

func mustEnv(key string) string {
	v := os.Getenv(key)

	if v == "" {
		log.Fatalf("Required environment variable %s is not set", key)
	}

	return v
}

func getEnv(key string, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}

func envInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}

	var n int
	if _, err := fmt.Sscan(v, &n); err != nil {
		return fallback
	}
	return n
}

func envFloat(key string, fallback float64) float64 {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}

	var n float64
	if _, err := fmt.Sscan(v, &n); err != nil {
		return fallback
	}
	return n
}
