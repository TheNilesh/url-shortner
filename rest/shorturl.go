package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/thenilesh/url-shortner/svc"
)

const MaxRequestBodySize = 1048576 // 1 MB

type ShortURLHandler struct {
	urlShortner svc.URLShortner
}

func NewShortURL(urlShortner svc.URLShortner) *ShortURLHandler {
	return &ShortURLHandler{
		urlShortner: urlShortner,
	}
}

type ShortURL struct {
	// Optional user provided short path, may be blank
	ShortPath string `json:"short_path"`
	TargetURL string `json:"target_url"`
}

func (s *ShortURLHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Prevent OOM/buffer overflow
	limitReader := io.LimitReader(r.Body, MaxRequestBodySize)
	var shortURL ShortURL
	decoder := json.NewDecoder(limitReader)
	if err := decoder.Decode(&shortURL); err != nil {
		if err == io.EOF {
			http.Error(w, "Request body is too large", http.StatusRequestEntityTooLarge)
			return
		}
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	shortPath, err := s.urlShortner.CreateShortPath(shortURL.ShortPath, shortURL.TargetURL)
	if err != nil {
		switch err {
		case svc.ErrServerError:
			// TODO: Include error message in response
			// TODO: Log error
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

func (s *ShortURLHandler) Get(w http.ResponseWriter, r *http.Request) {
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

func (s *ShortURLHandler) Put(w http.ResponseWriter, r *http.Request) {
	//TODO: Implement
}

func (s *ShortURLHandler) Delete(w http.ResponseWriter, r *http.Request) {
	//TODO: Implement
}
