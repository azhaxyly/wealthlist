package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Port string
	}
	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		DbName   string
	}
	Mail struct {
		Host     string
		Port     string
		User     string
		Password string
	}
	Google struct {
		SpreadsheetId     string
		SheetName         string
		PathToCredentials string
	}
}

func InitConfig(envPath string) (*Config, error) {
	if err := godotenv.Load(envPath); err != nil {
		return nil, err
	}

	cfg := &Config{}
	cfg.Server.Port = os.Getenv("SERVER_PORT")

	cfg.Database.Host = os.Getenv("DB_HOST")
	cfg.Database.Port, _ = strconv.Atoi(os.Getenv("DB_PORT"))
	cfg.Database.User = os.Getenv("DB_USER")
	cfg.Database.Password = os.Getenv("DB_PASSWORD")
	cfg.Database.DbName = os.Getenv("DB_NAME")

	cfg.Mail.Host = os.Getenv("MAIL_HOST")
	cfg.Mail.Port = os.Getenv("MAIL_PORT")
	cfg.Mail.User = os.Getenv("MAIL_USER")
	cfg.Mail.Password = os.Getenv("MAIL_PASSWORD")

	cfg.Google.SpreadsheetId = os.Getenv("GOOGLE_SPREADSHEET_ID")
	cfg.Google.SheetName = os.Getenv("GOOGLE_SHEET_NAME")
	cfg.Google.PathToCredentials = os.Getenv("GOOGLE_CREDENTIALS_PATH")

	return cfg, nil
}
