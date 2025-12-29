package config

import (
	"encoding/json"
	"os"

	"Notification-Service/internal/sender"
)

type Config struct {
	ServerPort string             `json:"server_port"`
	DBDSN      string             `json:"db_dsn"`
	Email      sender.EmailConfig `json:"email_config"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
