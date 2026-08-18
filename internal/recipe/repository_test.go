package recipe

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/ManoloEsS/hone-in/internal/database"
)

func TestRepositoryListsRecipesByName(t *testing.T) {
	db, err := database.Open(filepath.Join(t.TempDir(), "recipes.db"))
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	for _, name := range []string{"Zucchini Soup", "Apple Pie"} {
		_, err := db.Exec("INSERT INTO recipes (name, created_at, updated_at) VALUES (?, ?, ?)", name, "2026-08-16T00:00:00Z", "2026-08-16T00:00:00Z")
		if err != nil {
			t.Fatalf("insert recipe %q: %v", name, err)
		}
	}

	recipes, err := NewRepository(db).ListRecipes(context.Background())
	if err != nil {
		t.Fatalf("list recipes: %v", err)
	}

	if len(recipes) != 2 {
		t.Fatalf("expected 2 recipes, got %d", len(recipes))
	}
	if recipes[0].Name != "Apple Pie" || recipes[1].Name != "Zucchini Soup" {
		t.Fatalf("expected recipes ordered by name, got %#v", recipes)
	}
}

func TestRepositoryCreatesRecipe(t *testing.T) {
	db, err := database.Open(filepath.Join(t.TempDir(), "recipes.db"))
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	now := time.Date(2026, time.August, 17, 20, 0, 0, 0, time.UTC)
	created, err := NewRepository(db).CreateRecipe(context.Background(), "Apple Pie", now)
	if err != nil {
		t.Fatalf("create recipe: %v", err)
	}

	if created.ID == 0 {
		t.Fatal("expected created recipe ID")
	}
	if created.Name != "Apple Pie" {
		t.Fatalf("expected created recipe name %q, got %q", "Apple Pie", created.Name)
	}
	wantTimestamp := now.Format(time.RFC3339Nano)
	if created.CreatedAt != wantTimestamp || created.UpdatedAt != wantTimestamp {
		t.Fatalf("expected matching timestamps %q, got created_at %q and updated_at %q", wantTimestamp, created.CreatedAt, created.UpdatedAt)
	}

	recipes, err := NewRepository(db).ListRecipes(context.Background())
	if err != nil {
		t.Fatalf("list recipes: %v", err)
	}
	if len(recipes) != 1 || recipes[0] != created {
		t.Fatalf("expected persisted recipe %#v, got %#v", created, recipes)
	}
}
