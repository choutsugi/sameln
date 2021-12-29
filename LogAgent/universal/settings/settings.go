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

func Init() *error.Error {
	if initialized.Load() {
		record.Warn("The Settings module unable to re-initialize!")
		return error.NewError(nil, codes.InitSettingsFailed)
	}
	// get command line parameters
	filePath := flag.String("config", "./config.yaml", "log_agent -config=\"./config.yaml\"")
	flag.Parse()
	viper.SetConfigFile(*filePath)

	// read config file
	record.Info("The Settings module starts to load config file(%s).", *filePath)
	if err := viper.ReadInConfig(); err != nil {
		record.Warn("The Settings module loads config file(%s) unsuccessfully!", *filePath)
		return error.NewError(err, codes.InitSettingsFailed)
	}

	// parse config file
	record.Info("The Settings module starts to parse config file(%s).", *filePath)
	if err := viper.Unmarshal(Config); err != nil {
		record.Warn("The Settings module parses config file(%s) unsuccessfully!", *filePath)
		return error.NewError(err, codes.InitSettingsFailed)
	}

	// watch config file
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		if raw := viper.Unmarshal(Config); raw != nil {
			msg := generic.FileUpdateMsg{
				FileName:    *filePath,
				IsUnmarshal: false,
				Raw:         raw,
			}
			generic.ConfigFileUpdateChan <- msg
		} else {
			msg := generic.FileUpdateMsg{
				FileName:    *filePath,
				IsUnmarshal: true,
				Raw:         nil,
			}
			generic.ConfigFileUpdateChan <- msg
		}
	})
	initialized.Store(true)
	return error.Null()
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
