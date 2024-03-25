package main

import (
	"context"
	"fmt"
	"net/http"
	"series/pkg/config"
	"series/pkg/routes"
)

func main() {
	client := config.Connect()
	defer client.Disconnect(context.Background())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Func Handler!")
	})

	routes.SetupAuthRoutes(client)

	http.ListenAndServe(":8080", nil)
}
