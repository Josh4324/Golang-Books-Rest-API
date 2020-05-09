package main

import (
	"fmt"
	"net/http"
)

type books map[string]string

func (b books) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	
	switch r.URL.String(){
	case "/home":
		fmt.Fprintf(w, "Welcome to books server. You came via %s", r.URL.String())
	case "/about":
		fmt.Fprintf(w, "All about our library of books. You came via %s", r.URL.String())
	default:
		w.WriteHeader(404);
	}
}

func main() {
	//Http Server
	/* http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, Go HTTP\n")
	})

	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "About, HTTP in Go\n")
	}) */

	var bk books
	fmt.Println("server listening on port 8000")
	http.ListenAndServe(":8000", bk)
}
