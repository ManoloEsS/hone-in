package web

import (
	"embed"
	"errors"
	"html/template"
	"io/fs"
	"net/http"
	"time"

	"github.com/ManoloEsS/hone-in/internal/recipe"
)

//go:embed templates/catalog.html templates/new_recipe.html static/css/catalog.css
var assetFiles embed.FS

type Handler struct {
	repository *recipe.Repository
	templates  *template.Template
}

type catalogPage struct {
	Recipes []recipe.Recipe
}

type newRecipePage struct {
	Name  string
	Error string
}

func NewHandler(repository *recipe.Repository) (http.Handler, error) {
	templates, err := template.ParseFS(assetFiles, "templates/catalog.html", "templates/new_recipe.html")
	if err != nil {
		return nil, err
	}
	staticFiles, err := fs.Sub(assetFiles, "static")
	if err != nil {
		return nil, err
	}

	handler := &Handler{repository: repository, templates: templates}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", handler.health)
	mux.HandleFunc("GET /recipes", handler.catalogPage)
	mux.HandleFunc("GET /recipes/new", handler.newRecipeForm)
	mux.HandleFunc("POST /recipes", handler.createRecipe)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFiles))))
	mux.HandleFunc("/", handler.redirectToCatalog)

	return mux, nil
}

func (handler *Handler) health(response http.ResponseWriter, _ *http.Request) {
	response.WriteHeader(http.StatusOK)
	_, _ = response.Write([]byte("healthcheck"))
}

func (handler *Handler) redirectToCatalog(response http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		http.NotFound(response, request)
		return
	}
	http.Redirect(response, request, "/recipes", http.StatusSeeOther)
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
