package http_api

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/nsqio/nsq/internal/lg"
)

type logWriter struct {
	logf lg.AppLogFunc
}

func (l logWriter) Write(p []byte) (int, error) {
	l.logf(lg.WARN, "%s", string(p))
	return len(p), nil
}

//Param1监听对象
//Param2路由
//Param3协议类型
//Param4日志函数
func Serve(listener net.Listener, handler http.Handler, proto string, logf lg.AppLogFunc) error {
	logf(lg.INFO, "%s: listening on %s", proto, listener.Addr())
	//创建HTTP对象
	server := &http.Server{
		Handler: handler,
		//错误日志
		ErrorLog: log.New(logWriter{logf}, "", 0),
	}
	err := server.Serve(listener)
	// 排除正常关闭的情况
	// theres no direct way to detect this error because it is not exposed
	if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
		return fmt.Errorf("http.Serve() error - %s", err)
	}

	logf(lg.INFO, "%s: closing %s", proto, listener.Addr())

	return nil
}
