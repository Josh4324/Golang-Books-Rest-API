package main

import (
	"net/http"
	"time"
	"encoding/json"
)

type Book struct {
	Name string `json:"name"`
	Author string `json:"author"`
	PublishedAt string `json:"published_at"`
}

type books map[string]string

/* func (b books) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.URL.Path {
	case "/home":
		fmt.Fprintf(w, "Welcome to books server. You came via %s", r.URL.String())
	case "/about":
		fmt.Fprintf(w, "All about our library of books. You came via %s", r.URL.String())
	case "/book":
		item := r.URL.Query().Get("item")
		book, ok := b[item]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "book not found: %q\n", item)
			return
		}
		fmt.Fprintf(w, "The author is %s", book)

	default:
		w.WriteHeader(http.StatusNotFound)
	}
} */

func loadBooks(w http.ResponseWriter, r *http.Request){
	book := []Book{ Book{"The Silmarilion", "JRR Tolkein", time.Now().Local().String()},
					Book{"The Hobbit", "JRR Tolkein", time.Now().Local().String()},
	}

	bk, err := json.Marshal(book);
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bk)
}

func main() {
	//Http Server
	/* http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, Go HTTP\n")
	})

	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "About, HTTP in Go\n")
	}) */

	/* bk := books{}
	bk["jrrtolkein"] = "The Silmarilion"
	fmt.Println("server listening on port 8000") */
	//http.ListenAndServe(":8000", bk)

	mux := http.NewServeMux()
	mux.HandleFunc("/books", loadBooks)

	fs := http.FileServer(http.Dir("assets/"))
	mux.Handle("/static/", http.StripPrefix("/static", fs))
	http.ListenAndServe(":8000", mux)
}
