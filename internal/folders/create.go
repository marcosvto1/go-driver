package folders

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

func (h *handler) Create(rw http.ResponseWriter, r *http.Request) {
	folder := new(Folder)

	err := json.NewDecoder(r.Body).Decode(folder)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = folder.Validate()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := Insert(h.db, folder)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	folder.ID = id

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	json.NewEncoder(rw).Encode(folder)
}

func Insert(db *sql.DB, folder *Folder) (id int64, err error) {
	folder.ModifiedAt = time.Now()

	query := `INSERT INTO "folders" ("name", "parent_id", "modified_at") VALUES ($1, $2, $3) RETURNING "id"`

	var LastInsertId int64

	err = db.QueryRow(query, folder.Name, folder.ParentID, folder.ModifiedAt).Scan(&LastInsertId)
	if err != nil {
		return -1, err
	}

	return LastInsertId, err
}
