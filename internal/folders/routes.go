package folders

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/marcosvto1/go-driver/internal/auth"
)

type handler struct {
	db *sql.DB
}

func SetRoute(r chi.Router, db *sql.DB) {
	h := handler{
		db: db,
	}

	r.Route("/folders", func(r chi.Router) {
		r.Use(auth.Validate)

		r.Post("/", h.Create)
		r.Put("/{id}", h.Modify)
		r.Get("/{id}", h.Get)
		r.Delete("/{id}", h.Delete)
		r.Get("/", h.List)
	})

}
