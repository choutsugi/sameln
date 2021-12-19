package record

import (
	"LogAgent/common/error"
	"LogAgent/common/system"
	"fmt"
)

// 控制台打印级别
const (
	HINT    = "HIN"
	ERROR   = "ERR"
	FAILED  = "FAI"
	SUCCEED = "SUC"
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

	fmt.Printf("TIME: %s\t %s: %s\n", system.LocalTime(), stat, info)
}

func logError(info string, err error.RawErr, stack string) {
	fmt.Printf("TIME: %s\t %s: %s\t RAW_ERR: %s\t CALL_STACK: %s\n", system.LocalTime(), ERROR, info, err, stack)
}

// Hint 控制台打印提示信息
func Hint(args ...interface{}) {
	log(HINT, args...)
}

// Succeed 控制台打印成功信息
func Succeed(args ...interface{}) {
	log(SUCCEED, args...)
}

// Failed 控制台打印失败信息
func Failed(args ...interface{}) {
	log(FAILED, args...)
}

// Error 控制台打印自定义错误
func Error(err *error.Error) {
	logError(err.Info(), err.RawErr(), err.CallStack())
}
