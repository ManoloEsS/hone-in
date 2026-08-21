package recipe

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/ManoloEsS/hone-in/internal/database"
)

func TestCreateRecipeTrimsName(t *testing.T) {
	repository := newTestRepository(t)
	now := time.Date(2026, time.August, 17, 20, 0, 0, 0, time.UTC)

	created, err := CreateRecipe(context.Background(), repository, "  Apple Pie  ", now)
	if err != nil {
		t.Fatalf("create recipe: %v", err)
	}
	if created.Name != "Apple Pie" {
		t.Fatalf("expected trimmed recipe name %q, got %q", "Apple Pie", created.Name)
	}
}

func TestCreateRecipeRejectsBlankName(t *testing.T) {
	repository := newTestRepository(t)

	_, err := CreateRecipe(context.Background(), repository, " \t ", time.Now())
	if !errors.Is(err, ErrRecipeNameRequired) {
		t.Fatalf("expected required-name error, got %v", err)
	}

	recipes, err := repository.ListRecipes(context.Background())
	if err != nil {
		t.Fatalf("list recipes: %v", err)
	}
	if len(recipes) != 0 {
		t.Fatalf("expected blank name to create no recipes, got %#v", recipes)
	}
}

func newTestRepository(t *testing.T) *Repository {
	t.Helper()

	db, err := database.Open(filepath.Join(t.TempDir(), "recipes.db"))
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	return NewRepository(db)
}
