package controller

import (
	"encoding/json"
	"net/http"
	"model" // import the model package
)

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User
		json.NewDecoder(r.Body).Decode(&user)

		existingUser, err := model.GetUser(db, user.Username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if existingUser.Password != user.Password {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}

		json.NewEncoder(w).Encode(existingUser)
	}
}

func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User
		json.NewDecoder(r.Body).Decode(&user)

		err := model.StoreUser(db, user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(user)
	}
}
