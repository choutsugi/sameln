// Package record 日志模块初始化之前打印信息
package record

import (
	"LogAgent/universal/error"
	"LogAgent/universal/system"
	"fmt"
)

// 控制台打印级别
const (
	INFO  = "INFO"
	ERROR = "ERROR"
	FATAL = "FATAL"
	WARN  = "WARN"
)

func log(stat string, args ...interface{}) {
	var info string
	var isString bool

	length := len(args)

	if length > 1 {
		info = fmt.Sprintf(args[0].(string), args[1:]...)
	} else if len(args) == 1 {
		info, isString = args[0].(string)
		if !isString {
			return
		}
	} else {
		return
	}

	fmt.Printf("%s\t%s\t%s\n", system.UtcTime(), stat, info)
}

func logError(info string, err error.RawErr, stack string) {
	fmt.Printf("TIME: %s\t %s: %s\t RAW_ERR: %s\t CALL_STACK: %s\n", system.LocalTime(), ERROR, info, err, stack)
}

// Info 控制台打印提示信息
func Info(args ...interface{}) {
	log(INFO, args...)
}

// Warn 控制台打印警告信息
func Warn(args ...interface{}) {
	log(WARN, args...)
}

// Fatal 控制台打印错误信息
func Fatal(args ...interface{}) {
	log(FATAL, args...)
}

// Error 控制台打印自定义错误
func Error(err *error.Error) {
	logError(err.Info(), err.RawErr(), err.CallStack())
}
