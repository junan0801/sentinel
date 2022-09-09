package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

//定义api接口数据结构和序列化json字段
/*type Data struct {
	ID    int    `json:"id"`
	IP    string `json:"IP"`
	DESC  string `json:"desc"`
	OWNER string `json:"owner"`
}

type CloumnsData struct {
	NAME  string `json:"name"`
	ALIAS string `json:"alias"`
}

type Employee struct {
	CODE    int           `json:"code"`
	ISADMIN bool          `json:"isadmin"`
	DATA    []Data        `json:"data"`
	COLUMNS []CloumnsData `json:"columns"`
	MESSGAE string        `json:"messgae"`
}*/

type Data struct {
	IP          string `json:"ip"`
	App         string `json:"app"`
	SuccessQps  int    `json:"successQps"`
	MachineName string `json:"machineName"`
}

//type OauthCheckToken struct {
//	Resource     string `json:"resource"`
//	PassQPS      int    `json:"passQps"`
//	BlockQPS     int    `json:"blockQps"`
//	SuccessQPS   int    `json:"successQps"`
//	ExceptionQPS int    `json:"exceptionQps"`
//}
//type ResourceDetails struct {
//	OauthCheckToken OauthCheckToken `json:"/oauth/check_token"`
//}
//
//type sentinel struct {
//	Topeauthv2 []Data `json:"tope-auth-v2"`
//}

//发送http请求和json 序列化并打印数据结构;

//func main() {
//	url := "http://g-sentinel-dashboard.tope365.com/custom/metric/get"
//	resp, _ := http.Get(url)
//	s := sentinel{}
//	body, _ := ioutil.ReadAll(resp.Body)
//	resp.Body.Close()
//	ss := json.Unmarshal([]byte(body), &s)
//	fmt.Println(ss)
//	fmt.Println(fmt.Sprintf("%+v", s))
//
//}

/*func getMetrics(url string) (map[string][]) } {
	//url := "http://g-sentinel-dashboard.tope365.com/custom/metric/get"
	resp, _ := http.Get(url)
	// 因为json 数据key 不固定,使用map 获取
	dataMap := make(map[string][]Data)
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	json.Unmarshal([]byte(body), &dataMap)
	fmt.Printf("获取到的数据类型是%T", dataMap)
	for k, v := range dataMap {
		log.Printf(`%v,%+v`, k, v)
	}
	//m := map(dataMap)
	return dataMap
}*/
func main() {
	//getMetrics("http://g-sentinel-dashboard.tope365.com/custom/metric/get")
	url := "http://g-sentinel-dashboard.tope365.com/custom/metric/get"
	resp, _ := http.Get(url)
	// 因为json 数据key 不固定,使用map 获取
	dataMap := make(map[string][]Data)
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	json.Unmarshal([]byte(body), &dataMap)
	fmt.Printf("获取到的数据类型是%T", dataMap)
	m := make(map[string]string)
	for _, v := range dataMap {
		//log.Printf(`%v`, v)
		for _, v1 := range v {
			ip := v1.IP
			app := v1.App
			successqps := v1.SuccessQps
			machinename := v1.MachineName
			m["ip"] = ip
			m["app"] = app
			m["successqps"] = strconv.Itoa(successqps)
			m["machinename"] = machinename
			fmt.Println(m)
			//fmt.Println(v1.IP, v1)
		}

		//m[string(k1)] = int(v1)
		//log.Printf("%v,%+v", k1, v1)
		//fmt.Printf("%T %T %T", k1, v1, m)
		//fmt.Println()
		//time.Sleep(100)*/
		//fmt.Printf("%T", v)

		//first := v.Data[0]
		//fmt.Println(first.IP)

	}

}
