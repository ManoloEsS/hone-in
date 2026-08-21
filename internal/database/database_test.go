package database

import (
	"path/filepath"
	"testing"
)

func TestOpenMigratesRecipesAndEnablesForeignKeys(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "recipes.db"))
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	var foreignKeys int
	if err := db.QueryRow("PRAGMA foreign_keys").Scan(&foreignKeys); err != nil {
		t.Fatalf("read foreign-key setting: %v", err)
	}
	if foreignKeys != 1 {
		t.Fatalf("expected foreign keys enabled, got %d", foreignKeys)
	}

	var tableName string
	if err := db.QueryRow("SELECT name FROM sqlite_master WHERE type = 'table' AND name = 'recipes'").Scan(&tableName); err != nil {
		t.Fatalf("find recipes table: %v", err)
	}
	if tableName != "recipes" {
		t.Fatalf("expected recipes table, got %q", tableName)
	}
}
