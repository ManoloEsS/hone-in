package recipe

import (
	"context"
	"database/sql"

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
		recipes[index] = Recipe{
			ID:        row.ID,
			Name:      row.Name,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
		}
	}

	return recipes, nil
}
