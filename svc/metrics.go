package svc

type Metrics struct {
	domainCounts map[string]int64
	domainChan   chan string
}

func NewMetrics() *Metrics {
	m := Metrics{
		domainCounts: map[string]int64{},
		domainChan:   make(chan string, 10),
	}
	return &m
}

func (m *Metrics) Start() {
	go m.incDomainCounts()
}

func (m *Metrics) IncDomainCount(domain string) {
	m.domainChan <- domain
}

func (m *Metrics) incDomainCounts() {
	for domain := range m.domainChan {
		if _, ok := m.domainCounts[domain]; !ok {
			m.domainCounts[domain] = 0
		}
		m.domainCounts[domain]++
		// TODO: Keep only 3 domains with highest counts
	}
}

func (m *Metrics) GetDomainCounts() map[string]int64 {
	// TODO: Copy map to another map and return for immutability
	return m.domainCounts
}
