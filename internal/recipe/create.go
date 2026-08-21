package recipe

import (
	"context"
	"errors"
	"strings"
	"time"
)

var ErrRecipeNameRequired = errors.New("recipe name is required")

func CreateRecipe(ctx context.Context, repository *Repository, name string, now time.Time) (Recipe, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Recipe{}, ErrRecipeNameRequired
	}

	return repository.CreateRecipe(ctx, name, now)
}
