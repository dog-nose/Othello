package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	cfg := Load()
	if cfg.DBHost != "localhost" {
		t.Fatalf("expected DBHost localhost, got %s", cfg.DBHost)
	}
	if cfg.DBPort != "3306" {
		t.Fatalf("expected DBPort 3306, got %s", cfg.DBPort)
	}
	if cfg.DBName != "othello" {
		t.Fatalf("expected DBName othello, got %s", cfg.DBName)
	}
}

func TestLoadWithEnv(t *testing.T) {
	os.Setenv("DB_HOST", "myhost")
	os.Setenv("DB_PORT", "3307")
	defer os.Unsetenv("DB_HOST")
	defer os.Unsetenv("DB_PORT")

	cfg := Load()
	if cfg.DBHost != "myhost" {
		t.Fatalf("expected DBHost myhost, got %s", cfg.DBHost)
	}
	if cfg.DBPort != "3307" {
		t.Fatalf("expected DBPort 3307, got %s", cfg.DBPort)
	}
}

func TestLoadTest(t *testing.T) {
	cfg := LoadTest()
	if cfg.DBName != "othello_test" {
		t.Fatalf("expected DBName othello_test, got %s", cfg.DBName)
	}
}

func TestDSN(t *testing.T) {
	cfg := &Config{
		DBUser:     "user",
		DBPassword: "pass",
		DBHost:     "host",
		DBPort:     "3306",
		DBName:     "db",
	}
	expected := "user:pass@tcp(host:3306)/db?parseTime=true"
	if dsn := cfg.DSN(); dsn != expected {
		t.Fatalf("expected DSN %s, got %s", expected, dsn)
	}
}
