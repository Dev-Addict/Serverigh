package config

import (
	"testing"

	"serverigh/internal/apperror"
)

func TestNormalizeReturnsOperationalConfigError(t *testing.T) {
	cfg := Config{
		Root:            t.TempDir(),
		Port:            4173,
		MaxPreviewBytes: 1024,
	}

	err := cfg.Normalize()
	if !apperror.HasCode(err, apperror.CodeInvalidConfig) {
		t.Fatalf("expected invalid config operational error, got %v", err)
	}
}

func TestNormalizeReturnsOperationalNotFoundError(t *testing.T) {
	cfg := Config{
		Root:            "/definitely/missing/serverigh/root",
		Host:            "127.0.0.1",
		Port:            4173,
		MaxPreviewBytes: 1024,
	}

	err := cfg.Normalize()
	if !apperror.HasCode(err, apperror.CodeNotFound) {
		t.Fatalf("expected not found operational error, got %v", err)
	}
}

func TestAddressFormatsIPv4Host(t *testing.T) {
	cfg := Config{
		Host: "127.0.0.1",
		Port: 4173,
	}

	if cfg.Address() != "127.0.0.1:4173" {
		t.Fatalf("expected IPv4 address, got %q", cfg.Address())
	}
}

func TestAddressFormatsIPv6Host(t *testing.T) {
	cfg := Config{
		Host: "::1",
		Port: 4173,
	}

	if cfg.Address() != "[::1]:4173" {
		t.Fatalf("expected IPv6 address, got %q", cfg.Address())
	}
}
