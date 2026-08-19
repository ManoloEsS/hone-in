package web

import (
	"database/sql"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/ManoloEsS/hone-in/internal/database"
	"github.com/ManoloEsS/hone-in/internal/recipe"
)

func TestCatalogShowsEmptyState(t *testing.T) {
	server, db := newTestServer(t)
	defer server.Close()
	defer db.Close()

	response, err := server.Client().Get(server.URL + "/recipes")
	if err != nil {
		t.Fatalf("get empty catalog: %v", err)
	}
	body := readResponseBody(t, response)

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected empty catalog status %d, got %d", http.StatusOK, response.StatusCode)
	}
	if !strings.Contains(body, "No recipes yet.") {
		t.Fatalf("expected empty state in response, got %q", body)
	}
	if !strings.Contains(body, "/static/css/catalog.css") {
		t.Fatalf("expected catalog stylesheet link in response, got %q", body)
	}
	if !strings.Contains(body, "New recipe") {
		t.Fatalf("expected new-recipe action in response, got %q", body)
	}
}

func TestNewRecipeForm(t *testing.T) {
	server, db := newTestServer(t)
	defer server.Close()
	defer db.Close()

	response, err := server.Client().Get(server.URL + "/recipes/new")
	if err != nil {
		t.Fatalf("get new recipe form: %v", err)
	}
	body := readResponseBody(t, response)

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected new recipe form status %d, got %d", http.StatusOK, response.StatusCode)
	}
	if !strings.Contains(body, "<form method=\"post\" action=\"/recipes\">") {
		t.Fatalf("expected recipe form in response, got %q", body)
	}
}

func TestCreateRecipeRedirectsAndPersistsTrimmedName(t *testing.T) {
	server, db := newTestServer(t)
	defer server.Close()
	defer db.Close()

	client := *server.Client()
	client.CheckRedirect = func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}
	response, err := client.PostForm(server.URL+"/recipes", url.Values{"name": {"  Apple Pie  "}})
	if err != nil {
		t.Fatalf("submit new recipe: %v", err)
	}
	response.Body.Close()

	if response.StatusCode != http.StatusSeeOther {
		t.Fatalf("expected recipe creation status %d, got %d", http.StatusSeeOther, response.StatusCode)
	}
	if location := response.Header.Get("Location"); location != "/recipes" {
		t.Fatalf("expected creation redirect to /recipes, got %q", location)
	}

	var name string
	if err := db.QueryRow("SELECT name FROM recipes").Scan(&name); err != nil {
		t.Fatalf("read created recipe: %v", err)
	}
	if name != "Apple Pie" {
		t.Fatalf("expected trimmed stored name %q, got %q", "Apple Pie", name)
	}
}

func TestCreateRecipeRejectsBlankName(t *testing.T) {
	server, db := newTestServer(t)
	defer server.Close()
	defer db.Close()

	response, err := server.Client().PostForm(server.URL+"/recipes", url.Values{"name": {" \t "}})
	if err != nil {
		t.Fatalf("submit blank recipe: %v", err)
	}
	body := readResponseBody(t, response)

	if response.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("expected blank recipe status %d, got %d", http.StatusUnprocessableEntity, response.StatusCode)
	}
	if !strings.Contains(body, "Recipe name is required.") {
		t.Fatalf("expected validation error in response, got %q", body)
	}

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM recipes").Scan(&count); err != nil {
		t.Fatalf("count recipes: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected blank recipe to create no records, got %d", count)
	}
}

func TestCatalogServesStylesheet(t *testing.T) {
	server, db := newTestServer(t)
	defer server.Close()
	defer db.Close()

	response, err := server.Client().Get(server.URL + "/static/css/catalog.css")
	if err != nil {
		t.Fatalf("get catalog stylesheet: %v", err)
	}
	body := readResponseBody(t, response)

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected stylesheet status %d, got %d", http.StatusOK, response.StatusCode)
	}
	if !strings.Contains(body, "font-family") {
		t.Fatalf("expected stylesheet content, got %q", body)
	}
}

func TestHealthEndpointReturnsOK(t *testing.T) {
	server, db := newTestServer(t)
	defer server.Close()
	defer db.Close()

	response, err := server.Client().Get(server.URL + "/healthz")
	if err != nil {
		t.Fatalf("get health endpoint: %v", err)
	}
	body := readResponseBody(t, response)

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected health endpoint status %d, got %d", http.StatusOK, response.StatusCode)
	}
	if body != "healthcheck" {
		t.Fatalf("expected healthcheck body, got %q", body)
	}
}

