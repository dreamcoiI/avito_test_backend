package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	Port   string
	DBUser string
	DBPass string
	DBHost string
	DBPort int
	DBName string
}

func LoadConfigFromEnv(envFilePath string) Config {
	err := godotenv.Load(envFilePath)
	if err != nil {
		fmt.Println("Error loading end file")
	}

	var config Config

	config.Port = getEnvOrDefault("PORT", "8080")
	config.DBUser = getEnvOrDefault("DBUSER", "avito_tech")
	config.DBPass = getEnvOrDefault("DBPASS", "Avito")
	config.DBHost = getEnvOrDefault("DBHOST", "localhost")
	config.DBPort = getEnvOrDefaultInt("DBPORT", 5432)
	config.DBName = getEnvOrDefault("DBNAME", "test") //TODO created .env file

	return config
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvOrDefaultInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	valueInt, err := strconv.Atoi(valueStr)
	if err != nil {
		fmt.Printf("Error parsing %s as integer: %s\n", key, err)
		return defaultValue
	}
	return valueInt
}

func (config *Config) GetDBString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.DBUser, config.DBPass, config.DBHost, config.DBPort, config.DBName)
}
