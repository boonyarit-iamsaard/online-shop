package config

import (
	"os"
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

	_, err := Load()
	if err == nil {
		t.Fatal("Load() error = nil, want error")
	}

	if err.Error() != "config: DATABASE_URL is required" {
		t.Errorf("Load() error = %q, want %q", err.Error(), "config: DATABASE_URL is required")
	}
}

func TestLoadUsesDotenvFile(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("PORT", "")
	t.Setenv("DATABASE_URL", "")
	writeDotenv(t, "PORT=3000\nDATABASE_URL=postgres://dotenv\n")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Port != "3000" {
		t.Errorf("Port = %q, want %q", cfg.Port, "3000")
	}
	if cfg.DatabaseURL != "postgres://dotenv" {
		t.Errorf("DatabaseURL = %q, want %q", cfg.DatabaseURL, "postgres://dotenv")
	}
}

func TestLoadEnvironmentOverridesDotenvFile(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("PORT", "9999")
	t.Setenv("DATABASE_URL", "")
	writeDotenv(t, "PORT=3000\nDATABASE_URL=postgres://dotenv\n")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Port != "9999" {
		t.Errorf("Port = %q, want %q", cfg.Port, "9999")
	}
	if cfg.DatabaseURL != "postgres://dotenv" {
		t.Errorf("DatabaseURL = %q, want %q", cfg.DatabaseURL, "postgres://dotenv")
	}
}

func writeDotenv(t *testing.T, content string) {
	t.Helper()

	if err := os.WriteFile(".env", []byte(content), 0o600); err != nil {
		t.Fatalf("write .env: %v", err)
	}
}
