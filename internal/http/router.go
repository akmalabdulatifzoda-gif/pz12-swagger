package httpx

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/akmalabdulatifzoda-gif/pz12-swagger/internal/http/handlers"
	"github.com/akmalabdulatifzoda-gif/pz12-swagger/internal/http/redoc"
)

// NewRouter собирает маршруты. apiToken — токен для защищённых методов
// (пустая строка = авторизация выключена).
func NewRouter(h *handlers.Handler, apiToken string) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1/notes", func(r chi.Router) {
		r.Get("/", h.ListNotes)
		r.Get("/{id}", h.GetNote)

		// Изменяющие методы защищены Bearer-токеном (@Security BearerAuth).
		r.Group(func(r chi.Router) {
			r.Use(BearerAuth(apiToken))
			r.Post("/", h.CreateNote)
			r.Patch("/{id}", h.PatchNote)
			r.Delete("/{id}", h.DeleteNote)
		})
	})

	// Документация: Swagger UI на /docs, ReDoc на /redoc.
	r.Get("/docs", http.RedirectHandler("/docs/index.html", http.StatusMovedPermanently).ServeHTTP)
	r.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.URL("/docs/doc.json"),
		httpSwagger.DocExpansion("list"),
		httpSwagger.PersistAuthorization(true),
	))
	r.Get("/redoc", redoc.Page)
	r.Get("/redoc/redoc.standalone.js", redoc.Script)
	return r
}
