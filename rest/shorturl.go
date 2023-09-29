package rest

import (
	"fmt"
	"net/http"
)

type ShortURL struct {
	// TODO: Use svc.urlShortner
	// TODO: use store
}

func (s *ShortURL) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request received for create")
	w.Write([]byte("OK"))
}

func (s *ShortURL) Get(w http.ResponseWriter, r *http.Request) {
	// TODO: if do not redirect query param is included then return ShortURL resource
	// otherwise redirect user to the targetURL
	print("Request received for get")
}

func (s *ShortURL) Put(w http.ResponseWriter, r *http.Request) {
	//TODO: Implement
}

func (s *ShortURL) Delete(w http.ResponseWriter, r *http.Request) {
	//TODO: Implement
}
