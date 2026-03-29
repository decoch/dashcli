package exitcode

import (
	"errors"
	"testing"
)

func TestCode_Nil(t *testing.T) {
	t.Parallel()

	if got := Code(nil); got != CodeSuccess {
		t.Fatalf("Code(nil) = %d, want %d", got, CodeSuccess)
	}
}

func TestCode_RuntimeError(t *testing.T) {
	t.Parallel()

	err := Runtimef("something failed")
	if got := Code(err); got != CodeRuntime {
		t.Fatalf("Code(Runtimef) = %d, want %d", got, CodeRuntime)
	}
}

func TestCode_UsageError(t *testing.T) {
	t.Parallel()

	err := Usagef("bad input")
	if got := Code(err); got != CodeUsage {
		t.Fatalf("Code(Usagef) = %d, want %d", got, CodeUsage)
	}
}

func TestCode_WrapRuntime(t *testing.T) {
	t.Parallel()

	err := WrapRuntime(errors.New("inner"))
	if got := Code(err); got != CodeRuntime {
		t.Fatalf("Code(WrapRuntime) = %d, want %d", got, CodeRuntime)
	}
}

func TestCode_WrapUsage(t *testing.T) {
	t.Parallel()

	err := WrapUsage(errors.New("inner"))
	if got := Code(err); got != CodeUsage {
		t.Fatalf("Code(WrapUsage) = %d, want %d", got, CodeUsage)
	}
}

func TestCode_WrapNil(t *testing.T) {
	t.Parallel()

	if got := Code(WrapRuntime(nil)); got != CodeSuccess {
		t.Fatalf("Code(WrapRuntime(nil)) = %d, want %d", got, CodeSuccess)
	}
	if got := Code(WrapUsage(nil)); got != CodeSuccess {
		t.Fatalf("Code(WrapUsage(nil)) = %d, want %d", got, CodeSuccess)
	}
}

func TestCode_PlainError(t *testing.T) {
	t.Parallel()

	err := errors.New("plain error")
	if got := Code(err); got != CodeRuntime {
		t.Fatalf("Code(plain error) = %d, want %d (fallback)", got, CodeRuntime)
	}
}

func TestError_Message(t *testing.T) {
	t.Parallel()

	err := Runtimef("failed: %s", "reason")
	if got := err.Error(); got != "failed: reason" {
		t.Fatalf("Error() = %q, want %q", got, "failed: reason")
	}
}

func TestError_NilInner(t *testing.T) {
	t.Parallel()

	err := &Error{Code: CodeRuntime, Err: nil}
	if got := err.Error(); got != "unknown error" {
		t.Fatalf("Error() = %q, want %q", got, "unknown error")
	}
}

func TestError_Unwrap(t *testing.T) {
	t.Parallel()

	inner := errors.New("inner")
	err := WrapRuntime(inner)
	if !errors.Is(err, inner) {
		t.Fatalf("errors.Is(WrapRuntime(inner), inner) = false, want true")
	}
}
