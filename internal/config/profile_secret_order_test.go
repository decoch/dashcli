package config

import (
	"errors"
	"testing"

	"github.com/99designs/keyring"
)

func TestResolve_APIKeyOrder_FlagWins(t *testing.T) {
	resolved, err := Resolve(ResolveInput{
		Flags: Flags{
			APIKey:  "flag-key",
			Profile: "prod",
		},
		Config: File{
			Profiles: map[string]Profile{
				"prod": {APIKeyEnv: "REDASH_API_KEY_PROD"},
			},
		},
		GetSecret: func(profile string) (string, error) {
			return "secret-key", nil
		},
		LookupEnv: func(key string) (string, bool) {
			return "env-key", true
		},
	})
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if resolved.APIKey != "flag-key" {
		t.Fatalf("Resolve().APIKey = %q, want %q", resolved.APIKey, "flag-key")
	}
}

func TestResolve_APIKeyOrder_SecretWinsOverEnv(t *testing.T) {
	resolved, err := Resolve(ResolveInput{
		Flags: Flags{Profile: "prod"},
		Config: File{
			Profiles: map[string]Profile{
				"prod": {APIKeyEnv: "REDASH_API_KEY_PROD"},
			},
		},
		GetSecret: func(profile string) (string, error) {
			if profile != "prod" {
				t.Fatalf("profile = %q, want %q", profile, "prod")
			}
			return "secret-key", nil
		},
		LookupEnv: func(key string) (string, bool) {
			return "env-key", true
		},
	})
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if resolved.APIKey != "secret-key" {
		t.Fatalf("Resolve().APIKey = %q, want %q", resolved.APIKey, "secret-key")
	}
}

func TestResolve_APIKeyOrder_ProfileEnvThenGlobalEnv(t *testing.T) {
	resolved, err := Resolve(ResolveInput{
		Flags: Flags{Profile: "prod"},
		Config: File{
			Profiles: map[string]Profile{
				"prod": {APIKeyEnv: "REDASH_API_KEY_PROD"},
			},
		},
		GetSecret: func(profile string) (string, error) {
			return "", keyring.ErrKeyNotFound
		},
		LookupEnv: func(key string) (string, bool) {
			switch key {
			case "REDASH_API_KEY_PROD":
				return "profile-env-key", true
			case "REDASH_API_KEY":
				return "global-env-key", true
			default:
				return "", false
			}
		},
	})
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if resolved.APIKey != "profile-env-key" {
		t.Fatalf("Resolve().APIKey = %q, want %q", resolved.APIKey, "profile-env-key")
	}
}

func TestResolve_APIKeyOrder_GlobalEnvFallback(t *testing.T) {
	resolved, err := Resolve(ResolveInput{
		Flags:  Flags{},
		Config: File{},
		GetSecret: func(profile string) (string, error) {
			if profile != "default" {
				t.Fatalf("profile = %q, want %q", profile, "default")
			}
			return "", keyring.ErrKeyNotFound
		},
		LookupEnv: func(key string) (string, bool) {
			if key == "REDASH_API_KEY" {
				return "global-env-key", true
			}
			return "", false
		},
	})
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if resolved.APIKey != "global-env-key" {
		t.Fatalf("Resolve().APIKey = %q, want %q", resolved.APIKey, "global-env-key")
	}
}

func TestResolve_SecretRuntimeError(t *testing.T) {
	_, err := Resolve(ResolveInput{
		Flags: Flags{Profile: "prod"},
		Config: File{
			Profiles: map[string]Profile{
				"prod": {APIKeyEnv: "REDASH_API_KEY_PROD"},
			},
		},
		GetSecret: func(profile string) (string, error) {
			return "", errors.New("keyring backend failed")
		},
	})
	if err == nil {
		t.Fatal("Resolve() error = nil, want error")
	}
	var runtimeErr *RuntimeError
	if !errors.As(err, &runtimeErr) {
		t.Fatalf("Resolve() error type = %T, want %T", err, &RuntimeError{})
	}
}