func TestCatalogShowsStoredRecipes(t *testing.T) {
	server, db := newTestServer(t)
	defer server.Close()
	defer db.Close()

	recipeIDs := make(map[string]int64)
	for _, name := range []string{"Zucchini Soup", "Apple Pie"} {
		result, err := db.Exec("INSERT INTO recipes (name, created_at, updated_at) VALUES (?, ?, ?)", name, "2026-08-16T00:00:00Z", "2026-08-16T00:00:00Z")
		if err != nil {
			t.Fatalf("insert recipe %q: %v", name, err)
		}
		id, err := result.LastInsertId()
		if err != nil {
			t.Fatalf("read recipe ID %q: %v", name, err)
		}
		recipeIDs[name] = id
	}

	response, err := server.Client().Get(server.URL + "/recipes")
	if err != nil {
		t.Fatalf("get populated catalog: %v", err)
	}
	body := readResponseBody(t, response)

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected populated catalog status %d, got %d", http.StatusOK, response.StatusCode)
	}
	for _, name := range []string{"Apple Pie", "Zucchini Soup"} {
		if !strings.Contains(body, name) {
			t.Fatalf("expected %q in response, got %q", name, body)
		}
		link := `href="/recipes/` + strconv.FormatInt(recipeIDs[name], 10) + `"`
		if !strings.Contains(body, link) {
			t.Fatalf("expected recipe link %q in response, got %q", link, body)
		}
	}
}

func TestRecipeDetailShowsStoredRecipe(t *testing.T) {
	server, db := newTestServer(t)
	defer server.Close()
	defer db.Close()

	result, err := db.Exec("INSERT INTO recipes (name, created_at, updated_at) VALUES (?, ?, ?)", "Apple Pie", "2026-08-16T00:00:00Z", "2026-08-16T00:00:00Z")
	if err != nil {
		t.Fatalf("insert recipe: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("read recipe ID: %v", err)
	}

	response, err := server.Client().Get(server.URL + "/recipes/" + strconv.FormatInt(id, 10))
	if err != nil {
		t.Fatalf("get recipe detail: %v", err)
	}
	body := readResponseBody(t, response)

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected recipe detail status %d, got %d", http.StatusOK, response.StatusCode)
	}
	if !strings.Contains(body, "<h1>Apple Pie</h1>") {
		t.Fatalf("expected recipe name in response, got %q", body)
	}

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM recipes").Scan(&count); err != nil {
		t.Fatalf("count recipes after detail request: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected detail request to preserve recipe count, got %d", count)
	}
}

func TestRecipeDetailReturnsNotFoundForMissingRecipe(t *testing.T) {
	server, db := newTestServer(t)
	defer server.Close()
	defer db.Close()

	response, err := server.Client().Get(server.URL + "/recipes/999")
	if err != nil {
		t.Fatalf("get missing recipe detail: %v", err)
	}
	response.Body.Close()

	if response.StatusCode != http.StatusNotFound {
		t.Fatalf("expected missing recipe status %d, got %d", http.StatusNotFound, response.StatusCode)
	}
}

func TestRootRedirectsToCatalog(t *testing.T) {
	server, db := newTestServer(t)
	defer server.Close()
	defer db.Close()

	client := *server.Client()
	client.CheckRedirect = func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}
	response, err := client.Get(server.URL + "/")
	if err != nil {
		t.Fatalf("get root: %v", err)
	}
	response.Body.Close()

	if response.StatusCode != http.StatusSeeOther {
		t.Fatalf("expected root status %d, got %d", http.StatusSeeOther, response.StatusCode)
	}
	if location := response.Header.Get("Location"); location != "/recipes" {
		t.Fatalf("expected redirect to /recipes, got %q", location)
	}
}

func newTestServer(t *testing.T) (*httptest.Server, *sql.DB) {
	db, err := database.Open(filepath.Join(t.TempDir(), "recipes.db"))
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	handler, err := NewHandler(recipe.NewRepository(db))
	if err != nil {
		db.Close()
		t.Fatalf("create web handler: %v", err)
	}

	return httptest.NewServer(handler), db
}

func readResponseBody(t *testing.T, response *http.Response) string {
	t.Helper()
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("read response body: %v", err)
	}
	return string(body)
}
