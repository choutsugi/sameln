package logger

import (
	"LogAgent/common/error"
	"LogAgent/common/record"
	"LogAgent/common/settings"
	"os"

	"go.uber.org/zap"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap/zapcore"
)

var (
	log           *zap.SugaredLogger
	IsInitialized bool
)

// Init 日志模块初始化
func Init(logConfig *settings.LogConfigType, mode string) *error.Error {

	writer := newWriter(
		logConfig.FileName,
		logConfig.MaxSize,
		logConfig.MaxAge,
		logConfig.MaxBackups,
	)

	encoder := newEncoder(logConfig.Type)
	level := new(zapcore.Level)
	if err := level.UnmarshalText([]byte(logConfig.Level)); err != nil {
		record.Failed("日志模块解析级别配置信息失败")
		return error.NewError(err, error.CodeSysLoggerInitFailed)
	}

	var core zapcore.Core
	if mode == modeDevelop {
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		// 同时输出到日志文件和终端
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, writer, level),
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel),
		)
	} else {
		core = zapcore.NewCore(encoder, writer, level)
	}

	logger := zap.New(core, zap.AddCaller())
	// 替换zap包中全局的logger实例
	zap.ReplaceGlobals(logger)
	log = logger.Sugar()
	IsInitialized = true
	return error.Null()
}

// Sync 刷新日志缓存
func Sync() {
	err := log.Sync()
	if err != nil {
		//TODO
		return
	}
}

func L() *zap.SugaredLogger {
	return log
}

func newEncoder(logType string) zapcore.Encoder {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.TimeKey = "time"
	config.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncodeDuration = zapcore.SecondsDurationEncoder
	config.EncodeCaller = zapcore.ShortCallerEncoder

	switch logType {
	case outTypeJson:
		return zapcore.NewJSONEncoder(config)
	case outTypeConsole:
		return zapcore.NewConsoleEncoder(config)
	default:
		return zapcore.NewConsoleEncoder(config)
	}
}

func newWriter(filename string, maxSize, maxAge, maxBackup int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxAge:     maxAge,
		MaxBackups: maxBackup,
	}
	return zapcore.AddSync(lumberJackLogger)
}
