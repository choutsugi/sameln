// Package logger 日志模块，基于zap库的日志器
package logger

import (
	"LogAgent/universal/error"
	"LogAgent/universal/generic"
	"LogAgent/universal/record"
	"LogAgent/universal/settings"
	"go.uber.org/atomic"
	"os"
	"time"

	"go.uber.org/zap"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap/zapcore"
)

var (
	log         *zap.SugaredLogger
	initialized atomic.Bool
)

// Init 日志模块初始化
func Init(config *settings.LogConfigType, mode string) *error.Error {
	if initialized.Load() {
		return error.Null()
	}

	level := new(zapcore.Level)
	if err := level.UnmarshalText([]byte(config.Level)); err != nil {
		record.Fatal("日志模块解析级别配置信息失败")
		return error.NewError(err, error.CodeSysLoggerInitFailed)
	}

	writeSyncer := createWriteSyncer(config.FileName, config.MaxSize, config.MaxAge, config.MaxBackups)
	fileEncoder := createEncoder(config.Type)

	var core zapcore.Core
	if mode == modeDevelop {
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		core = zapcore.NewTee(
			zapcore.NewCore(fileEncoder, writeSyncer, level),
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel),
		)
	} else {
		core = zapcore.NewCore(fileEncoder, writeSyncer, level)
	}

	logger := zap.New(core, zap.AddCaller())
	// 替换zap包中全局的logger实例
	zap.ReplaceGlobals(logger)
	log = logger.Sugar()

	initialized.Store(true)
	return error.Null()
}

// L 日志器
func L() *zap.SugaredLogger {
	return log
}

// Sync 刷新日志缓存
func Sync() {
	for tick := 0; tick < generic.TrySyncWithMaxTime; tick++ {
		if raw := log.Sync(); raw == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// IsInitialized 是否已初始化
func IsInitialized() bool {
	return initialized.Load()
}

// 创建日志写同步器
func createWriteSyncer(filename string, maxSize, maxAge, maxBackup int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxAge:     maxAge,
		MaxBackups: maxBackup,
	}
	return zapcore.AddSync(lumberJackLogger)
}

// 创建日志编码器
func createEncoder(logType string) zapcore.Encoder {
	config := zapcore.EncoderConfig{
		MessageKey:     "MSG",
		LevelKey:       "LEVEL",
		TimeKey:        "TIME",
		NameKey:        "NAME",
		CallerKey:      "CALLER",
		FunctionKey:    "FUNC",
		StacktraceKey:  "STACKTRACE",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     createTimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	// 根据配置设置日志输出位置
	switch logType {
	case outTypeJson:
		return zapcore.NewJSONEncoder(config)
	case outTypeConsole:
		return zapcore.NewConsoleEncoder(config)
	default:
		return zapcore.NewConsoleEncoder(config)
	}
}

// 创建日志时间格式编码器：自定义
func createTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}
