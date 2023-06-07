package folders

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

func (h *handler) Modify(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	folder := new(Folder)
	err = json.NewDecoder(r.Body).Decode(folder)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = folder.Validate()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = Update(h.db, int64(id), folder)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(folder)
}

func Update(db *sql.DB, id int64, folder *Folder) error {
	folder.ModifiedAt = time.Now()

	query := `UPDATE folders SET name=$1, modified_at=$2 WHERE id=$3`

	_, err := db.Exec(query, folder.Name, folder.ModifiedAt, id)
	if err != nil {
		return err
	}

	return nil
}
