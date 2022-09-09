package main

/*import (
	"export/collect"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func main() {
	prometheus.MustRegister(collect.NewloadavagCollector())
	http.Handle("/metrics", promhttp.Handler())
	log.Print("expose /metrics use port:8085")
	log.Fatal(http.ListenAndServe(":8085", nil))
}
*/
/*import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"time"
)

func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}

var (
	opsProcessed = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "myapp",
		Name:      "processed_ops_total",
		Help:      "The total number",
	})
)

func main() {
	prometheus.MustRegister(opsProcessed)
	recordMetrics()

	http.Handle("/metrics1", promhttp.Handler())
	log.Print("export /metrics on port 8085")
	http.ListenAndServe(":8086", nil)

}
*/

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics() {
	go func() {
		for {
			opsProcessed1.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events  rrrrrrrrr",
	})
	opsProcessed1 = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "My_example_guage_data",
		Help:        "My example guage data",
		ConstLabels: map[string]string{"error": "111"},
	})
)

func main() {
	recordMetrics()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
