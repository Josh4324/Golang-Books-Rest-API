package router

import (
	"go-rest-api/middleware"

	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/v1/books/{id}", middleware.GetBook).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/books", middleware.GetAllBook).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/books", middleware.CreateBook).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/books/{id}", middleware.UpdateBook).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/v1/books/{id}", middleware.DeleteBook).Methods("DELETE", "OPTIONS")

	return router
}
