package web

import (
	"errors"
	"net/http"
	"time"

	"github.com/ManoloEsS/hone-in/internal/recipe"
)

type newRecipePage struct {
	Name  string
	Error string
}

func (handler *Handler) newRecipeForm(response http.ResponseWriter, _ *http.Request) {
	handler.renderNewRecipeForm(response, http.StatusOK, newRecipePage{})
}

func (handler *Handler) createRecipe(response http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		http.Error(response, "invalid recipe form", http.StatusBadRequest)
		return
	}

	name := request.FormValue("name")
	_, err := recipe.CreateRecipe(request.Context(), handler.repository, name, time.Now())
	if errors.Is(err, recipe.ErrRecipeNameRequired) {
		handler.renderNewRecipeForm(response, http.StatusUnprocessableEntity, newRecipePage{
			Name:  name,
			Error: "Recipe name is required.",
		})
		return
	}
	if err != nil {
		http.Error(response, "unable to create recipe", http.StatusInternalServerError)
		return
	}

	http.Redirect(response, request, "/recipes", http.StatusSeeOther)
}

func (handler *Handler) renderNewRecipeForm(response http.ResponseWriter, status int, page newRecipePage) {
	response.WriteHeader(status)
	if err := handler.templates.ExecuteTemplate(response, "new_recipe.html", page); err != nil {
		http.Error(response, "unable to render recipe form", http.StatusInternalServerError)
	}
}
