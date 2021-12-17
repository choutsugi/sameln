package settings

import (
	"LogAgent/error"
	"flag"
	"fmt"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"

	"github.com/spf13/viper"
)

var (
	Config = new(ConfigType)
)

// Init /* 配置模块初始化：从文件读取配置信息并赋值于settings.Config
func Init() *error.Error {
	// 通过命令行参数指定配置文件路径
	filePath := flag.String("config", "./config.yaml", "log_agent -config=\"./config.yaml\"")
	flag.Parse()
	viper.SetConfigFile(*filePath)

	// 读取配置信息
	if err := viper.ReadInConfig(); err != nil {
		return error.NewError(err, error.CodeSysSettingsInitFailed)
	}

	// 解析配置文件
	if err := viper.Unmarshal(Config); err != nil {
		return error.NewError(err, error.CodeSysSettingsInitFailed)
	}

	// 监控配置文件
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		if err := viper.Unmarshal(Config); err != nil {
			zap.L().Debug(error.GetInfo(error.CodeSysSettingsConfigUpdated))
		}
	})

	return error.NullWithCode(error.CodeSysSettingsInitSucceed)
}

func GetGlobalMode() string {
	var mode string
	switch Config.App.Mode {
	case ModeRelease:
		mode = ModeRelease
	case ModeDevelop:
		mode = ModeDevelop
	default:
		fmt.Println(error.GetInfo(error.CodeSysUnknownAppMode))
		mode = ModeDevelop
	}

	return mode
}
