package collect

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

//var procPath = "/proc/loadavg"
var procPath = "d:/test.txt"

// 获取负载

func GetLoad() (loads []float64, err error) {
	// 根据文件位置读取到数据切片
	data, err := ioutil.ReadFile(procPath)
	if err != nil {
		return nil, err
	}
	// 调用下面的parseLoad方法,返回负载
	loads, err = parseLoad(string(data))
	if err != nil {
		return nil, err
	}
	return loads, nil

}

// 具体的解析方法,接收一个切片数据,处理后返回负载和是否error
func parseLoad(data string) (loads []float64, err error) {
	// 创建一个长度为3的floatt64 类型的空切片
	loads = make([]float64, 3)
	//0.00 0.00 0.00 2/167 30355
	//将数据分割(默认就是空格分割)
	parts := strings.Fields(data)
	// load 分别显示1分钟 5分钟 15分钟,所以必须不少于3条
	if len(parts) < 3 {
		return nil, fmt.Errorf("unexpected content in %s", procPath)
	}
	// range 循环取数据
	for i, load := range parts[0:3] {
		loads[i], err = strconv.ParseFloat(load, 64)
		if err != nil {
			return nil, fmt.Errorf("\"could not parse load '%s': %w\", load, err")
		}
	}
	return loads, nil
}
