package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Env      string
	Server   ServerConfig
	Database DBConfig
	SMTP     SMTPConfig
	Google   GoogleConfig
}

type ServerConfig struct {
	Host string
	Port int
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	To       string
}

type GoogleConfig struct {
	SpreadsheetID     string
	SheetName         string
	PathToCredentials string
}

func InitConfig(envPath string) (*Config, error) {
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Warning: .env file not found, using default values")
	}

	serverPort, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	if err != nil {
		log.Fatalf("Invalid SERVER_PORT value: %v", err)
	}

	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		log.Fatalf("Invalid DB_PORT value: %v", err)
	}

	mailPort, err := strconv.Atoi(getEnv("MAIL_PORT", "587"))
	if err != nil {
		log.Fatalf("Invalid MAIL_PORT value: %v", err)
	}

	cfg := &Config{
		Env: getEnv("APP_ENV", "local"),

		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: serverPort,
		},
		Database: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     dbPort,
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "12345678"),
			DBName:   getEnv("DB_NAME", "MILLIONAIRE"),
		},
		SMTP: SMTPConfig{
			Host:     getEnv("MAIL_HOST", "smtp.gmail.com"),
			Port:     mailPort,
			Username: getEnv("MAIL_USER", ""),
			Password: getEnv("MAIL_PASSWORD", ""),
			From:     getEnv("MAIL_FROM", ""),
			To:       getEnv("MAIL_TO", ""),
		},
		Google: GoogleConfig{
			SpreadsheetID:     getEnv("GOOGLE_SPREADSHEET_ID", ""),
			SheetName:         getEnv("GOOGLE_SHEET_NAME", ""),
			PathToCredentials: getEnv("GOOGLE_CREDENTIALS_PATH", ""),
		},
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
