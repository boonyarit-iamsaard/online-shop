package config

import (
	"os"
	"strings"
	"testing"
)

func TestLoadRequiresPort(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("PORT", "")
	t.Setenv("DATABASE_URL", "")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() error = nil, want error")
	}

	if err.Error() != "config: PORT is required" {
		t.Errorf("Load() error = %q, want %q", err.Error(), "config: PORT is required")
	}
}

func TestLoadRequiresDatabaseURL(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("PORT", "3000")
	t.Setenv("DATABASE_URL", "")

	setRequiredAuthEnv(t)

	_, err := Load()
	if err == nil {
		t.Fatal("Load() error = nil, want error")
	}

	if err.Error() != "config: DATABASE_URL is required" {
		t.Errorf("Load() error = %q, want %q", err.Error(), "config: DATABASE_URL is required")
	}
}

func TestLoadRequiresAppEnv(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("PORT", "3000")
	t.Setenv("DATABASE_URL", "postgres://localhost")
	t.Setenv("APP_ENV", "")
	t.Setenv("AUTH_CUSTOMER_PORTAL_ORIGIN", "http://customer.localhost")
	t.Setenv("AUTH_STAFF_PORTAL_ORIGIN", "http://staff.localhost")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() error = nil, want error")
	}

	if err.Error() != "config: APP_ENV is required" {
		t.Errorf("Load() error = %q, want %q", err.Error(), "config: APP_ENV is required")
	}
}

func TestLoadRequiresCustomerPortalOrigin(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("PORT", "3000")
	t.Setenv("DATABASE_URL", "postgres://localhost")
	t.Setenv("APP_ENV", "test")
	t.Setenv("AUTH_CUSTOMER_PORTAL_ORIGIN", "")
	t.Setenv("AUTH_STAFF_PORTAL_ORIGIN", "http://staff.localhost")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() error = nil, want error")
	}

	if err.Error() != "config: AUTH_CUSTOMER_PORTAL_ORIGIN is required" {
		t.Errorf("Load() error = %q, want %q", err.Error(), "config: AUTH_CUSTOMER_PORTAL_ORIGIN is required")
	}
}

func TestLoadRequiresStaffPortalOrigin(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("PORT", "3000")
	t.Setenv("DATABASE_URL", "postgres://localhost")
	t.Setenv("APP_ENV", "test")
	t.Setenv("AUTH_CUSTOMER_PORTAL_ORIGIN", "http://customer.localhost")
	t.Setenv("AUTH_STAFF_PORTAL_ORIGIN", "")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() error = nil, want error")
	}

	if err.Error() != "config: AUTH_STAFF_PORTAL_ORIGIN is required" {
		t.Errorf("Load() error = %q, want %q", err.Error(), "config: AUTH_STAFF_PORTAL_ORIGIN is required")
	}
}

func TestLoadUsesDotenvFile(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("PORT", "")
	t.Setenv("DATABASE_URL", "")
	writeDotenv(t, requiredDotenv())

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Port != "3000" {
		t.Errorf("Port = %q, want %q", cfg.Port, "3000")
	}
	if cfg.DatabaseURL != "postgres://localhost" {
		t.Errorf("DatabaseURL = %q, want %q", cfg.DatabaseURL, "postgres://localhost")
	}
}

func TestLoadEnvironmentOverridesDotenvFile(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("PORT", "9999")
	t.Setenv("DATABASE_URL", "")
	writeDotenv(t, requiredDotenv())

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Port != "9999" {
		t.Errorf("Port = %q, want %q", cfg.Port, "9999")
	}
	if cfg.DatabaseURL != "postgres://localhost" {
		t.Errorf("DatabaseURL = %q, want %q", cfg.DatabaseURL, "postgres://localhost")
	}
}

func writeDotenv(t *testing.T, content string) {
	t.Helper()

	if err := os.WriteFile(".env", []byte(content), 0o600); err != nil {
		t.Fatalf("write .env: %v", err)
	}
}

func requiredDotenv() string {
	return strings.Join([]string{
		"PORT=3000",
		"DATABASE_URL=postgres://localhost",
		"APP_ENV=test",
		"AUTH_CUSTOMER_PORTAL_ORIGIN=http://customer.localhost",
		"AUTH_STAFF_PORTAL_ORIGIN=http://staff.localhost",
	}, "\n") + "\n"
}

func setRequiredAuthEnv(t *testing.T) {
	t.Helper()

	t.Setenv("APP_ENV", "test")
	t.Setenv("AUTH_CUSTOMER_PORTAL_ORIGIN", "http://customer.localhost")
	t.Setenv("AUTH_STAFF_PORTAL_ORIGIN", "http://staff.localhost")
}
