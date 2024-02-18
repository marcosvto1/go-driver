package users

import (
	"fmt"
	"time"
)

func (h *handler) authenticate(login, password string) (*User, error) {
	stmt := `select * from users where login = $1 and password = $2`

	fmt.Println(login, password)

	row := h.db.QueryRow(stmt, login, encPass(password))
	u := User{}
	err := row.Scan(&u.ID, &u.Name, &u.Login, &u.Password, &u.CreatedAt, &u.ModifiedAt, &u.LastLogin, &u.Deleted)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (h *handler) updateLastLogin(u *User) error {
	u.LastLogin = time.Now()
	return Update(h.db, u.ID, u)
}

func Authenticate(login, password string) (u *User, err error) {
	u, err = gh.authenticate(login, password)
	if err != nil {
		return
	}

	err = gh.updateLastLogin(u)
	if err != nil {
		println("Error updating last login")
		return nil, err
	}

	return u, nil
}
