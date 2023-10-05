package metrics

import (
	"testing"
)

func TestMetrics(t *testing.T) {
	m := NewMetrics()
	m.Start()

	collector := m.GetCollector("domain_shortens")
	if collector == nil {
		t.Errorf("Expected collector to not be nil")
	}

	collector.Inc("example.com")
	collector.Inc("example.com")
	collector.Inc("example.org")
	collector.Inc("example.net")

	pairs := collector.GetMaxValuePairs(2)
	if len(pairs) != 2 {
		t.Errorf("Expected 2 pairs, got %d", len(pairs))
	}
	if pairs[0].Key != "example.com" || pairs[0].Value != 2 {
		t.Errorf("Expected pair 0 to be {\"example.com\", 2}, got %v", pairs[0])
	}
	if pairs[1].Key != "example.org" || pairs[1].Value != 1 {
		t.Errorf("Expected pair 1 to be {\"example.org\", 1}, got %v", pairs[1])
	}
}
