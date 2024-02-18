package users

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func (h *handler) Create(rw http.ResponseWriter, r *http.Request) {
	u := new(User)

	err := json.NewDecoder(r.Body).Decode(u)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	u.SetPassword(u.Password)

	err = u.Validate()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := Insert(h.db, u)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	u.ID = id

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	json.NewEncoder(rw).Encode(u)
}

func Insert(db *sql.DB, u *User) (int64, error) {
	query := `INSERT INTO "users" ("name", "login", "password", "modified_at") VALUES ($1, $2, $3, $4) RETURNING id`

	stmt, err := db.Prepare(query)
	if err != nil {
		return -1, err
	}

	err = stmt.QueryRow(u.Name, u.Login, u.Password, u.ModifiedAt).Scan(&u.ID)
	if err != nil {
		return -1, err
	}

	return u.ID, nil
}
