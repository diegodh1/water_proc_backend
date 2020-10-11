package config

import (
	"os"
	str "strconv"
)

//Config struct
type Config struct {
	DB         *DBConfig
	SecretSeed string
}

//DBConfig struct
type DBConfig struct {
	Username string
	Password string
	Database string
	Port     int
	Host     string
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

// Simple helper function to convert string to int
func convert(value string) int {
	port, err := str.Atoi(value)
	if err != nil {
		return 5432
	}
	return port
}

//GetConfig returns a db configuration
func GetConfig() *Config {
	return &Config{
		DB: &DBConfig{
			Username: getEnv("DB_USER_NAME", "sa"),
			Password: getEnv("DB_PASS", "Cristiano1994"),
			Database: getEnv("DB_NAME", "Water"),
			Port:     convert(getEnv("DB_PORT", "1433")),
			Host:     getEnv("DB_HOST", "localhost"),
		},
		SecretSeed: getEnv("SECRET_KEY", "integrapps"),
	}
}
