package routes

import (
	"net/http"
	"series/pkg/controllers"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupAuthRoutes(client *mongo.Client) {
	http.HandleFunc("/register", controllers.Register(client))
	// http.HandleFunc("/login", controllers.Login(client))
	// http.HandleFunc("/logout", controllers.Logout())
}
