package config

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/99designs/keyring"
)

const (
	envBaseURL = "REDASH_BASE_URL"
	envAPIKey  = "REDASH_API_KEY"
)

type LookupEnvFunc func(string) (string, bool)
type GetSecretFunc func(string) (string, error)

type Flags struct {
	BaseURL string
	APIKey  string
	Profile string
	Timeout time.Duration
	Debug   bool
}

type ResolveInput struct {
	Flags     Flags
	Config    File
	LookupEnv LookupEnvFunc
	GetSecret GetSecretFunc
}

type Resolved struct {
	BaseURL string
	APIKey  string
	Profile string
	Timeout time.Duration
	Debug   bool
}

type RuntimeError struct {
	Err error
}

func (err *RuntimeError) Error() string {
	if err == nil || err.Err == nil {
		return "runtime error"
	}
	return err.Err.Error()
}

func (err *RuntimeError) Unwrap() error {
	if err == nil {
		return nil
	}
	return err.Err
}

func Resolve(input ResolveInput) (Resolved, error) {
	lookupEnv := input.LookupEnv
	if lookupEnv == nil {
		lookupEnv = func(string) (string, bool) { return "", false }
	}
	getSecret := input.GetSecret
	if getSecret == nil {
		getSecret = func(string) (string, error) {
			return "", keyring.ErrKeyNotFound
		}
	}

	selectedProfile, profile, err := resolveProfile(input.Config, strings.TrimSpace(input.Flags.Profile))
	if err != nil {
		return Resolved{}, err
	}

	baseURL := strings.TrimSpace(input.Flags.BaseURL)
	if baseURL == "" {
		if value, ok := lookupEnv(envBaseURL); ok {
			baseURL = strings.TrimSpace(value)
		}
	}
	if baseURL == "" {
		baseURL = strings.TrimSpace(profile.BaseURL)
	}

	apiKey := strings.TrimSpace(input.Flags.APIKey)
	if apiKey == "" {
		secretProfile := selectedProfile
		if secretProfile == "" {
			secretProfile = "default"
		}

		secretValue, err := getSecret(secretProfile)
		if err == nil {
			apiKey = strings.TrimSpace(secretValue)
		} else if !errors.Is(err, keyring.ErrKeyNotFound) {
			return Resolved{}, &RuntimeError{Err: fmt.Errorf("load API key from keyring: %w", err)}
		}
	}
	if apiKey == "" {
		keyEnv := strings.TrimSpace(profile.APIKeyEnv)
		if keyEnv != "" {
			if value, ok := lookupEnv(keyEnv); ok {
				apiKey = strings.TrimSpace(value)
			}
		}
	}
	if apiKey == "" {
		if value, ok := lookupEnv(envAPIKey); ok {
			apiKey = strings.TrimSpace(value)
		}
	}

	return Resolved{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Profile: selectedProfile,
		Timeout: input.Flags.Timeout,
		Debug:   input.Flags.Debug,
	}, nil
}

func resolveProfile(configFile File, flagProfile string) (string, Profile, error) {
	profiles := configFile.Profiles
	if profiles == nil {
		profiles = map[string]Profile{}
	}

	selectedProfile := ""
	if flagProfile != "" {
		selectedProfile = flagProfile
	} else if strings.TrimSpace(configFile.DefaultProfile) != "" {
		selectedProfile = strings.TrimSpace(configFile.DefaultProfile)
	} else if len(profiles) > 0 {
		selectedProfile = "default"
	}

	if selectedProfile == "" {
		return "", Profile{}, nil
	}
	profile, ok := profiles[selectedProfile]
	if !ok {
		return "", Profile{}, fmt.Errorf("profile %q not found", selectedProfile)
	}
	return selectedProfile, profile, nil
}
