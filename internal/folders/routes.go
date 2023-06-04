package folders

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
)

type handler struct {
	db *sql.DB
}

func SetupRoute(r chi.Router, db *sql.DB) {
	h := handler{
		db: db,
	}

	r.Post("", h.Create)
}
