package error

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"time"
)

type Error struct {
	StatusCode uint64    // 错误码
	Message    string    // 错误信息
	rawErr     error     // 原生错误
	callStack  []uintptr // 函数调用栈指针
}

var (
	null = &Error{
		StatusCode: CodeSysSuccess,
		Message:    StatusSuccess,
		rawErr:     nil,
		callStack:  nil,
	}
)

// Info 返回自定义错误信息
func (err *Error) Info() string {
	return err.Message
}

// RawErr 返回原生error
func (err *Error) RawErr() error {
	return err.rawErr
}

// CallStack 获取函数调用栈
func (err *Error) CallStack() string {
	frames := runtime.CallersFrames(err.callStack)
	var (
		f      runtime.Frame
		more   bool
		result string
		index  int
	)

	for {
		f, more = frames.Next()
		if index = strings.Index(f.File, "src"); index != -1 {
			f.File = f.File[index+4:]
		}
		result = fmt.Sprintf("%s%s\n\t%s:%d\n", result, f.Function, f.File, f.Line)
		if !more {
			break
		}
	}
	return result
}

func GetInfo(code uint64) string {
	message, isExist := errMsg[code]
	if !isExist {
		return ""
	}
	return message
}

// NewError 错误打包
func NewError(raw error, code uint64) *Error {
	pcs := make([]uintptr, 32)
	// skip the first 3 invocations
	count := runtime.Callers(3, pcs)
	return &Error{
		Message:    errMsg[code],
		StatusCode: code,
		rawErr:     raw,
		callStack:  pcs[:count],
	}
}

// NewCustomError 自定义错误打包
func NewCustomError(raw error, args ...interface{}) *Error {
	message := fmtErrMsg(args)
	if raw == nil {
		raw = errors.New(message)
	}

	if message == "" {
		message = raw.Error()
	}

	pcs := make([]uintptr, 32)
	// skip the first 3 invocations
	count := runtime.Callers(3, pcs)

	return &Error{
		Message:    message,
		StatusCode: 0,
		rawErr:     raw,
		callStack:  pcs[:count],
	}
}

// Null 空错误
func Null() *Error {
	return null
}

// NullWithCode 带状态码的空错误
func NullWithCode(code uint64) *Error {
	message, isExist := errMsg[code]
	if !isExist {
		message = StatusUnknown
	}
	null.Message = message
	return null
}

// NullWithMessage 带有自定义信息的空错误
func NullWithMessage(args ...interface{}) *Error {
	message := fmtErrMsg(args)
	if message != "" {
		null.Message = message
	}
	return null
}

func LogError(err Error) {
	fmt.Printf("[E%d] info: %s\nraw err: %s\ncall stack: %s\n",
		err.StatusCode,
		err.Info(),
		err.RawErr(),
		err.CallStack(),
	)
}

func LogErrorWithTime(err Error) {
	fmt.Printf("[E%d] time:%v info: %s\nraw err: %s\ncall stack: %s\n",
		err.StatusCode,
		time.Now().Local(),
		err.Info(),
		err.RawErr(),
		err.CallStack(),
	)
}

// fmtErrMsg 格式化自定义错误信息
func fmtErrMsg(args ...interface{}) string {
	if len(args) > 1 {
		return fmt.Sprintf(args[0].(string), args[1:]...)
	}
	if len(args) == 1 {
		if v, ok := args[0].(string); ok {
			return v
		}
		if v, ok := args[0].(error); ok {
			return v.Error()
		}
	}
	return ""
}
