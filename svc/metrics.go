package svc

import "github.com/prometheus/client_golang/prometheus"

type metrics struct {
	domainCounts *prometheus.CounterVec
}

func NewMetrics() *metrics {
	m := metrics{
		domainCounts: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "url_shortner",
				Name:      "shortened_domain_count",
				Help:      "Counts the number of times each domain is shortened.",
			},
			[]string{"domain"},
		),
	}
	defaultRegistry := prometheus.NewRegistry()
	prometheus.DefaultRegisterer = defaultRegistry
	prometheus.DefaultGatherer = defaultRegistry
	prometheus.MustRegister(m.domainCounts)
	return &m
}

func (m *metrics) IncDomainCount(domain string) {
	m.domainCounts.With(prometheus.Labels{"domain": domain}).Inc()
}
