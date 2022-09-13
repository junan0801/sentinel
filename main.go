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

/*import (
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
}*/

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io/ioutil"
	"net/http"
)

type Data2 struct {
	// json 转struct  必须以大写字母开头,否者会被忽略,数据类型也要设置正确,否者会取不到值
	IP          string  `json:"ip"`
	App         string  `json:"app"`
	SuccessQps  float64 `json:"successQps"`
	MachineName string  `json:"machineName"`
}

type fooCollector struct {
	fooMetric *prometheus.Desc
	barMetric *prometheus.Desc
}

func newFoolCollector() *fooCollector {
	m1 := make(map[string]string)
	m1["env"] = "prod"
	v := []string{"ip", "app", "machineName"}
	return &fooCollector{
		fooMetric: prometheus.NewDesc("FF_METRICS", "SHOW metrics a for", nil, nil),
		barMetric: prometheus.NewDesc("sentinel", "successqps of app from sentinel", v, m1),
	}

}
func (collect *fooCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collect.barMetric
	ch <- collect.fooMetric

}
func (collect *fooCollector) Collect(ch chan<- prometheus.Metric) {
	url := "http://g-sentinel-dashboard.tope365.com/custom/metric/get"
	resp, _ := http.Get(url)
	// 因为json 数据key 不固定,使用map 获取
	dataMap := make(map[string][]Data2)
	//for range time.Tick(1 * time.Second) {
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	json.Unmarshal([]byte(body), &dataMap)
	fmt.Printf("获取到的数据类型是%T\t", dataMap)
	//fmt.Println(dataMap)
	for _, v := range dataMap {
		for _, v1 := range v {
			ip := v1.IP
			app := v1.App
			//successqps := strconv.FormatFloat(v1.SuccessQps, 'g', 2, 32)
			successqps := v1.SuccessQps
			machinename := v1.MachineName
			fmt.Println(ip, app, successqps, machinename)
			ch <- prometheus.MustNewConstMetric(collect.barMetric, prometheus.GaugeValue, successqps, ip, app, machinename)
		}
	}

}
func main() {
	foo := newFoolCollector()
	prometheus.MustRegister(foo)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":18080", nil)
}
