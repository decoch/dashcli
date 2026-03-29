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

func TestAuthStatus_NotFoundIsSuccess(t *testing.T) {
	previousGet := authGetSecret
	previousInput := authInput
	authGetSecret = func(profile string) (string, error) {
		if profile != "default" {
			t.Fatalf("profile = %q, want %q", profile, "default")
		}
		return "", keyring.ErrKeyNotFound
	}
	authInput = strings.NewReader("")
	t.Cleanup(func() {
		authGetSecret = previousGet
		authInput = previousInput
	})

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	code := Run(context.Background(), []string{"auth", "status"}, stdout, stderr)

	if code != exitcode.CodeSuccess {
		t.Fatalf("Run() code = %d, want %d", code, exitcode.CodeSuccess)
	}
	if got, want := stdout.String(), "No API key stored for profile default\n"; got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestAuthStatus_ProfileFlag(t *testing.T) {
	previousGet := authGetSecret
	authGetSecret = func(profile string) (string, error) {
		if profile != "prod" {
			t.Fatalf("profile = %q, want %q", profile, "prod")
		}
		return "secret", nil
	}
	t.Cleanup(func() {
		authGetSecret = previousGet
	})

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	code := Run(context.Background(), []string{"--profile", "prod", "auth", "status"}, stdout, stderr)

	if code != exitcode.CodeSuccess {
		t.Fatalf("Run() code = %d, want %d", code, exitcode.CodeSuccess)
	}
	if got, want := stdout.String(), "API key is set for profile prod\n"; got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
}

func TestAuthSet_StoresForDefaultProfile(t *testing.T) {
	previousSet := authSetSecret
	previousInput := authInput

	calledProfile := ""
	calledKey := ""
	authSetSecret = func(profile, apiKey string) error {
		calledProfile = profile
		calledKey = apiKey
		return nil
	}
	authInput = strings.NewReader("abc123\n")
	t.Cleanup(func() {
		authSetSecret = previousSet
		authInput = previousInput
	})

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	code := Run(context.Background(), []string{"auth", "set"}, stdout, stderr)

	if code != exitcode.CodeSuccess {
		t.Fatalf("Run() code = %d, want %d", code, exitcode.CodeSuccess)
	}
	if calledProfile != "default" {
		t.Fatalf("set profile = %q, want %q", calledProfile, "default")
	}
	if calledKey != "abc123" {
		t.Fatalf("set apiKey = %q, want %q", calledKey, "abc123")
	}
	if got, want := stdout.String(), "API key stored for profile default\n"; got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
}

func TestAuthDelete_UsesProfileFlag(t *testing.T) {
	previousDelete := authDeleteSecret
	calledProfile := ""
	authDeleteSecret = func(profile string) error {
		calledProfile = profile
		return nil
	}
	t.Cleanup(func() {
		authDeleteSecret = previousDelete
	})

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	code := Run(context.Background(), []string{"--profile", "stg", "auth", "delete"}, stdout, stderr)

	if code != exitcode.CodeSuccess {
		t.Fatalf("Run() code = %d, want %d", code, exitcode.CodeSuccess)
	}
	if calledProfile != "stg" {
		t.Fatalf("delete profile = %q, want %q", calledProfile, "stg")
	}
	if got, want := stdout.String(), "API key removed for profile stg\n"; got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
}

func TestAuthStatus_RuntimeError(t *testing.T) {
	previousGet := authGetSecret
	authGetSecret = func(profile string) (string, error) {
		return "", errors.New("backend failure")
	}
	t.Cleanup(func() {
		authGetSecret = previousGet
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
