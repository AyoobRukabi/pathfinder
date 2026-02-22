package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Env            string `json:"env"`
	NetworkMapPath string `json:"-"`
	StartStation   string `json:"-"`
	EndStation     string `json:"-"`
	NumTrains      int    `json:"-"`
}

func MustLoad(mapPath string, start string, destination string, trains int) *Config {
	var cfg Config

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config/local/local.json"
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("can't read config file: %s, error: %v", configPath, err)
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("can't unmarshal config file: %s, error: %v", configPath, err)
	}

	cfg.NetworkMapPath = mapPath
	cfg.StartStation = start
	cfg.EndStation = destination
	cfg.NumTrains = trains

	return &cfg
}
