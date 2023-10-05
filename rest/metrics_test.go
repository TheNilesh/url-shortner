package rest_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/thenilesh/url-shortner/metrics"
	"github.com/thenilesh/url-shortner/mocks"
	"github.com/thenilesh/url-shortner/rest"
)

func TestMetricsHandler_Get(t *testing.T) {
	log := logrus.New()
	m := new(mocks.Metrics)
	c := new(mocks.Collector)
	handler := rest.NewMetricsHandler(log, m)

	m.On("GetCollector", "domain_shortens").Return(c)
	c.On("GetMaxValuePairs", 3).Return([]metrics.KeyValuePair{{Key: "test", Value: 1}})
	req, err := http.NewRequest("GET", "/metrics", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.Get(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "text/plain", rr.Header().Get("Content-Type"))
	assert.Equal(t, "test: 1\n", rr.Body.String())
}
