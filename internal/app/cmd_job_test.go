package app

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/99designs/keyring"

	"github.com/decoch/dashcli/internal/exitcode"
	"github.com/decoch/dashcli/internal/redash"
)

func stubJobCredentials(t *testing.T) {
	t.Helper()
	prevLookupEnv := lookupEnv
	prevGetAPIKey := getAPIKey
	prevGetBaseURL := getBaseURL
	lookupEnv = func(key string) (string, bool) {
		switch key {
		case "REDASH_BASE_URL":
			return "https://redash.example.com", true
		case "REDASH_API_KEY":
			return "test-key", true
		}
		return "", false
	}
	getAPIKey = func() (string, error) { return "", keyring.ErrKeyNotFound }
	getBaseURL = func() (string, error) { return "", keyring.ErrKeyNotFound }
	t.Cleanup(func() {
		lookupEnv = prevLookupEnv
		getAPIKey = prevGetAPIKey
		getBaseURL = prevGetBaseURL
	})
}

func stubJobGetJob(t *testing.T, responses []map[string]any) {
	t.Helper()
	prev := jobGetJob
	call := 0
	jobGetJob = func(_ context.Context, _ *redash.Client, _ string) (map[string]any, error) {
		resp := responses[call]
		if call < len(responses)-1 {
			call++
		}
		return resp, nil
	}
	t.Cleanup(func() { jobGetJob = prev })
}

func jobResponse(status int) map[string]any {
	return map[string]any{
		"job": map[string]any{
			"id":     "42",
			"status": float64(status),
			"error":  "",
		},
	}
}

func TestJobWait_CompletesImmediately(t *testing.T) {
	stubJobCredentials(t)
	stubJobGetJob(t, []map[string]any{jobResponse(3)})

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	code := Run(context.Background(), []string{"job", "wait", "42", "--interval=1ms", "--max-wait=5s"}, stdout, stderr)

	if code != exitcode.CodeSuccess {
		t.Fatalf("Run() code = %d, want %d", code, exitcode.CodeSuccess)
	}
	if !strings.Contains(stdout.String(), "42") {
		t.Fatalf("stdout = %q, want job id in output", stdout.String())
	}
}

func TestJobWait_PollingLoop(t *testing.T) {
	stubJobCredentials(t)
	stubJobGetJob(t, []map[string]any{
		jobResponse(1),
		jobResponse(1),
		jobResponse(3),
	})

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	code := Run(context.Background(), []string{"job", "wait", "42", "--interval=1ms", "--max-wait=5s"}, stdout, stderr)

	if code != exitcode.CodeSuccess {
		t.Fatalf("Run() code = %d, want %d", code, exitcode.CodeSuccess)
	}
}

func TestJobWait_Timeout(t *testing.T) {
	stubJobCredentials(t)

	prev := jobGetJob
	jobGetJob = func(_ context.Context, _ *redash.Client, _ string) (map[string]any, error) {
		return jobResponse(1), nil
	}
	t.Cleanup(func() { jobGetJob = prev })

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	code := Run(context.Background(), []string{"job", "wait", "42", "--interval=1ms", "--max-wait=50ms"}, stdout, stderr)

	if code != exitcode.CodeRuntime {
		t.Fatalf("Run() code = %d, want %d (timeout)", code, exitcode.CodeRuntime)
	}
}

func TestJobWait_JSONOutput(t *testing.T) {
	stubJobCredentials(t)
	stubJobGetJob(t, []map[string]any{jobResponse(3)})

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	code := Run(context.Background(), []string{"--json", "job", "wait", "42", "--interval=1ms", "--max-wait=5s"}, stdout, stderr)

	if code != exitcode.CodeSuccess {
		t.Fatalf("Run() code = %d, want %d", code, exitcode.CodeSuccess)
	}
	if !strings.Contains(stdout.String(), `"job"`) {
		t.Fatalf("stdout = %q, want JSON with job field", stdout.String())
	}
}

func TestJobWait_ContextCancelled(t *testing.T) {
	stubJobCredentials(t)

	ctx, cancel := context.WithCancel(context.Background())
	callCount := 0
	prev := jobGetJob
	jobGetJob = func(callCtx context.Context, _ *redash.Client, _ string) (map[string]any, error) {
		callCount++
		if callCount >= 2 {
			cancel()
		}
		return jobResponse(1), nil
	}
	t.Cleanup(func() { jobGetJob = prev })

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	code := Run(ctx, []string{"job", "wait", "42", "--interval=1ms", "--max-wait=5s"}, stdout, stderr)

	if code != exitcode.CodeRuntime {
		t.Fatalf("Run() code = %d, want %d (context cancelled)", code, exitcode.CodeRuntime)
	}
}
