package main

import (
	"LogAgent/app"
	"LogAgent/error"
	"LogAgent/logger"
	"LogAgent/settings"
	"fmt"
	"time"
)

var (
	ret = new(error.Error)
)

func main() {
	// 1.配置模块初始化
	if ret = settings.Init(); ret != error.Null() {
		fmt.Println(time.Now().Local(), ret.Info(), ret.RawErr())
		return
	}
	fmt.Println(time.Now().Local(), ret.Info())

	// 2.设置APP模式
	app.GlobalMode = settings.GetGlobalMode()

	// 3.初始化日志模块
	if ret = logger.Init(settings.Config.Log, app.GlobalMode); ret != error.Null() {
		fmt.Println(time.Now().Local(), ret.Info(), ret.RawErr())
		return
	}
	defer logger.Sync()
	logger.Log.Infof("he")

}
