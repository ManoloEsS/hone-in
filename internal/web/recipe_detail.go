package web

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/ManoloEsS/hone-in/internal/recipe"
)

type recipeDetailPage struct {
	Recipe recipe.Recipe
}

func (handler *Handler) recipeDetail(response http.ResponseWriter, request *http.Request) {
	id, err := strconv.ParseInt(request.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		http.NotFound(response, request)
		return
	}

	recipe, err := handler.repository.GetRecipe(request.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		http.NotFound(response, request)
		return
	}
	if err != nil {
		http.Error(response, "unable to load recipe", http.StatusInternalServerError)
		return
	}

	if err := handler.templates.ExecuteTemplate(response, "recipe_detail.html", recipeDetailPage{Recipe: recipe}); err != nil {
		http.Error(response, "unable to render recipe", http.StatusInternalServerError)
	}
}
