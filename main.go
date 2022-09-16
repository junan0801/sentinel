package main

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Data struct {
	// json 转struct  必须以大写字母开头,否者会被忽略,数据类型也要设置正确,否者会取不到值
	IP              string                     `json:"ip"`
	Port            int                        `json:"port"`
	App             string                     `json:"app"`
	MachineName     string                     `json:"machineName"`
	SuccessQps      float64                    `json:"successQps"`
	PassQps         float64                    `json:"passQps"`
	BlockQps        float64                    `json:"blockQps"`
	ExceptionQps    float64                    `json:"exceptionQps"`
	ResourceDetails map[string]resourceDetails `json:"resourceDetails"`
}
type resourceDetails struct {
	Resource     string  `json:"resource"`
	PassQPS      float64 `json:"passQps"`
	BlockQPS     float64 `json:"blockQps"`
	SuccessQPS   float64 `json:"successQps"`
	ExceptionQPS float64 `json:"exceptionQps"`
}

//定义获取数据函数
func getMetric() (mp map[string][]Data) {
	//domain := os.Getenv("domain")
	domain := "http://g-sentinel-dashboard.tope365.com"

	fmt.Println("qqqqqqqqqqqqqqqqqqqqq")
	url := domain + "/custom/metric/get"
	fmt.Println(url)
	resp, _ := http.Get(url)
	// 因为json 数据key 不固定,使用map 获取
	dataMap := make(map[string][]Data)
	//for range time.Tick(1 * time.Second) {
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	json.Unmarshal([]byte(body), &dataMap)
	//fmt.Printf("获取到的数据类型是%T", dataMap)
	//fmt.Println(dataMap)
	return dataMap
}

type successqpsCollector struct {
	successMetric *prometheus.Desc
	passMetric    *prometheus.Desc
	blockMetric   *prometheus.Desc
}

func newqpsCollector() *successqpsCollector {
	return &successqpsCollector{
		successMetric: prometheus.NewDesc(prometheus.BuildFQName("", "", "successqps"),
			"help",
			[]string{"ip", "port", "app", "machineName", "interface"}, nil),
		passMetric: prometheus.NewDesc(prometheus.BuildFQName("", "", "passqps"),
			"help",
			[]string{"ip", "port", "app", "machineName", "interface"}, nil),
		blockMetric: prometheus.NewDesc(prometheus.BuildFQName("", "", "blockqps"),
			"help",
			[]string{"ip", "port", "app", "machineName", "interface"}, nil),
	}
}

func (collect *successqpsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collect.successMetric
	ch <- collect.passMetric
	ch <- collect.blockMetric
}
func (collect *successqpsCollector) Collect(ch chan<- prometheus.Metric) {
	mp := getMetric()
	//第一个循环取出程序名
	for _, v := range mp {
		//第二个循环取出数据详细信息
		for _, v1 := range v {
			ip := v1.IP
			port := strconv.Itoa(v1.Port)
			app := v1.App
			machinename := v1.MachineName
			//第三个循环取出resource里面的数据
			for _, v2 := range v1.ResourceDetails {
				resource := v2.Resource
				successqps := v2.SuccessQPS
				blockqps := v2.BlockQPS
				passqps := v2.PassQPS
				fmt.Printf("ip %s\t port %s\t machinename %s\t  app %s\t  r %s\t   s%.2f\t ", ip, port, machinename, app, resource, successqps)
				ch <- prometheus.MustNewConstMetric(collect.successMetric, prometheus.GaugeValue, successqps, ip, port, app, machinename, resource)
				ch <- prometheus.MustNewConstMetric(collect.passMetric, prometheus.GaugeValue, passqps, ip, port, app, machinename, resource)
				ch <- prometheus.MustNewConstMetric(collect.blockMetric, prometheus.GaugeValue, blockqps, ip, port, app, machinename, resource)
			}
		}
	}
}

/*type passqpsCollector struct {
	passMetric *prometheus.Desc
}

func newPassqpsCollector() *passqpsCollector {
	m1 := make(map[string]string)
	m1["env"] = "prod"
	v := []string{"ip", "port", "app", "machineName", "interface"}
	return &passqpsCollector{
		passMetric: prometheus.NewDesc("passqps", "passqps of app from sentinel", v, m1),
	}
}*/

/*func (collect *successqpsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collect.passMetric
}

func (collect *successqpsCollector) Collect(ch chan<- prometheus.Metric) {
	//url := "http://g-sentinel-dashboard.tope365.com/custom/metric/get"
	mp := getMetric()
	for _, v := range mp {
		//第二个循环取出数据详细信息
		for _, v1 := range v {
			ip := v1.IP
			port := strconv.Itoa(v1.Port)
			app := v1.App
			machinename := v1.MachineName
			//第三个循环取出resource里面的数据
			for _, v2 := range v1.ResourceDetails {
				resource := v2.Resource
				passqps := v2.PassQPS
				fmt.Printf("ip %s\t port %s\t machinename %s\t  app %s\t  resource %s\t   passqps%.2f\t ", ip, port, machinename, app, resource, passqps)
				ch <- prometheus.MustNewConstMetric(collect.passMetric, prometheus.GaugeValue, passqps, ip, port, app, machinename, resource)
			}
		}
	}
}

type blockqpsCollector struct {
	blockMetric *prometheus.Desc
}

func newBlockqpsCollector() *blockqpsCollector {
	m1 := make(map[string]string)
	m1["env"] = "prod"
	v := []string{"ip", "port", "app", "machineName", "interface"}
	return &blockqpsCollector{
		blockMetric: prometheus.NewDesc("blockqps", "blockqps of app from sentinel", v, m1),
	}
}

func (collect *blockqpsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collect.blockMetric

}
func (collect *blockqpsCollector) Collect(ch chan<- prometheus.Metric) {
	//url := "http://g-sentinel-dashboard.tope365.com/custom/metric/get"
	mp := getMetric()
	for _, v := range mp {
		//第二个循环取出数据详细信息
		for _, v1 := range v {
			ip := v1.IP
			port := strconv.Itoa(v1.Port)
			app := v1.App
			machinename := v1.MachineName
			//第三个循环取出resource里面的数据
			for _, v2 := range v1.ResourceDetails {
				resource := v2.Resource
				blockqps := v2.PassQPS
				fmt.Printf("ip %s\t port %s\t machinename %s\t  app %s\t  resource %s\t   blockqps %.2f\t ", ip, port, machinename, app, resource, blockqps)
				ch <- prometheus.MustNewConstMetric(collect.blockMetric, prometheus.GaugeValue, blockqps, ip, port, app, machinename, resource)
			}
		}
	}

}
*/

func main() {
	//registry := prometheus.NewRegistry()
	qpsRegistry := newqpsCollector()
	//passqpsRegistry := newPassqpsCollector()
	//blockqpsRegistry := newBlockqpsCollector()
	//prometheus.MustRegister(successqpsRegistry, passqpsRegistry, blockqpsRegistry)
	prometheus.MustRegister(qpsRegistry)
	//registry.MustRegister(successqpsRegistry)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":18080", nil)
}
