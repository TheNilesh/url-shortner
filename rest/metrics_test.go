// BEGIN: q7r6t5y2plo9
package rest_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thenilesh/url-shortner/metrics"
	"github.com/thenilesh/url-shortner/rest"
)

type mockMetrics struct {
	mock.Mock
}

func (m *mockMetrics) Start() {
	m.Called()
}

func (m *mockMetrics) IncDomainCount(domain string) {
	m.Called(domain)
}

func (m *mockMetrics) GetDomainCounts() []metrics.KeyValuePair {
	args := m.Called()
	return args.Get(0).([]metrics.KeyValuePair)
}

func TestMetricsHandler_Get(t *testing.T) {
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)
	mockMetrics := new(mockMetrics)
	handler := rest.NewMetricsHandler(log, mockMetrics)

	req, _ := http.NewRequest("GET", "/metrics", nil)
	expectedDomainCounts := []metrics.KeyValuePair{
		{
			Key:   "example.com",
			Value: 2,
		},
		{
			Key:   "example.org",
			Value: 1,
		},
	}
	mockMetrics.On("GetDomainCounts").Return(expectedDomainCounts, nil)

	rr := httptest.NewRecorder()
	handler.Get(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "text/plain", rr.Header().Get("Content-Type"))
	assert.Equal(t, "example.com: 2\nexample.org: 1\n", rr.Body.String())
}
