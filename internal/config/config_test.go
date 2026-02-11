package config

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"
)

// helper — создаем тестовый yaml
func writeTempConfig(t *testing.T, content string) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	return path
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}

func TestLoad_Success(t *testing.T) {
	yaml := `
interval: 10s
port: 8080
log_level: info

sites:
  - id: site-1
    name: Google
    url: https://google.com

database:
  host: localhost
  port: 5432
  user: postgres
  name: testdb
`

	path := writeTempConfig(t, yaml)

	cfg, err := Load(path, testLogger())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Interval != 10_000_000_000 { // 10s
		t.Fatalf("expected interval 10s, got %v", cfg.Interval)
	}

	if len(cfg.Sites) != 1 {
		t.Fatalf("expected 1 site")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("not-exists.yaml", testLogger())

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestLoad_EmptySites(t *testing.T) {
	yaml := `
interval: 10s

database:
  host: localhost
  port: 5432
  user: postgres
  name: testdb
`

	path := writeTempConfig(t, yaml)

	_, err := Load(path, testLogger())
	if err == nil {
		t.Fatal("expected error for empty sites")
	}
}

func TestLoad_InvalidInterval(t *testing.T) {
	yaml := `
interval: 0s

sites:
  - id: site-1
    name: Google
    url: https://google.com

database:
  host: localhost
  port: 5432
  user: postgres
  name: testdb
`

	path := writeTempConfig(t, yaml)

	_, err := Load(path, testLogger())
	if err == nil {
		t.Fatal("expected interval validation error")
	}
}

func TestLoad_DBValidation(t *testing.T) {
	yaml := `
interval: 10s

sites:
  - id: site-1
    name: Google
    url: https://google.com

database: {}
`

	path := writeTempConfig(t, yaml)

	_, err := Load(path, testLogger())
	if err == nil {
		t.Fatal("expected DB validation error")
	}
}

func TestLoad_EnvOverride(t *testing.T) {
	yaml := `
interval: 10s

sites:
  - id: site-1
    name: Google
    url: https://google.com

database:
  host: localhost
  port: 5432
  user: postgres
  name: testdb
`

	path := writeTempConfig(t, yaml)

	t.Setenv("APP_PORT", "9999")
	t.Setenv("CHECK_INTERVAL", "30s")

	cfg, err := Load(path, testLogger())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Port != 9999 {
		t.Fatalf("env override failed, expected 9999 got %d", cfg.Port)
	}

	if cfg.Interval.String() != "30s" {
		t.Fatalf("interval env override failed")
	}
}

func TestLoad_DatabaseURL(t *testing.T) {
	yaml := `
interval: 10s

sites:
  - id: site-1
    name: Google
    url: https://google.com

database: {}
`

	path := writeTempConfig(t, yaml)

	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db")

	cfg, err := Load(path, testLogger())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.DB.URL == "" {
		t.Fatal("expected DATABASE_URL to be set")
	}
}
