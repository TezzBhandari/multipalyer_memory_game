package metrics

type Metrics struct {
	metrics map[string]int
}

func NewMetrics() *Metrics {
	return &Metrics{
		metrics: make(map[string]int),
	}
}

func (m *Metrics) registerMetrics(key string) {
	m.metrics[key] = 0
}

func (m *Metrics) SetIfGreater(key string, value int) {
	if value > m.metrics[key] {
		m.metrics[key] = value
	}
}

func (m *Metrics) Inc(key string) {
	m.metrics[key]++
}

func (m *Metrics) Get(key string) int {
	return m.metrics[key]
}

func (m *Metrics) GetAll() map[string]int {
	// return a copy to avoid concurrent map access, this should be more performant than using a sync.Map
	metricsCopy := make(map[string]int)
	for k, v := range m.metrics {
		metricsCopy[k] = v
	}
	return metricsCopy
}
