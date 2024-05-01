package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	DBURI     	string
	DBPORT		int
    DBNAME     	string
    TESTDBNAME 	string
}

func New() *Config {
	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
    if err != nil {
        log.Fatalf("Error parsing DB_PORT: %v", err)
    }
    return &Config{
        DBURI:     	getEnv("DBURI", ""),
        DBPORT: 	dbPort,
		DBNAME:     getEnv("DBNAME", "mydatabase"),
        TESTDBNAME: getEnv("TESTDBNAME", ""),
    }
}

// Helper function to read an environment variable or return a default value
func getEnv(key, fallback string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return fallback
}
