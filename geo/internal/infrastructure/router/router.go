package router

import (
	"net/http"
	"projects/LDmitryLD/geo-nginx/geo/internal/modules"

	"github.com/go-chi/chi/v5"
)

func NewRouter(controllers *modules.Controllers) *chi.Mux {
	r := chi.NewRouter()

	setDefaultRoutes(r)

	r.Post("/api/address/search", controllers.Geo.Search)
	r.Post("/api/address/geocode", controllers.Geo.Geocode)

	return r
}

func setDefaultRoutes(r *chi.Mux) {
	r.Get("/swagger", swaggerUI)
	r.Get("/public/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))).ServeHTTP(w, r)
	})
}
