package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrint_JSONMode(t *testing.T) {
	t.Parallel()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	out := New(Options{
		JSON:   true,
		Stdout: stdout,
		Stderr: stderr,
	})

	err := out.Print(map[string]any{"name": "alice"})
	if err != nil {
		t.Fatalf("Print() error = %v", err)
	}

	got := stdout.String()
	if !strings.Contains(got, `"name":"alice"`) {
		t.Fatalf("Print() output = %q, want JSON payload", got)
	}
}

func TestPrint_TextMode(t *testing.T) {
	t.Parallel()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	out := New(Options{
		JSON:   false,
		Stdout: stdout,
		Stderr: stderr,
	})

	err := out.Print("hello")
	if err != nil {
		t.Fatalf("Print() error = %v", err)
	}

	if got, want := stdout.String(), "hello\n"; got != want {
		t.Fatalf("Print() = %q, want %q", got, want)
	}
}

func TestErrorf_WritesStderr(t *testing.T) {
	t.Parallel()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	out := New(Options{
		JSON:   false,
		Stdout: stdout,
		Stderr: stderr,
	})

	out.Errorf("failed: %s", "boom")

	if got, want := stderr.String(), "failed: boom\n"; got != want {
		t.Fatalf("Errorf() = %q, want %q", got, want)
	}
}
