package web

import (
	"net/http"

	"github.com/ManoloEsS/hone-in/internal/recipe"
)

type catalogPage struct {
	Recipes []recipe.Recipe
}

func (handler *Handler) catalogPage(response http.ResponseWriter, request *http.Request) {
	recipes, err := handler.repository.ListRecipes(request.Context())
	if err != nil {
		http.Error(response, "unable to load recipes", http.StatusInternalServerError)
		return
	}

	if err := handler.templates.ExecuteTemplate(response, "catalog.html", catalogPage{Recipes: recipes}); err != nil {
		http.Error(response, "unable to render recipes", http.StatusInternalServerError)
	}
}
