package http_api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/nsqio/nsq/internal/lg"
)

//装饰器使用的函数
//将原有的APIHandler外层装饰一层其他逻辑返回装饰后的APIHandler
type Decorator func(APIHandler) APIHandler

//api逻辑处理函数的定义（httprouter 的回调格式）
type APIHandler func(http.ResponseWriter, *http.Request, httprouter.Params) (interface{}, error)

//记录错误的结构体
type Err struct {
	Code int
	Text string
}

//返回错误字符
func (e Err) Error() string {
	return e.Text
}

//版本检测
func acceptVersion(req *http.Request) int {
	if req.Header.Get("accept") == "application/vnd.nsq; version=1.0" {
		return 1
	}

	return 0
}

//直接把处理结果打印出来不通过协议包装（文本格式）
func PlainText(f APIHandler) APIHandler {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
		code := 200
		data, err := f(w, req, ps)
		if err != nil {
			code = err.(Err).Code
			data = err.Error()
		}
		switch d := data.(type) {
		case string:
			w.WriteHeader(code)
			io.WriteString(w, d)
		case []byte:
			w.WriteHeader(code)
			w.Write(d)
		default:
			panic(fmt.Sprintf("unknown response type %T", data))
		}
		return nil, nil
	}
}

//返回V1版本的协议处理函数（json格式）
func V1(f APIHandler) APIHandler {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
		//调用内层的逻辑
		data, err := f(w, req, ps)
		if err != nil {
			RespondV1(w, err.(Err).Code, err)
			return nil, nil
		}
		RespondV1(w, 200, data)
		return nil, nil
	}
}

//具体的逻辑执行
func RespondV1(w http.ResponseWriter, code int, data interface{}) {
	var response []byte
	var err error
	var isJSON bool

	if code == 200 {
		switch data.(type) {
		case string:
			response = []byte(data.(string))
		case []byte:
			response = data.([]byte)
		case nil:
			response = []byte{}
		//默认使用json打印错误消息
		default:
			isJSON = true
			response, err = json.Marshal(data)
			if err != nil {
				code = 500
				data = err
			}
		}
	}
	//错误情况下打印错误消息
	if code != 200 {
		isJSON = true
		response, _ = json.Marshal(struct {
			Message string `json:"message"`
		}{fmt.Sprintf("%s", data)})
	}

	//标记内容是json格式
	if isJSON {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
	//写入请求版本
	w.Header().Set("X-NSQ-Content-Type", "nsq; version=1.0")
	w.WriteHeader(code)
	w.Write(response)
}

//返回层层装饰后的匿名函数（httprouter的回调函数）
func Decorate(f APIHandler, ds ...Decorator) httprouter.Handle {
	decorated := f
	for _, decorate := range ds {
		//每次遍历都是将原来的APIHandler外面再装饰一层，最后的调用顺序FILO
		decorated = decorate(decorated)
	}
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		//调用最后封装的结果
		decorated(w, req, ps)
	}
}

//返回log的装饰函数
//嵌套了2层
//第一次调用返回一个能生成装饰函数
//再次调用返回的APIHandler返回真正日志的函数
func Log(logf lg.AppLogFunc) Decorator {
	return func(f APIHandler) APIHandler {
		return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
			//耗时计算
			start := time.Now()
			//调用内层的逻辑
			response, err := f(w, req, ps)
			elapsed := time.Since(start)
			status := 200
			if e, ok := err.(Err); ok {
				status = e.Code
			}
			logf(lg.INFO, "%d %s %s (%s) %s",
				status, req.Method, req.URL.RequestURI(), req.RemoteAddr, elapsed)
			return response, err
		}
	}
}

//返回错误处理的匿名函数
func LogPanicHandler(logf lg.AppLogFunc) func(w http.ResponseWriter, req *http.Request, p interface{}) {
	return func(w http.ResponseWriter, req *http.Request, p interface{}) {
		logf(lg.ERROR, "panic in HTTP handler - %s", p)
		Decorate(func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
			return nil, Err{500, "INTERNAL_ERROR"}
		}, Log(logf), V1)(w, req, nil)
	}
}

//返回处理找不到路由的匿名函数
func LogNotFoundHandler(logf lg.AppLogFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		Decorate(func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
			return nil, Err{404, "NOT_FOUND"}
		}, Log(logf), V1)(w, req, nil)
	})
}

//返回处理不允许方法的匿名
func LogMethodNotAllowedHandler(logf lg.AppLogFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		Decorate(func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
			return nil, Err{405, "METHOD_NOT_ALLOWED"}
		}, Log(logf), V1)(w, req, nil)
	})
}
