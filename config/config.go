package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DBConfig
	SMTP     SMTPConfig
	Google   GoogleConfig
}

type ServerConfig struct {
	Host string
	Port string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

type GoogleConfig struct {
	SpreadsheetID     string
	SheetName         string
	PathToCredentials string
}

func InitConfig(envPath string) (*Config, error) {
	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	cfg := &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "MILLIONAIRE"),
		},
		SMTP: SMTPConfig{
			Host:     getEnv("MAIL_HOST", "smtp.gmail.com"),
			Port:     getEnv("MAIL_PORT", "587"),
			Username: getEnv("MAIL_USER", ""),
			Password: getEnv("MAIL_PASSWORD", ""),
			From:     getEnv("MAIL_FROM", ""),
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
