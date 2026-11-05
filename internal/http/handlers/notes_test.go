package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/akmalabdulatifzoda-gif/pz12-swagger/internal/core/service"
	httpx "github.com/akmalabdulatifzoda-gif/pz12-swagger/internal/http"
	"github.com/akmalabdulatifzoda-gif/pz12-swagger/internal/http/handlers"
	"github.com/akmalabdulatifzoda-gif/pz12-swagger/internal/repo"

	_ "github.com/akmalabdulatifzoda-gif/pz12-swagger/docs"
)

func newServer(token string) *httptest.Server {
	h := &handlers.Handler{Svc: service.NewNoteService(repo.NewNoteRepoMem())}
	return httptest.NewServer(httpx.NewRouter(h, token))
}

func request(t *testing.T, srv *httptest.Server, method, path, body, auth string) *http.Response {
	t.Helper()
	req, _ := http.NewRequest(method, srv.URL+path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if auth != "" {
		req.Header.Set("Authorization", auth)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	return resp
}

func TestCRUD(t *testing.T) {
	srv := newServer("")
	defer srv.Close()

	do := func(method, path, body string) int {
		req, _ := http.NewRequest(method, srv.URL+path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		return resp.StatusCode
	}

	cases := []struct {
		method, path, body string
		want               int
	}{
		{"POST", "/api/v1/notes", `{"title":"Первая","content":"тест"}`, 201},
		{"POST", "/api/v1/notes", `{"title":""}`, 400},
		{"POST", "/api/v1/notes", `{bad json`, 400},
		{"GET", "/api/v1/notes", "", 200},
		{"GET", "/api/v1/notes/1", "", 200},
		{"GET", "/api/v1/notes/99", "", 404},
		{"GET", "/api/v1/notes/abc", "", 400},
		{"PATCH", "/api/v1/notes/1", `{"content":"обновлено"}`, 200},
		{"PATCH", "/api/v1/notes/99", `{"content":"x"}`, 404},
		{"DELETE", "/api/v1/notes/1", "", 204},
		{"DELETE", "/api/v1/notes/1", "", 404},
	}
	for _, c := range cases {
		if got := do(c.method, c.path, c.body); got != c.want {
			t.Errorf("%s %s: got %d, want %d", c.method, c.path, got, c.want)
		}
	}
}

func TestListPaginationAndSearch(t *testing.T) {
	srv := newServer("")
	defer srv.Close()
	for _, title := range []string{"Go", "Swagger", "Go модули"} {
		request(t, srv, "POST", "/api/v1/notes", `{"title":"`+title+`"}`, "")
	}
	resp := request(t, srv, "GET", "/api/v1/notes?q=go&limit=1", "", "")
	if got := resp.Header.Get("X-Total-Count"); got != "2" {
		t.Errorf("X-Total-Count = %s, want 2", got)
	}
	if resp := request(t, srv, "GET", "/api/v1/notes?page=abc", "", ""); resp.StatusCode != 400 {
		t.Errorf("bad page: got %d, want 400", resp.StatusCode)
	}
}

func TestBearerAuth(t *testing.T) {
	srv := newServer("secret")
	defer srv.Close()
	cases := []struct {
		auth string
		want int
	}{
		{"", 401},
		{"Bearer wrong", 401},
		{"Bearer secret", 201},
	}
	for _, c := range cases {
		if got := request(t, srv, "POST", "/api/v1/notes", `{"title":"x"}`, c.auth).StatusCode; got != c.want {
			t.Errorf("auth %q: got %d, want %d", c.auth, got, c.want)
		}
	}
	// чтение доступно без токена
	if got := request(t, srv, "GET", "/api/v1/notes", "", "").StatusCode; got != 200 {
		t.Errorf("GET without token: got %d, want 200", got)
	}
}

func TestDocsRoutes(t *testing.T) {
	srv := newServer("")
	defer srv.Close()
	for _, p := range []string{"/docs/index.html", "/docs/doc.json", "/redoc", "/redoc/redoc.standalone.js"} {
		if got := request(t, srv, "GET", p, "", "").StatusCode; got != 200 {
			t.Errorf("GET %s: got %d, want 200", p, got)
		}
	}
}
