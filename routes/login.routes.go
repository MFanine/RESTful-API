package main

import (
	"net/http"
	"controller" // import the controller package
)

func InitRoutes(db *sql.DB) {
	http.HandleFunc("/login", controller.LoginHandler(db))
	http.HandleFunc("/register", controller.RegisterHandler(db))
}
