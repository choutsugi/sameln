package logger

import (
	"LogAgent/settings"

	"go.uber.org/zap"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap/zapcore"
)

func Init(config *settings.ConfigType) {

}

func getEncoder() zapcore.Encoder {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.TimeKey = "time"
	config.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncodeDuration = zapcore.SecondsDurationEncoder
	config.EncodeCaller = zapcore.ShortCallerEncoder

	switch settings.Config.Log.Type {
	case logOutTypeJson:
		return zapcore.NewJSONEncoder(config)
	case logOutTypeConsole:
		return zapcore.NewConsoleEncoder(config)
	default:
		return zapcore.NewConsoleEncoder(config)
	}
}

func getLogWriter(filename string, maxSize, maxAge, maxBackup int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxAge:     maxAge,
		MaxBackups: maxBackup,
	}
	return zapcore.AddSync(lumberJackLogger)
}
