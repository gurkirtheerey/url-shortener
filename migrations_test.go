package main

import (
	"context"
	"testing"
)

func TestRunMigrations_Idempotent(t *testing.T) {
	if err := runMigrations(context.Background(), testStore.pool); err != nil {
		t.Fatalf("expected migrations to be idempotent: %v", err)
	}

	var exists bool
	err := testStore.pool.QueryRow(
		context.Background(),
		"SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE name = $1)",
		"001_init.sql",
	).Scan(&exists)
	if err != nil {
		t.Fatalf("failed to query schema_migrations: %v", err)
	}

	if !exists {
		t.Fatal("expected 001_init.sql to be recorded in schema_migrations")
	}
}
