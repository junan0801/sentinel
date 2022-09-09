package main

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io/ioutil"
	"net/http"
)

type Data struct {
	// json 转struct  必须以大写字母开头,否者会被忽略,数据类型也要设置正确,否者会取不到值
	IP          string  `json:"ip"`
	App         string  `json:"app"`
	SuccessQps  float64 `json:"successQps"`
	MachineName string  `json:"machineName"`
}

func main() {
	registry := prometheus.NewRegistry()
	queueLength := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		//以下两个参数会组合成 monitoring_demo_sentinel
		//Namespace: "monitoring",
		//Subsystem: "demo",
		Name: "sentinel",
		Help: "successqps of app from sentinel",
		//静态标签
		//ConstLabels: map[string]string{
		//	"module": "http-server",
	}, []string{"ip", "app", "machinename"})

	url := "http://g-sentinel-dashboard.tope365.com/custom/metric/get"
	resp, _ := http.Get(url)
	// 因为json 数据key 不固定,使用map 获取
	dataMap := make(map[string][]Data)
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	json.Unmarshal([]byte(body), &dataMap)
	fmt.Printf("获取到的数据类型是%T", dataMap)
	//m := make(map[string]string)
	for _, v := range dataMap {
		//log.Printf(`%v`, v)
		for _, v1 := range v {
			ip := v1.IP
			app := v1.App
			successqps := v1.SuccessQps
			machinename := v1.MachineName
			/*			m["ip"] = ip
						m["app"] = app
						m["successqps"] = strconv.Itoa(successqps)
						m["machinename"] = machinename*/
			//fmt.Println(m)
			//fmt.Println(v1.IP, v1)
			queueLength.WithLabelValues(ip, app, machinename).Set(successqps)
		}

	}
	//queueLength.WithLabelValues("127.0.0.1", "11", "TESTMACHINE").Set(100)

	registry.MustRegister(queueLength)
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{Registry: registry}))
	http.ListenAndServe(":8050", nil)
}
