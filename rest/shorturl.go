package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/thenilesh/url-shortner/svc"
)

type ShortURL struct {
	urlShortner svc.URLShortner
}

func NewShortURL(urlShortner svc.URLShortner) *ShortURL {
	return &ShortURL{
		urlShortner: urlShortner,
	}
}

type URL struct {
	ShortPath string `json:"short_path"` // User provided short path
	TargetURL string `json:"target_url"`
}

func (s *ShortURL) Create(w http.ResponseWriter, r *http.Request) {

	// TODO: Prevent OOM/buffer overflow by not parsing large request body
	var req URL
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	// TODO: Validate request body, shortPath should be alphanumeric, no spaces allowed
	shortPath, err := s.urlShortner.CreateShortPath(req.ShortPath, req.TargetURL)
	if err != nil {
		switch err {
		case svc.ErrServerError:
			// TODO: Include error message in response
			w.WriteHeader(http.StatusInternalServerError)
		case svc.ErrConflict:
			w.WriteHeader(http.StatusConflict)
		// case svc.ErrValidationFailed:
		// 	w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Location", fmt.Sprintf("/%s", shortPath))
	w.WriteHeader(http.StatusCreated)
}

func (s *ShortURL) Get(w http.ResponseWriter, r *http.Request) {
	// TODO: if do not redirect query param is included then return ShortURL resource
	// otherwise redirect user to the targetURL
	vars := mux.Vars(r)
	shortPath := vars["id"]
	targetURL, err := s.urlShortner.GetTargetURL(shortPath)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	http.Redirect(w, r, targetURL, http.StatusMovedPermanently)
}

func (s *ShortURL) Put(w http.ResponseWriter, r *http.Request) {
	//TODO: Implement
}

func (s *ShortURL) Delete(w http.ResponseWriter, r *http.Request) {
	//TODO: Implement
}
