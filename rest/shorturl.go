package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/thenilesh/url-shortner/svc"
)

const (
	maxRequestBodySize = 104875 // 1 MB
)

type RequestIDKey string

type ShortURLHandler struct {
	urlShortner *svc.URLShortner
	log         *logrus.Logger
}

func NewShortURLHandler(log *logrus.Logger, urlShortner *svc.URLShortner) *ShortURLHandler {
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

type Response struct {
	RequestID string `json:"request_id"`
	Message   string `json:"message"`
}

func (s *ShortURLHandler) Create(w http.ResponseWriter, r *http.Request) {
	requestID, _ := r.Context().Value(RequestIDKey("requestID")).(string)
	log := s.log.WithField("requestID", requestID)
	log.Infof("Received request. %s %s", r.Method, r.URL.Path)
	// TODO: Handle case when request is cancelled by client using context
	// Prevent OOM/buffer overflow
	limitReader := io.LimitReader(r.Body, maxRequestBodySize)
	var shortURL ShortURL
	decoder := json.NewDecoder(limitReader)
	if err := decoder.Decode(&shortURL); err != nil {
		log.Error(err)
		if err == io.ErrUnexpectedEOF {
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			w.Header().Set("Content-Type", "application/json")
			w.Write(marshalMessage(requestID, "Request body is too large or invalid JSON"))
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.Write(marshalMessage(requestID, "Failed to decode JSON"))
		return
	}
	defer r.Body.Close()

	shortPath, err := s.urlShortner.CreateShortPath(r.Context(), shortURL.ShortPath, shortURL.TargetURL)
	if err != nil {
		log.Error(err)
		switch err.(type) {
		case *svc.ErrValidation:
			w.WriteHeader(http.StatusBadRequest)
		case *svc.ErrConflict:
			w.WriteHeader(http.StatusConflict)
		case *svc.ErrServerError:
			w.WriteHeader(http.StatusInternalServerError)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(marshalMessage(requestID, err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Location", fmt.Sprintf("/%s", shortPath))
	w.Header().Set("Content-Type", "application/json")
	w.Write(marshalMessage(requestID, fmt.Sprintf("Created short URL: /%s", shortPath)))
	log.Infof("Sent response. shortPath:%s", shortPath)
}

func (s *ShortURLHandler) Get(w http.ResponseWriter, r *http.Request) {
	requestID, _ := r.Context().Value(RequestIDKey("requestID")).(string)
	log := s.log.WithField("requestID", requestID)
	log.Infof("Received request. %s %s", r.Method, r.URL.Path)
	// TODO: if do not redirect query param is included then return ShortURL resource
	// otherwise redirect user to the targetURL
	vars := mux.Vars(r)
	shortPath := vars["id"]
	targetURL, err := s.urlShortner.GetTargetURL(r.Context(), shortPath)
	if err != nil {
		log.Errorf("Failed to get targetURL for shortPath: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	http.Redirect(w, r, targetURL, http.StatusMovedPermanently)
	log.Infof("Redirected. %s->%s", shortPath, targetURL)
}

func (s *ShortURLHandler) Put(w http.ResponseWriter, r *http.Request) {
	//TODO: Implement
}

func (s *ShortURLHandler) Delete(w http.ResponseWriter, r *http.Request) {
	//TODO: Implement
}

func marshalMessage(requestID string, msg string) []byte {
	Response := Response{
		RequestID: requestID,
		Message:   msg,
	}
	dataBytes, _ := json.Marshal(Response)
	return dataBytes
}
