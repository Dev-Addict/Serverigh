package apperror

import (
	"errors"
	"testing"
)

func TestOperationalErrorWrapsCause(t *testing.T) {
	cause := errors.New("read failed")
	err := Wrap(CodeFilesystem, "filesystem error", cause)

	if !errors.Is(err, cause) {
		t.Fatalf("expected wrapped cause")
	}

	opErr, ok := AsOperational(err)
	if !ok {
		t.Fatalf("expected operational error")
	}

	if opErr.Code != CodeFilesystem {
		t.Fatalf("expected filesystem code, got %q", opErr.Code)
	}
}

func TestOperationalErrorMatchesByCode(t *testing.T) {
	err := Wrap(
		CodeOutsideRoot,
		"path is outside configured root",
		errors.New("nested"),
	)
	target := New(CodeOutsideRoot, "different message")

	if !errors.Is(err, target) {
		t.Fatalf("expected operational errors to match by code")
	}
}

func TestWrapOperationPreservesCause(t *testing.T) {
	cause := errors.New("stat failed")
	err := WrapOperation(CodeFilesystem, "filesystem error", "inspect path", cause)

	if !errors.Is(err, cause) {
		t.Fatalf("expected wrapped operation cause")
	}

	if err.Error() != "filesystem error: inspect path: stat failed" {
		t.Fatalf("unexpected error message %q", err.Error())
	}
}
