package common

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Env            string
	DBHost         string
	DBPort         int
	DBUser         string
	DBPassword     string
	DBName         string
	GHToken        string
	ParserInterval int
}

func LoadConfig() (*Config, error) {
    _ = godotenv.Load()

    getRequiredEnv := func(key string) (string, error) {
        value, exists := os.LookupEnv(key)
        if !exists || value == "" {
            return "", fmt.Errorf("required environment variable %s is not set", key)
        }
        return value, nil
    }

    getEnvWithDefault := func(key, defaultValue string) string {
        if value, exists := os.LookupEnv(key); exists && value != "" {
            return value
        }
        return defaultValue
    }

    // Validate and get required values
    dbHost, err := getRequiredEnv("DB_HOST")
    if err != nil {
        return nil, err
    }

    dbUser, err := getRequiredEnv("DB_USER")
    if err != nil {
        return nil, err
    }

    dbPassword, err := getRequiredEnv("DB_PASSWORD")
    if err != nil {
        return nil, err
    }

    dbName, err := getRequiredEnv("DB_NAME")
    if err != nil {
        return nil, err
    }

    ghToken, err := getRequiredEnv("GH_TOKEN")
    if err != nil {
        return nil, err
    }

    dbPortStr, err := getRequiredEnv("DB_PORT")
    if err != nil {
        return nil, err
    }
    dbPort, err := strconv.Atoi(dbPortStr)
    if err != nil {
        return nil, fmt.Errorf("invalid DB_PORT: %v", err)
    }

    parserIntervalStr := getEnvWithDefault("PARSER_INTERVAL", "30")
    parserInterval, err := strconv.Atoi(parserIntervalStr)
    if err != nil {
        return nil, fmt.Errorf("invalid PARSER_INTERVAL: %v", err)
    }

    env := getEnvWithDefault("ENV", "prod")

    return &Config{
        Env:            env,
        DBHost:         dbHost,
        DBPort:         dbPort,
        DBUser:         dbUser,
        DBPassword:     dbPassword,
        DBName:         dbName,
        GHToken:        ghToken,
        ParserInterval: parserInterval,
    }, nil
}
