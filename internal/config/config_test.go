package config

import "testing"

func TestLoad_MissingToken(t *testing.T) {
	t.Setenv("BOT_TOKEN", "")

	if _, err := Load(); err == nil {
		t.Fatal("expected error when BOT_TOKEN is empty, got nil")
	}
}

func TestLoad_WithToken(t *testing.T) {
	t.Setenv("BOT_TOKEN", "test-token")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.BotToken != "test-token" {
		t.Fatalf("BotToken = %q, want %q", cfg.BotToken, "test-token")
	}
}
