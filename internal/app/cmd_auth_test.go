package app

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/99designs/keyring"

	"github.com/decoch/dashcli/internal/exitcode"
)

func TestAuthSet_StoresBaseURLAndAPIKey(t *testing.T) {
	previousSetBaseURL := authSetBaseURL
	previousSetAPIKey := authSetAPIKey
	previousInput := authInput

	calledBaseURL := ""
	calledAPIKey := ""
	authSetBaseURL = func(baseURL string) error {
		calledBaseURL = baseURL
		return nil
	}
	authSetAPIKey = func(apiKey string) error {
		calledAPIKey = apiKey
		return nil
	}
	authInput = strings.NewReader("https://redash.example.com\nabc123\n")
	t.Cleanup(func() {
		authSetBaseURL = previousSetBaseURL
		authSetAPIKey = previousSetAPIKey
		authInput = previousInput
	})

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	code := Run(context.Background(), []string{"auth", "set"}, stdout, stderr)

	if code != exitcode.CodeSuccess {
		t.Fatalf("Run() code = %d, want %d", code, exitcode.CodeSuccess)
	}
	if calledBaseURL != "https://redash.example.com" {
		t.Fatalf("baseURL = %q, want %q", calledBaseURL, "https://redash.example.com")
	}
	if calledAPIKey != "abc123" {
		t.Fatalf("apiKey = %q, want %q", calledAPIKey, "abc123")
	}
	if !strings.Contains(stdout.String(), "Credentials stored in keyring\n") {
		t.Fatalf("stdout = %q, want success message", stdout.String())
	}
}

func TestAuthSet_RejectsHTTPBaseURL(t *testing.T) {
	previousInput := authInput
	authInput = strings.NewReader("http://redash.example.com\nabc123\n")
	t.Cleanup(func() {
		authInput = previousInput
	})

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	code := Run(context.Background(), []string{"auth", "set"}, stdout, stderr)

	if code != exitcode.CodeUsage {
		t.Fatalf("Run() code = %d, want %d", code, exitcode.CodeUsage)
	}
}

func TestAuthDelete_SuccessWhenKeysMissing(t *testing.T) {
	previousDeleteBaseURL := authDeleteBaseURL
	previousDeleteAPIKey := authDeleteAPIKey
	authDeleteBaseURL = func() error {
		return keyring.ErrKeyNotFound
	}
	authDeleteAPIKey = func() error {
		return keyring.ErrKeyNotFound
	}
	t.Cleanup(func() {
		authDeleteBaseURL = previousDeleteBaseURL
		authDeleteAPIKey = previousDeleteAPIKey
	})

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	code := Run(context.Background(), []string{"auth", "delete"}, stdout, stderr)

	if code != exitcode.CodeSuccess {
		t.Fatalf("Run() code = %d, want %d", code, exitcode.CodeSuccess)
	}
	if got, want := stdout.String(), "Credentials removed from keyring\n"; !strings.Contains(got, want) {
		t.Fatalf("stdout = %q, want message containing %q", got, want)
	}
}

func TestAuthStatus_AllMissing(t *testing.T) {
	previousGetBaseURL := authGetBaseURL
	previousGetAPIKey := authGetAPIKey
	authGetBaseURL = func() (string, error) {
		return "", keyring.ErrKeyNotFound
	}
	authGetAPIKey = func() (string, error) {
		return "", keyring.ErrKeyNotFound
	}
	t.Cleanup(func() {
		authGetBaseURL = previousGetBaseURL
		authGetAPIKey = previousGetAPIKey
	})

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	code := Run(context.Background(), []string{"auth", "status"}, stdout, stderr)

	if code != exitcode.CodeSuccess {
		t.Fatalf("Run() code = %d, want %d", code, exitcode.CodeSuccess)
	}
	want := "No base URL stored\nNo API key stored\n"
	if got := stdout.String(); got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
}

func TestAuthStatus_RuntimeError(t *testing.T) {
	previousGetBaseURL := authGetBaseURL
	authGetBaseURL = func() (string, error) {
		return "", errors.New("backend failure")
	}
	t.Cleanup(func() {
		authGetBaseURL = previousGetBaseURL
	})

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	code := Run(context.Background(), []string{"auth", "status"}, stdout, stderr)

	if code != exitcode.CodeRuntime {
		t.Fatalf("Run() code = %d, want %d", code, exitcode.CodeRuntime)
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
}
