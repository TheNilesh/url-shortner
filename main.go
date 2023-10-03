package main

import (
	"context"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/thenilesh/url-shortner/rest"
	"github.com/thenilesh/url-shortner/store"
	"github.com/thenilesh/url-shortner/svc"
)

func main() {
	log := logrus.New()
	logLevel := os.Getenv("LOG_LEVEL")
	if len(logLevel) != 0 {
		logLevel, err := logrus.ParseLevel(logLevel)
		if err != nil {
			log.WithError(err).Fatal("Failed to parse log level")
		}
		log.Level = logLevel
	}

	log.Info("Starting server")
	r := mux.NewRouter()
	r.Use(RequestIDMiddleware)
	urlShortner := buildURLShortner(log)
	s := rest.NewShortURLHandler(log, urlShortner)

	log.Info("Registering routes")
	r.HandleFunc("/", s.Create).Methods("POST")
	r.HandleFunc("/{id}", s.Get).Methods("GET")
	r.HandleFunc("/{id}", s.Put).Methods("PUT")
	r.HandleFunc("/{id}", s.Delete).Methods("DELETE")
	http.Handle("/", r)

	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}
	log.Infof("Starting listening on %s", listenAddr)
	http.ListenAndServe(listenAddr, nil)
}

func buildURLShortner(log *logrus.Logger) *svc.URLShortner {
	redis, err := store.NewRedisClient("localhost:6379", "", 0)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to redis")
	}
	err = store.CheckRedisConnection(redis)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to redis")
	}
	targetURLStore, err := store.NewRedisKVStore(redis, "target")
	if err != nil {
		log.WithError(err).Fatal("Failed to create targetURLStore")
	}
	shortPathStore, err := store.NewRedisKVStore(redis, "short")
	if err != nil {
		log.WithError(err).Fatal("Failed to create shortPathStore")
	}
	urlShortner := svc.NewURLShortner(6, targetURLStore, shortPathStore)
	return &urlShortner
}

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		ctx := r.Context()
		ctx = context.WithValue(ctx, rest.RequestIDKey("requestID"), requestID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
