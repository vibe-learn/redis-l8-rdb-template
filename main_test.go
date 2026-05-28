package main

import (
	"context"
	"os"
	"testing"
	"time"
)

// TestAddrDefault — pure unit test. No Redis required.
func TestAddrDefault(t *testing.T) {
	if os.Getenv("REDIS_ADDR") == "" && Addr() != "localhost:6379" {
		t.Errorf("Addr() default = %q, want localhost:6379", Addr())
	}
}

// TestNewClient — конструктор не должен возвращать nil.
func TestNewClient(t *testing.T) {
	rdb := NewClient()
	if rdb == nil {
		t.Fatal("NewClient() returned nil")
	}
	_ = rdb.Close()
}

// TestIntegration — требует запущенный Redis (docker compose up -d).
// SKIPPED по умолчанию, чтобы CI был зелёным из коробки; убери t.Skip,
// когда реализуешь TODO-функции, и пропингуй сервер.
func TestIntegration(t *testing.T) {
	if os.Getenv("REDIS_INTEGRATION") == "" {
		t.Skip("set REDIS_INTEGRATION=1 and run `docker compose up -d` to enable")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	rdb := NewClient()
	defer rdb.Close()
	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Fatalf("ping failed: %v", err)
	}
	// TODO: вызови свои реализованные функции и проверь поведение урока «RDB-снэпшоты».
}
