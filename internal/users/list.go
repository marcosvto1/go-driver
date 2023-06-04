package users

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func (h *handler) List(rw http.ResponseWriter, r *http.Request) {
	users, err := FindAll(h.db)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(users)
}

func FindAll(db *sql.DB) ([]User, error) {
	query := `SELECT id, name, login, password, created_at, modified_at, deleted, last_login FROM users
	WHERE deleted = false`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(
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

		users = append(users, user)
	}

	return users, nil
}
