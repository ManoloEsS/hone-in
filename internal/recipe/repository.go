package recipe

import (
	"context"
	"database/sql"
	"time"

	dbgen "github.com/ManoloEsS/hone-in/internal/database/sqlc"
)

type Recipe struct {
	ID        int64
	Name      string
	CreatedAt string
	UpdatedAt string
}

type Repository struct {
	queries *dbgen.Queries
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{queries: dbgen.New(db)}
}

func (repository *Repository) ListRecipes(ctx context.Context) ([]Recipe, error) {
	rows, err := repository.queries.ListRecipes(ctx)
	if err != nil {
		return nil, err
	}

	recipes := make([]Recipe, len(rows))
	for index, row := range rows {
		recipes[index] = recipeFromRow(row)
	}

	return recipes, nil
}

func (repository *Repository) CreateRecipe(ctx context.Context, name string, now time.Time) (Recipe, error) {
	timestamp := now.UTC().Format(time.RFC3339Nano)
	row, err := repository.queries.CreateRecipe(ctx, dbgen.CreateRecipeParams{
		Name:      name,
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
	})
	if err != nil {
		return Recipe{}, err
	}

	return recipeFromRow(row), nil
}

func recipeFromRow(row dbgen.Recipe) Recipe {
	return Recipe{
		ID:        row.ID,
		Name:      row.Name,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}
