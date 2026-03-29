package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	DefaultDirName  = "dashcli"
	DefaultFileName = "config.json"
)

type File struct {
	DefaultProfile string             `json:"default_profile"`
	Profiles       map[string]Profile `json:"profiles"`
}

type Profile struct {
	BaseURL   string `json:"base_url"`
	APIKeyEnv string `json:"api_key_env"`
}

func DefaultPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("resolve user config dir: %w", err)
	}
	return filepath.Join(configDir, DefaultDirName, DefaultFileName), nil
}

func LoadDefault() (File, error) {
	path, err := DefaultPath()
	if err != nil {
		return File{}, err
	}
	return Load(path)
}

func Load(path string) (File, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return File{}, nil
		}
		return File{}, fmt.Errorf("read config file %q: %w", path, err)
	}

	var parsed File
	if err := json.Unmarshal(content, &parsed); err != nil {
		return File{}, fmt.Errorf("parse config file %q: %w", path, err)
	}
	if parsed.Profiles == nil {
		parsed.Profiles = map[string]Profile{}
	}
	return parsed, nil
}

