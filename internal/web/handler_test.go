package web

import (
	"database/sql"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
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

	for _, name := range []string{"Zucchini Soup", "Apple Pie"} {
		_, err := db.Exec("INSERT INTO recipes (name, created_at, updated_at) VALUES (?, ?, ?)", name, "2026-08-16T00:00:00Z", "2026-08-16T00:00:00Z")
		if err != nil {
			t.Fatalf("insert recipe %q: %v", name, err)
		}
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
