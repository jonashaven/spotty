package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	ClientID     string    `json:"client_id"`
	ClientSecret string    `json:"client_secret"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Expiry       time.Time `json:"expiry"`
}

func configDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "spotify-cli")
}

func configPath() string {
	return filepath.Join(configDir(), "config.json")
}

func cachePath() string {
	return filepath.Join(configDir(), "cache.json")
}

type Cache struct {
	Text      string    `json:"text"`
	IsPlaying bool      `json:"is_playing"`
	FetchedAt time.Time `json:"fetched_at"`
}

func LoadCache() (*Cache, error) {
	data, err := os.ReadFile(cachePath())
	if err != nil {
		return nil, err
	}
	var c Cache
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func SaveCache(c *Cache) error {
	p := cachePath()
	os.MkdirAll(filepath.Dir(p), 0700)
	data, _ := json.Marshal(c)
	return os.WriteFile(p, data, 0600)
}

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile(configPath())
	if err != nil {
		return &Config{}, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return &Config{}, err
	}
	return &cfg, nil
}

func SaveConfig(cfg *Config) error {
	p := configPath()
	if err := os.MkdirAll(filepath.Dir(p), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0600)
}
