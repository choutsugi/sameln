package settings

import (
	"LogAgent/blunder"
	"flag"
	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"

	"github.com/spf13/viper"
)

var (
	Config = new(ConfigType)
)

// Init /* 配置模块初始化：从文件读取配置信息并赋值于settings.Config
func Init() *blunder.Error {
	// 通过命令行参数指定配置文件路径
	filePath := flag.String("config", "./config.yaml", "log_agent -config=\"./config.yaml\"")
	flag.Parse()
	viper.SetConfigFile(*filePath)

	// 读取配置信息
	if err := viper.ReadInConfig(); err != nil {
		return blunder.NewError(blunder.CodeSysSettingsInitFailed, err)
	}

	// 解析配置文件
	if err := viper.Unmarshal(Config); err != nil {
		return blunder.NewError(blunder.CodeSysSettingsInitFailed, err)
	}

	// 监控配置文件
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		if err := viper.Unmarshal(Config); err != nil {
			zap.L().Debug(blunder.GetMsg(blunder.CodeSysSettingsConfigUpdated))
		}
	})

	return blunder.NewSuccess(blunder.CodeSysSettingsInitSucceed)
}
