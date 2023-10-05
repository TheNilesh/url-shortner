package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/thenilesh/url-shortner/metrics"
	"github.com/thenilesh/url-shortner/rest"
	"github.com/thenilesh/url-shortner/store"
	"github.com/thenilesh/url-shortner/svc"
)

const appName = "url-shortner"

func main() {
	// TODO: Main function is too big, need to refactor by adding app.go
	initViper()
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	log := logrus.New()
	strLogLvl := viper.GetString("log_level")
	logLevel, err := logrus.ParseLevel(strLogLvl)
	if err != nil {
		log.WithError(err).Fatal("Failed to parse log level")
	}
	log.Level = logLevel

	log.Info("Starting server")
	r := mux.NewRouter()
	r.Use(RequestIDMiddleware)
	metrics := metrics.NewMetrics()
	metrics.Start()
	metricsHandler := rest.NewMetricsHandler(log, metrics)
	urlShortner := buildURLShortner(log, metrics)
	s := rest.NewShortURLHandler(log, urlShortner)
	log.Info("Registering metrics route")
	r.HandleFunc("/metrics", metricsHandler.Get).Methods("GET")
	log.Info("Registering other routes")
	r.HandleFunc("/", s.Create).Methods("POST")
	r.HandleFunc("/{id}", s.Get).Methods("GET")
	r.HandleFunc("/{id}", s.Put).Methods("PUT")
	r.HandleFunc("/{id}", s.Delete).Methods("DELETE")
	http.Handle("/", r)

	listenAddr := viper.GetString("listen_addr")
	log.Infof("Starting listening on %s", listenAddr)
	http.ListenAndServe(listenAddr, nil)
}

func buildURLShortner(log *logrus.Logger, metrics metrics.Metrics) svc.URLShortner {
	redisAddr := viper.GetString("redis_addr")
	redisPassword := viper.GetString("redis_addr")
	redisDB := viper.GetInt("redis_db")
	redis, err := store.NewRedisClient(redisAddr, redisPassword, redisDB)
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
	return svc.NewURLShortner(6, targetURLStore, shortPathStore, metrics)
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

func initViper() {
	viper.SetDefault("log_level", "info")
	viper.SetDefault("listen_addr", ":8080")
	viper.SetDefault("redis_addr", "localhost:6379")
	viper.SetDefault("redis_password", "")
	viper.SetDefault("redis_db", 0)

	viper.SetConfigName("config")
	viper.AddConfigPath(fmt.Sprintf("/etc/%s", appName))
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("us")
	viper.AutomaticEnv()
}
