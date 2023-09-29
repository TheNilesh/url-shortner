package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/thenilesh/url-shortner/rest"
)

func main() {
	// TODO: add gorilla mux code
	r := mux.NewRouter()
	s := rest.ShortURL{}

	r.HandleFunc("/", s.Create).Methods("POST")
	r.HandleFunc("/{id}", s.Get).Methods("GET")
	// r.HandleFunc("/{id}", s.Create).Methods("PUT")
	r.HandleFunc("/{id}", s.Delete).Methods("DELETE")
	http.Handle("/", r)

	http.ListenAndServe(":8080", nil)
}
