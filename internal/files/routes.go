package files

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/marcosvto1/go-driver/internal/bucket"
	"github.com/marcosvto1/go-driver/internal/queue"
)

type handler struct {
	db     *sql.DB
	bucket *bucket.Bucket
	queue  *queue.Queue
}

func SetRoutes(r chi.Router, db *sql.DB, bucket *bucket.Bucket, queue *queue.Queue) {
	h := handler{
		db:     db,
		bucket: bucket,
		queue:  queue,
	}

	r.Put("/{id}", h.Modify)
	r.Post("", h.Create)
	r.Delete("/{id}", h.Delete)
}
