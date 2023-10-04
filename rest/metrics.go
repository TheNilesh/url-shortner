package rest

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/thenilesh/url-shortner/metrics"
)

type MetricsHandler interface {
	Get(w http.ResponseWriter, r *http.Request)
}

type metricsHandler struct {
	log     *logrus.Logger
	metrics metrics.Metrics
}

func NewMetricsHandler(log *logrus.Logger, m metrics.Metrics) MetricsHandler {
	return &metricsHandler{
		log:     log,
		metrics: m,
	}
}

func (m *metricsHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	for _, kv := range m.metrics.GetDomainCounts() {
		w.Write([]byte(kv.Key))
		w.Write([]byte(": "))
		w.Write([]byte(fmt.Sprintf("%d", kv.Value)))
		w.Write([]byte("\n"))
	}
}
