package svc

import "container/heap"

type Metrics struct {
	domainHeap MaxHeap
	domainChan chan string
}

func NewMetrics() *Metrics {
	m := Metrics{
		domainHeap: MaxHeap{},
		domainChan: make(chan string, 10),
	}
	heap.Init(&m.domainHeap)
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
		m.domainHeap.IncOrPush(domain)
	}
}

func (m *Metrics) GetDomainCounts() []KeyValuePair {
	return m.domainHeap.Top3MaxHeapKeysValues()
}
