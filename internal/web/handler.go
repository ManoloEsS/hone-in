package web

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/ManoloEsS/hone-in/internal/recipe"
)

//go:embed templates/catalog.html templates/new_recipe.html templates/recipe_detail.html static/css/catalog.css
var assetFiles embed.FS

type Handler struct {
	repository *recipe.Repository
	templates  *template.Template
}

func NewHandler(repository *recipe.Repository) (http.Handler, error) {
	templates, err := template.ParseFS(assetFiles, "templates/catalog.html", "templates/new_recipe.html", "templates/recipe_detail.html")
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
	mux.HandleFunc("GET /recipes/{id}", handler.recipeDetail)
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
