package rest

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/thenilesh/url-shortner/svc"
)

type metrics struct {
	log     *logrus.Logger
	metrics *svc.Metrics
}

func NewMetricsHandler(log *logrus.Logger, m *svc.Metrics) *metrics {
	return &metrics{
		log:     log,
		metrics: m,
	}
}

func (m *metrics) Get(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	for _, kv := range m.metrics.GetDomainCounts() {
		w.Write([]byte(kv.Key))
		w.Write([]byte(": "))
		w.Write([]byte(fmt.Sprintf("%d", kv.Value)))
		w.Write([]byte("\n"))
	}
}
