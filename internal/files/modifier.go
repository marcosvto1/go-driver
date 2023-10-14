package files

import (
	"database/sql"
	"encoding/json"
	"fmt"
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

	file, err := Get(h.db, int64(id))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&file)
	if err != nil {
		fmt.Println("decode err")

		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = file.Validate()
	if err != nil {
		fmt.Println("validate err")
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = Update(h.db, int64(id), file)
	if err != nil {
		fmt.Println("update err")
		fmt.Println(err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
	rw.Header().Add("Content-Type", "application/json")
}

func Update(db *sql.DB, id int64, f *File) error {
	f.ModifiedAt = time.Now()

	query := `UPDATE "files" SET "name"=$1, "modified_at"=$2, "deleted"=$3 where id = $4`

	_, err := db.Exec(query, f.Name, f.ModifiedAt, f.Deleted, id)
	if err != nil {
		return err
	}

	return nil
}
