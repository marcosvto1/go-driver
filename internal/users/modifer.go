package users

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

func (h *handler) Modify(rw http.ResponseWriter, r *http.Request) {
	user := new(User)

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if user.Name == "" {
		http.Error(rw, ErrNameRequired.Error(), http.StatusBadRequest)
		return
	}

	err = Update(h.db, int64(id), user)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: GetId

	rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(user)
}

func Update(db *sql.DB, id int64, u *User) error {
	u.ModifiedAt = time.Now()

	stmt := `UPDATE "users" SET "name"=$1, "login"=$2, "password"=$3 "modified_at"=$4 WHERE id=$5`

	_, err := db.Exec(stmt, u.Name, u.Login, u.Password, u.ModifiedAt, id)
	if err != nil {
		return err
	}

	return nil
}
