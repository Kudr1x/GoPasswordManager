package storage

import (
	"encoding/json"
	"fmt"
	"goPasswordManager/internal/models"
	"os"
	"path/filepath"
)

type StorageConfig struct {
	Version int               `json:"version"`
	Hash    []byte            `json:"hash"`
	Salt    []byte            `json:"salt"`
	Entries []models.Password `json:"entries"`
}

func AddEntry(name string, entry models.Password) error {
	cfg, err := LoadConfig(name)
	if err != nil {
		return fmt.Errorf("failed to load storage: %v", err)
	}

	cfg.Entries = append(cfg.Entries, entry)

	return SaveConfig(name, cfg)
}

func SaveConfig(name string, cfg *StorageConfig) error {
	path, err := getStoragePath(name)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

func getStoragePath(name string) (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "pm", name+".json"), nil
}

func Exists(name string) bool {
	path, err := getStoragePath(name)
	if err != nil {
		return false
	}
	_, err = os.Stat(path)
	return !os.IsNotExist(err)
}

func Create(name string, hash, salt []byte) error {
	config := StorageConfig{
		Version: 1,
		Hash:    hash,
		Salt:    salt,
	}

	path, err := getStoragePath(name)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

func LoadConfig(name string) (*StorageConfig, error) {
	path, err := getStoragePath(name)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config StorageConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
