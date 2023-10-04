package metrics

import "container/heap"

type Metrics interface {
	Start()
	IncDomainCount(domain string)
	GetDomainCounts() []KeyValuePair
}

type metrics struct {
	domainHeap MaxHeap
	domainChan chan string
}

func NewMetrics() Metrics {
	m := &metrics{
		domainHeap: MaxHeap{},
		domainChan: make(chan string, 10),
	}
	heap.Init(&m.domainHeap)
	return m
}

func (m *metrics) Start() {
	// TODO: introduce stop channel
	go m.incDomainCounts()
}

func (m *metrics) IncDomainCount(domain string) {
	m.domainChan <- domain
}

func (m *metrics) incDomainCounts() {
	for domain := range m.domainChan {
		m.domainHeap.IncOrPush(domain)
	}
}

func (m *metrics) GetDomainCounts() []KeyValuePair {
	return m.domainHeap.Top3MaxHeapKeysValues()
}
