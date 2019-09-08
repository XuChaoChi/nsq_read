//nsq对标准库日志的简单封装
// short for "log"
package lg

import (
	"fmt"
	"log"
	"os"
	"strings"
)

//日志等级
const (
	DEBUG = LogLevel(1)
	INFO  = LogLevel(2)
	WARN  = LogLevel(3)
	ERROR = LogLevel(4)
	FATAL = LogLevel(5)
)

//定义用者使用日志函数的标准
type AppLogFunc func(lvl LogLevel, f string, args ...interface{})

//日志接口
type Logger interface {
	Output(maxdepth int, s string) error
}

//定义一个空的日志
type NilLogger struct{}

//空日志必须实现的接口
func (l NilLogger) Output(maxdepth int, s string) error {
	return nil
}

//定义日志等级类型
type LogLevel int

//获取日志等级
func (l *LogLevel) Get() interface{} { return *l }

//设置日志等级
func (l *LogLevel) Set(s string) error {
	lvl, err := ParseLogLevel(s)
	if err != nil {
		return err
	}
	*l = lvl
	return nil
}

//获取当前日志等级
func (l *LogLevel) String() string {
	switch *l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARNING"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	}
	return "invalid"
}

//通过字符串获取日志等级
func ParseLogLevel(levelstr string) (LogLevel, error) {
	switch strings.ToLower(levelstr) {
	case "debug":
		return DEBUG, nil
	case "info":
		return INFO, nil
	case "warn":
		return WARN, nil
	case "error":
		return ERROR, nil
	case "fatal":
		return FATAL, nil
	}
	return 0, fmt.Errorf("invalid log level '%s' (debug, info, warn, error, fatal)", levelstr)
}

//日志调用
//Param1 log对象
//Param2 配置的级别
//Param3 当前消息的等级
func Logf(logger Logger, cfgLevel LogLevel, msgLevel LogLevel, f string, args ...interface{}) {
	//消息级别没有配置高直接返回
	if cfgLevel > msgLevel {
		return
	}
	logger.Output(3, fmt.Sprintf(msgLevel.String()+": "+f, args...))
}

//在没有日志对象的情况下输出严:重错误日志
func LogFatal(prefix string, f string, args ...interface{}) {
	logger := log.New(os.Stderr, prefix, log.Ldate|log.Ltime|log.Lmicroseconds)
	Logf(logger, FATAL, FATAL, f, args...)
	os.Exit(1)
}
