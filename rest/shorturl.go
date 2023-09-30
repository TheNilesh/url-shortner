package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/thenilesh/url-shortner/store"
	"github.com/thenilesh/url-shortner/svc"
)

type ShortURL struct {
	// TODO: Use svc.urlShortner
	// TODO: use store
	store       store.KVStore
	urlShortner svc.URLShortner
}

type URL struct {
	ShortPath string `json:"short_path"` // User provided short path
	TargetURL string `json:"target_url"`
}

func (s *ShortURL) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request received for create")

	// TODO: Prevent OOM/buffer overflow  by not parsing large request body
	var data URL
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	// TODO: Cache backend may not be fast, consider using context
	// and map it with http timeout. ie, if http request times out,
	// cancel the context

	var shortURL string
	// FIXME: Following loop will execute indefinitely if shortPaths are exhausted
	for exists := true; exists; {
		shortURL = s.urlShortner.Shorten(data.TargetURL)

		var err error
		exists, err = s.store.Exists(shortURL)
		if err != nil {
			// TODO: Log
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	err := s.store.Put(shortURL, data.TargetURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (s *ShortURL) Get(w http.ResponseWriter, r *http.Request) {
	// TODO: if do not redirect query param is included then return ShortURL resource
	// otherwise redirect user to the targetURL

}

func (s *ShortURL) Put(w http.ResponseWriter, r *http.Request) {
	//TODO: Implement
}

func (s *ShortURL) Delete(w http.ResponseWriter, r *http.Request) {
	//TODO: Implement
}
