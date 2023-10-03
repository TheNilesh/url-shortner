package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/thenilesh/url-shortner/svc"
)

const (
	MaxRequestBodySize = 1048576 // 1 MB
)

type RequestIDKey string

type ShortURLHandler struct {
	urlShortner svc.URLShortner
	log         *logrus.Logger
}

func NewShortURLHandler(log *logrus.Logger, urlShortner svc.URLShortner) *ShortURLHandler {
	return &ShortURLHandler{
		log:         log,
		urlShortner: urlShortner,
	}
}

type ShortURL struct {
	// Optional user provided short path
	ShortPath string `json:"short_path"`
	TargetURL string `json:"target_url"`
}

func (s *ShortURLHandler) Create(w http.ResponseWriter, r *http.Request) {
	requestID, _ := r.Context().Value(RequestIDKey("requestID")).(string)
	log := s.log.WithField("requestID", requestID)
	log.Infof("Received request. %s %s", r.Method, r.URL.Path)

	// TODO: Handle case when request is cancelled by client

	// Prevent OOM/buffer overflow
	limitReader := io.LimitReader(r.Body, MaxRequestBodySize)
	var shortURL ShortURL
	decoder := json.NewDecoder(limitReader)
	if err := decoder.Decode(&shortURL); err != nil {
		if err == io.EOF {
			log.Error(err)
			http.Error(w, "Request body is too large", http.StatusRequestEntityTooLarge)
			return
		}
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	shortPath, err := s.urlShortner.CreateShortPath(context.TODO(), shortURL.ShortPath, shortURL.TargetURL)
	if err != nil {
		log.Error(err)
		switch err.(type) {
		case *svc.ErrValidation:
			// TODO: Return validation error in response body
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Location", fmt.Sprintf("/%s", shortPath))
	w.WriteHeader(http.StatusCreated)
	log.Infof("Sent response.Location: %s", shortPath)
}

func (s *ShortURLHandler) Get(w http.ResponseWriter, r *http.Request) {
	requestID, _ := r.Context().Value(RequestIDKey("requestID")).(string)
	log := s.log.WithField("requestID", requestID)
	log.Infof("Received request. %s %s", r.Method, r.URL.Path)

	// TODO: if do not redirect query param is included then return ShortURL resource
	// otherwise redirect user to the targetURL
	vars := mux.Vars(r)
	shortPath := vars["id"]
	targetURL, err := s.urlShortner.GetTargetURL(context.TODO(), shortPath)
	if err != nil {
		log.Errorf("Failed to get targetURL for shortPath: %s", shortPath)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	http.Redirect(w, r, targetURL, http.StatusMovedPermanently)
	log.Infof("Redirected to targetURL. %s->%s", shortPath, targetURL)
}

func (s *ShortURLHandler) Put(w http.ResponseWriter, r *http.Request) {
	//TODO: Implement
}

func (s *ShortURLHandler) Delete(w http.ResponseWriter, r *http.Request) {
	//TODO: Implement
}
