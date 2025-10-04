package main

import (
	"log"
	"net/http"

	"github.com/<your-username>/go-practice2/internal/handlers"
	"github.com/<your-username>/go-practice2/internal/middleware"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/user", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetUser)))
	mux.Handle("/user/create", middleware.AuthMiddleware(http.HandlerFunc(handlers.CreateUser)))

	log.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
