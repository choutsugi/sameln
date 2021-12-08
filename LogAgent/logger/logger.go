package logger

import (
	"LogAgent/settings"
	"os"

	"go.uber.org/zap"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap/zapcore"
)

func Init(config *settings.ConfigType) (err error) {

	writer := newWriter(
		config.Log.FileName,
		config.Log.MaxSize,
		config.Log.MaxAge,
		config.Log.MaxBackups,
	)

	encoder := newEncoder(config.Log.Type)
	level := new(zapcore.Level)
	if err = level.UnmarshalText([]byte(config.Log.Level)); err != nil {
		return err
	}

	var core zapcore.Core
	if config.App.Mode == appModeDev {
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
	return nil
}

func Sync() {
	err := zap.L().Sync()
	if err != nil {
		//TODO
		return
	}
}

func newEncoder(logType string) zapcore.Encoder {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.TimeKey = "time"
	config.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncodeDuration = zapcore.SecondsDurationEncoder
	config.EncodeCaller = zapcore.ShortCallerEncoder

	switch logType {
	case logOutTypeJson:
		return zapcore.NewJSONEncoder(config)
	case logOutTypeConsole:
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
