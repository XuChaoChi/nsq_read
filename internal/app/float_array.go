package app

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
)

type FloatArray []float64

//获取整个数组
func (a *FloatArray) Get() interface{} { return []float64(*a) }

//使用时候用，来分割
func (a *FloatArray) Set(param string) error {
	for _, s := range strings.Split(param, ",") {
		//将string解析成float32
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			log.Fatalf("Could not parse: %s", s)
			return nil
		}
		*a = append(*a, v)
	}
	//每次添加后都排下序
	sort.Sort(*a)
	return nil
}

//sort接口的实现
func (a FloatArray) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a FloatArray) Less(i, j int) bool { return a[i] > a[j] }
func (a FloatArray) Len() int           { return len(a) }

//将数组内容转成string返回
func (a *FloatArray) String() string {
	var s []string
	for _, v := range *a {
		s = append(s, fmt.Sprintf("%f", v))
	}
	return strings.Join(s, ",")
}
