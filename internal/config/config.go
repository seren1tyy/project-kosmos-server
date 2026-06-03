package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Server   ServerCfg   `json:"server"`
	Database DatabaseCfg `json:"database"`
}

type ServerCfg struct {
	Port      int    `json:"port"`
	CryptoKey string `json:"crypto_key"`
}

type DatabaseCfg struct {
	DSN string `json:"dsn"`
}

func Load(path string) Config {
	// Дефолтные значения (безопасные для локальной разработки)
	cfg := Config{
		Server: ServerCfg{
			Port:      12345,
			CryptoKey: "12345678901234567890123456789012",
		},
		Database: DatabaseCfg{
			DSN: "root:password@tcp(127.0.0.1:3306)/kosmos_db?parseTime=true&loc=Local",
		},
	}

	// Читаем файл, если есть
	if data, err := os.ReadFile(path); err == nil {
		if err := json.Unmarshal(data, &cfg); err != nil {
			log.Fatalf("❌ Failed to parse config.json: %v", err)
		}
	} else {
		log.Printf("⚠️ %s not found, using defaults", path)
	}

	// Валидация
	if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
		log.Fatalf("❌ Invalid server port: %d", cfg.Server.Port)
	}
	if len(cfg.Server.CryptoKey) != 32 {
		log.Fatalf("❌ Crypto key must be exactly 32 chars. Got: %d", len(cfg.Server.CryptoKey))
	}

	return cfg
}
