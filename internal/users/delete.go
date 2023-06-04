package users

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

func (h *handler) Delete(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = Delete(h.db, int64(id))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
}

func Delete(db *sql.DB, id int64) error {
	deleted := true
	modifiedAt := time.Now()

	smtp := `UPDATE users SET "deleted"=$1, "modified_at"=$2`

	_, err := db.Exec(smtp, deleted, modifiedAt)
	if err != nil {
		return err
	}

	return nil
}
