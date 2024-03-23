package main

import (
	"database/sql"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func StoreUser(db *sql.DB, user User) error {
	sql_additem := `
	INSERT OR REPLACE INTO users(
		username,
		password) values(?, ?)
	`

	stmt, err := db.Prepare(sql_additem)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err2 := stmt.Exec(user.Username, user.Password)
	if err2 != nil {
		return err2
	}
	return nil
}

func GetUser(db *sql.DB, username string) (*User, error) {
	sql_getitem := `
	SELECT Password FROM users WHERE username = ?
	`

	stmt, err := db.Prepare(sql_getitem)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var password string
	err2 := stmt.QueryRow(username).Scan(&password)
	if err2 != nil {
		if err2 != sql.ErrNoRows {
			return nil, err2
		}
	}

	user := &User{
		Username: username,
		Password: password,
	}

	return user, nil
}
