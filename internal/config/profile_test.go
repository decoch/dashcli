package config

import (
	"testing"
	"time"
)

func TestResolve_ProfileEnvBeatsGlobalAPIKey(t *testing.T) {
	t.Parallel()

	resolved, err := Resolve(ResolveInput{
		Flags: Flags{
			Profile: "prod",
			Timeout: 3 * time.Second,
		},
		Config: File{
			Profiles: map[string]Profile{
				"prod": {
					BaseURL:   "https://redash.example.com",
					APIKeyEnv: "REDASH_API_KEY_PROD",
				},
			},
		},
		LookupEnv: func(key string) (string, bool) {
			switch key {
			case "REDASH_API_KEY_PROD":
				return "profile-key", true
			case "REDASH_API_KEY":
				return "global-key", true
			default:
				return "", false
			}
		},
	})
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}

	if resolved.APIKey != "profile-key" {
		t.Fatalf("Resolve().APIKey = %q, want %q", resolved.APIKey, "profile-key")
	}
}

func TestResolve_GlobalFallbackWhenProfileKeyMissing(t *testing.T) {
	t.Parallel()

	resolved, err := Resolve(ResolveInput{
		Flags: Flags{Profile: "prod"},
		Config: File{
			Profiles: map[string]Profile{
				"prod": {
					APIKeyEnv: "REDASH_API_KEY_PROD",
				},
			},
		},
		LookupEnv: func(key string) (string, bool) {
			if key == "REDASH_API_KEY" {
				return "global-key", true
			}
			return "", false
		},
	})
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}

	if resolved.APIKey != "global-key" {
		t.Fatalf("Resolve().APIKey = %q, want %q", resolved.APIKey, "global-key")
	}
}

func TestResolve_FlagAPIKeyWins(t *testing.T) {
	t.Parallel()

	resolved, err := Resolve(ResolveInput{
		Flags: Flags{
			Profile: "prod",
			APIKey:  "flag-key",
		},
		Config: File{
			Profiles: map[string]Profile{
				"prod": {
					APIKeyEnv: "REDASH_API_KEY_PROD",
				},
			},
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

func TestResolve_MissingSelectedProfile(t *testing.T) {
	t.Parallel()

	_, err := Resolve(ResolveInput{
		Flags: Flags{Profile: "missing"},
		Config: File{
			Profiles: map[string]Profile{
				"prod": {},
			},
		},
	})
	if err == nil {
		t.Fatal("Resolve() error = nil, want error")
	}
}

func TestResolve_DefaultProfile(t *testing.T) {
	t.Parallel()

	resolved, err := Resolve(ResolveInput{
		Flags: Flags{},
		Config: File{
			DefaultProfile: "stg",
			Profiles: map[string]Profile{
				"stg": {
					BaseURL: "https://redash-stg.example.com",
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}

	if resolved.Profile != "stg" {
		t.Fatalf("Resolve().Profile = %q, want %q", resolved.Profile, "stg")
	}
	if resolved.BaseURL != "https://redash-stg.example.com" {
		t.Fatalf("Resolve().BaseURL = %q, want %q", resolved.BaseURL, "https://redash-stg.example.com")
	}
}

