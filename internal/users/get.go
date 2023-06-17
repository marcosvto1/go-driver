package users

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *handler) GetById(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := Get(h.db, int64(id))
	if err != nil {
		// TODO: validar se o err e pq não existe nenhum registro
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Header().Add("Content-Type", "application/type")
	json.NewEncoder(rw).Encode(user)
}

func Get(db *sql.DB, id int64) (*User, error) {
	user := User{}
	q := `
		SELECT id, name, login, password, created_at, modified_at, deleted, last_login
		FROM users
		WHERE id = $1"
	`
	err := db.QueryRow(q, id).Scan(
		&user.ID,
		&user.Name,
		&user.Login,
		&user.Password,
		&user.CreatedAt,
		&user.ModifiedAt,
		&user.Deleted,
		&user.LastLogin,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
