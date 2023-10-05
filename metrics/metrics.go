package metrics

type Metrics interface {
	Start()
	GetCollector(name string) Collector
}

type metrics struct {
	collectors map[string]Collector
}

func NewMetrics() Metrics {
	m := &metrics{
		// TODO: Externalize registration of collectors
		collectors: map[string]Collector{
			"domain_shortens": newCollector(10),
		},
	}
	return m
}

func (m *metrics) Start() {
	for _, c := range m.collectors {
		c.Start()
	}
}

func (m *metrics) GetCollector(name string) Collector {
	return m.collectors[name]
}

type Collector interface {
	Start()
	Inc(key string)
	GetMaxValuePairs(n int) []KeyValuePair
}

type collector struct {
	bufferChan chan string
	heap       Heap
}

func newCollector(bufferSize int) Collector {
	c := &collector{
		bufferChan: make(chan string, bufferSize),
		heap:       NewHeap(),
	}
	return c
}

func (c *collector) Start() {
	// TODO: introduce stop channel
	go func() {
		for domain := range c.bufferChan {
			c.heap.IncOrPush(domain)
		}
	}()
}

func (c *collector) Inc(key string) {
	c.heap.IncOrPush(key)
}

func (c collector) GetMaxValuePairs(n int) []KeyValuePair {
	return c.heap.GetMaxValuePairs(n)
}
