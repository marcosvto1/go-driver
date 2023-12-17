package users

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/marcosvto1/go-driver/internal/auth"
)

type handler struct {
	db *sql.DB
}

var gh handler

// SetRoutes sets the routes for the given chi.Router and initializes the handler with the provided *sql.DB.
//
// Parameters:
// - r: the chi.Router to set the routes on.
// - db: the *sql.DB to initialize the handler with.
func SetRoutes(r chi.Router, db *sql.DB) {
	gh = handler{db}

	r.Route("/users", func(r chi.Router) {
		r.Post("", gh.Create)

		r.Group(func(r chi.Router) {
			r.Use(auth.Validate)
			r.Put("/{id}", gh.Modify)
			r.Delete("/{id}", gh.Delete)
			r.Get("/{id}", gh.GetById)
			r.Get("", gh.List)
		})
	})
}
