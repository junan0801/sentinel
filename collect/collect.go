package collect

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
)

var namespace = "node"

type loadavagCollector struct {
	metrics []typeDesc
}
type typeDesc struct {
	desc      *prometheus.Desc
	valueType prometheus.ValueType
}

func NewloadavagCollector() *loadavagCollector {
	return &loadavagCollector{
		metrics: []typeDesc{
			{prometheus.NewDesc(namespace+"_load1", "1m load average", nil, nil), prometheus.GaugeValue},
			{prometheus.NewDesc(namespace+"_load5", "5m load average.", nil, nil), prometheus.GaugeValue},
			{prometheus.NewDesc(namespace+"_load15", "15m load average.", nil, nil), prometheus.GaugeValue},
		},
	}

}
func (collector *loadavagCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.metrics[1].desc
}
func (collector *loadavagCollector) Collect(ch chan<- prometheus.Metric) {
	loads, err := GetLoad()
	if err != nil {
		log.Print("get loadavag error:", err)
	}
	for i, load := range loads {
		ch <- prometheus.MustNewConstMetric(collector.metrics[i].desc, prometheus.GaugeValue, load)
	}
}
