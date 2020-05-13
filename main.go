package main

import (
	"fmt"
	"go-rest-api/router"
	"log"
	"net/http"
)

func main() {
	r := router.Router()
	// fs := http.FileServer(http.Dir("build"))
	// http.Handle("/", fs)
	fmt.Println("Starting server on the port 6500...")

	log.Fatal(http.ListenAndServe(":6500", r))
}
