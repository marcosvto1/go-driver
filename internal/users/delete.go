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

	rw.WriteHeader(http.StatusNoContent)
	rw.Header().Set("Content-Type", "application/json")
}

func Delete(db *sql.DB, id int64) error {
	modifiedAt := time.Now()

	smtp := `UPDATE users SET "modified_at"=$1, "deleted"=true, WHERE id=$2`

	_, err := db.Exec(smtp, modifiedAt, id)
	if err != nil {
		return err
	}

	return nil
}
