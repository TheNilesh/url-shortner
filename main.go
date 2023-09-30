package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/thenilesh/url-shortner/rest"
	"github.com/thenilesh/url-shortner/store"
	"github.com/thenilesh/url-shortner/svc"
)

func main() {
	r := mux.NewRouter()
	s := rest.NewShortURL(store.NewKVStore(), store.NewKVStore(), svc.NewRandomURLShortner(5))

	r.HandleFunc("/", s.Create).Methods("POST")
	//TODO: non existent short path returns method not allowed, it should return 404
	r.HandleFunc("/{id}", s.Get).Methods("GET")
	// r.HandleFunc("/{id}", s.Create).Methods("PUT")
	r.HandleFunc("/{id}", s.Delete).Methods("DELETE")
	http.Handle("/", r)

	http.ListenAndServe(":8080", nil)
}
