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
	url := "http://g-sentinel-dashboard.tope365.com/custom/metric/get"
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
	fooMetric *prometheus.Desc
}

func newSuccessqpsCollector() *successqpsCollector {
	m1 := make(map[string]string)
	m1["env"] = "prod"
	v := []string{"ip", "port", "app", "machineName", "interface"}
	return &successqpsCollector{
		fooMetric: prometheus.NewDesc("successqps", "successqps of app from sentinel", v, m1),
	}
}

func (collect *successqpsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collect.fooMetric
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
			fmt.Println("aaaaaaaaaaaaaaaaaaaaaa")
			fmt.Println(v1)
			fmt.Println("111111111111111111111111")
			fmt.Println(v1.ResourceDetails)
			//第三个循环取出resource里面的数据
			for _, v2 := range v1.ResourceDetails {
				resource := v2.Resource
				successqps := v2.SuccessQPS
				//success := v2.SuccessQPS
				//s := strconv.FormatFloat(success,'g',5, 32)
				fmt.Println("bbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
				fmt.Printf("ip %s\t port %s\t machinename %s\t  app %s\t  r %s\t   s%d\t ", ip, port, machinename, app, resource, successqps)
				ch <- prometheus.MustNewConstMetric(collect.fooMetric, prometheus.GaugeValue, successqps, ip, port, app, machinename, resource)
			}
		}
	}
}

type passqpsCollector struct {
	fooMetric *prometheus.Desc
}

func newPassqpsCollector() *passqpsCollector {
	m1 := make(map[string]string)
	m1["env"] = "prod"
	v := []string{"ip", "port", "app", "machineName"}
	return &passqpsCollector{
		fooMetric: prometheus.NewDesc("passqps", "passqps of app from sentinel", v, m1),
	}
}

func (collect *passqpsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collect.fooMetric
}

func (collect *passqpsCollector) Collect(ch chan<- prometheus.Metric) {
	//url := "http://g-sentinel-dashboard.tope365.com/custom/metric/get"
	mp := getMetric()
	for _, v := range mp {
		for _, v1 := range v {
			ip := v1.IP
			port := strconv.Itoa(v1.Port)
			app := v1.App

			machinename := v1.MachineName
			passqps := v1.PassQps
			fmt.Println(ip, app, machinename, passqps)
			ch <- prometheus.MustNewConstMetric(collect.fooMetric, prometheus.GaugeValue, passqps, ip, port, app, machinename)
		}
	}
}

type blockqpsCollector struct {
	fooMetric *prometheus.Desc
}

func newBlockqpsCollector() *blockqpsCollector {
	m1 := make(map[string]string)
	m1["env"] = "prod"
	v := []string{"ip", "port", "app", "machineName"}
	return &blockqpsCollector{
		fooMetric: prometheus.NewDesc("blockqps", "blockqps of app from sentinel", v, m1),
	}
}

func (collect *blockqpsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collect.fooMetric

}
func (collect *blockqpsCollector) Collect(ch chan<- prometheus.Metric) {
	//url := "http://g-sentinel-dashboard.tope365.com/custom/metric/get"
	mp := getMetric()
	for _, v := range mp {
		for _, v1 := range v {
			ip := v1.IP
			port := strconv.Itoa(v1.Port)
			app := v1.App
			machinename := v1.MachineName
			blockqps := v1.BlockQps
			fmt.Println(ip, app, machinename, blockqps)
			ch <- prometheus.MustNewConstMetric(collect.fooMetric, prometheus.GaugeValue, blockqps, ip, port, app, machinename)
		}
	}

}
func main() {
	successqpsRegistry := newSuccessqpsCollector()
	/*	passqpsRegistry := newPassqpsCollector()
		blockqpsRegistry := newBlockqpsCollector()*/
	//prometheus.MustRegister(successqpsRegistry, passqpsRegistry, blockqpsRegistry)
	prometheus.MustRegister(successqpsRegistry)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":18080", nil)
}
