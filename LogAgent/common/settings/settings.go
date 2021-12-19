package settings

import (
	"LogAgent/common/bundles"
	"LogAgent/common/error"
	"LogAgent/common/record"
	"flag"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	Config = new(ConfigType)
)

// Init 配置模块初始化：从文件读取配置信息并赋值于settings.Config
func Init() *error.Error {
	// 通过命令行参数指定配置文件路径
	filePath := flag.String("config", "./config.yaml", "log_agent -config=\"./config.yaml\"")
	flag.Parse()
	viper.SetConfigFile(*filePath)

	// 读取配置信息
	record.Hint("开始加载配置文件%s", *filePath)
	if err := viper.ReadInConfig(); err != nil {
		record.Failed("加载配置文件%s失败", *filePath)
		return error.NewError(err, error.CodeSysSettingsInitFailed)
	}

	// 解析配置文件
	record.Hint("开始解析配置文件%s", *filePath)
	if err := viper.Unmarshal(Config); err != nil {
		record.Failed("解析配置文件%s失败", *filePath)
		return error.NewError(err, error.CodeSysSettingsInitFailed)
	}

	// 监控配置文件
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		if err := viper.Unmarshal(Config); err != nil {
			msg := bundles.FileUpdateMsg{
				FileName:    *filePath,
				IsUnmarshal: false,
			}
			bundles.ConfigFileUpdateChan <- msg
		} else {
			msg := bundles.FileUpdateMsg{
				FileName:    *filePath,
				IsUnmarshal: true,
			}
			bundles.ConfigFileUpdateChan <- msg
		}
	})

	return error.NullWithCode(error.CodeSysSettingsInitSucceed)
}

// GetGlobalMode 获取全局运行模式
func GetGlobalMode() (mode string) {
	switch Config.App.Mode {
	case ModeRelease:
		mode = ModeRelease
	case ModeDevelop:
		mode = ModeDevelop
	default:
		mode = ModeDevelop
		record.Failed("解析运行模式错误，使用默认值%s", mode)
	}
	return
}
