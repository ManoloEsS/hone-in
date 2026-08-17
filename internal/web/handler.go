package web

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/ManoloEsS/hone-in/internal/recipe"
)

//go:embed templates/catalog.html static/css/catalog.css
var assetFiles embed.FS

type Handler struct {
	repository *recipe.Repository
	catalog    *template.Template
}

type catalogPage struct {
	Recipes []recipe.Recipe
}

func NewHandler(repository *recipe.Repository) (http.Handler, error) {
	catalog, err := template.ParseFS(assetFiles, "templates/catalog.html")
	if err != nil {
		return nil, err
	}
	staticFiles, err := fs.Sub(assetFiles, "static")
	if err != nil {
		return nil, err
	}

	handler := &Handler{repository: repository, catalog: catalog}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", handler.health)
	mux.HandleFunc("/recipes", handler.catalogPage)
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

	if err := handler.catalog.Execute(response, catalogPage{Recipes: recipes}); err != nil {
		http.Error(response, "unable to render recipes", http.StatusInternalServerError)
	}
}
