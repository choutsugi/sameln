package settings

import (
	"LogAgent/universal/codes"
	"LogAgent/universal/error"
	"LogAgent/universal/generic"
	"LogAgent/universal/record"
	"flag"
	"go.uber.org/atomic"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	Config      = new(ConfigType)
	initialized atomic.Bool
)

// Init 配置模块初始化：从文件读取配置信息并赋值于settings.Config
func Init() *error.Error {
	if initialized.Load() {
		return error.Null()
	}
	// 通过命令行参数指定配置文件路径
	filePath := flag.String("config", "./config.yaml", "log_agent -config=\"./config.yaml\"")
	flag.Parse()
	viper.SetConfigFile(*filePath)

	// 读取配置信息
	record.Info("The Settings module starts to load config file(%s).", *filePath)
	if err := viper.ReadInConfig(); err != nil {
		record.Warn("The Settings module loads config file(%s) unsuccessfully!", *filePath)
		return error.NewError(err, codes.InitSettingsFailed)
	}

	// 解析配置文件
	record.Info("The Settings module starts to parse config file(%s).", *filePath)
	if err := viper.Unmarshal(Config); err != nil {
		record.Warn("The Settings module parses config file(%s) unsuccessfully!", *filePath)
		return error.NewError(err, codes.InitSettingsFailed)
	}

	// 监控配置文件
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		if err := viper.Unmarshal(Config); err != nil {
			msg := generic.FileUpdateMsg{
				FileName:    *filePath,
				IsUnmarshal: false,
			}
			generic.ConfigFileUpdateChan <- msg
		} else {
			msg := generic.FileUpdateMsg{
				FileName:    *filePath,
				IsUnmarshal: true,
			}
			generic.ConfigFileUpdateChan <- msg
		}
	})
	initialized.Store(true)
	return error.NullWithCode(codes.InitSettingsSucceed)
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
		record.Warn("The Settings module parses running mode(default:%s) unsuccessfully!", mode)
	}
	return
}
