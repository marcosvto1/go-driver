package folders

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
)

type handler struct {
	db *sql.DB
}

func SetRoute(r chi.Router, db *sql.DB) {
	h := handler{
		db: db,
	}

	r.Post("", h.Create)
	r.Put("/{id}", h.Modify)
	r.Get("/{id}", h.Get)
	r.Delete("/{id}", h.Delete)
	r.Get("/", h.List)
}
