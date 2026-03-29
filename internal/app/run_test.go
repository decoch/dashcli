package app

import (
	"bytes"
	"context"
	"testing"

	"github.com/decoch/dashcli/internal/config"
	"github.com/decoch/dashcli/internal/exitcode"
)

func TestRun_Version(t *testing.T) {
	previousVersion := Version
	Version = "1.2.3"
	t.Cleanup(func() {
		Version = previousVersion
	})

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	code := Run(context.Background(), []string{"version"}, stdout, stderr)
	if code != exitcode.CodeSuccess {
		t.Fatalf("Run() code = %d, want %d", code, exitcode.CodeSuccess)
	}
	if got, want := stdout.String(), "1.2.3\n"; got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestRun_MeRequiresCredentials(t *testing.T) {
	previousLoadConfig := loadConfig
	previousLookupEnv := lookupEnv
	loadConfig = func() (config.File, error) {
		return config.File{}, nil
	}
	lookupEnv = func(string) (string, bool) {
		return "", false
	}
	t.Cleanup(func() {
		loadConfig = previousLoadConfig
		lookupEnv = previousLookupEnv
	})

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	code := Run(context.Background(), []string{"me"}, stdout, stderr)
	if code != exitcode.CodeUsage {
		t.Fatalf("Run() code = %d, want %d", code, exitcode.CodeUsage)
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
	if stderr.Len() == 0 {
		t.Fatal("stderr is empty, want usage error output")
	}
}
