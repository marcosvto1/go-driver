package files

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
	rw.Header().Add("Content-Type", "application/json")
}

func Delete(db *sql.DB, id int64) error {
	query := `UPDATE files SET "deleted"=$1, "modified_at"=$2 WHERE "id"=$3`

	_, err := db.Exec(query, true, time.Now(), id)
	if err != nil {
		return err
	}

	return nil
}
