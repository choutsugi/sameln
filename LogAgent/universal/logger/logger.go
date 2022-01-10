// Package logger
package logger

import (
	"LogAgent/universal/codes"
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

// Init Initialize the log module
func Init(config *settings.LogConfigType, mode string) *error.Error {
	if initialized.Load() {
		L().Error("The Logger module unable to re-initialize!")
		return error.NewError(nil, codes.InitLoggerFailed)
	}

	level := new(zapcore.Level)
	if raw := level.UnmarshalText([]byte(config.Level)); raw != nil {
		record.Warn("The Logger module unmarshal the running level unsuccessfully, please check the configuration file(%s).", config.FileName)
		return error.NewError(raw, codes.InitLoggerFailed)
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
	zap.ReplaceGlobals(logger)
	log = logger.Sugar()

	initialized.Store(true)
	return error.Null()
}

func L() *zap.SugaredLogger {
	return log
}

// Sync flush log buffer
func Sync() {
	for tick := 0; tick < generic.TrySyncWithMaxTime; tick++ {
		if raw := log.Sync(); raw == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// IsInitialized initialized or not
func IsInitialized() bool {
	return initialized.Load()
}

func createWriteSyncer(filename string, maxSize, maxAge, maxBackup int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxAge:     maxAge,
		MaxBackups: maxBackup,
	}
	return zapcore.AddSync(lumberJackLogger)
}

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

	switch logType {
	case outTypeJson:
		return zapcore.NewJSONEncoder(config)
	case outTypeConsole:
		return zapcore.NewConsoleEncoder(config)
	default:
		return zapcore.NewConsoleEncoder(config)
	}
}

func createTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000000"))
}
