package main

import (
	"LogAgent/app"
	"LogAgent/common/record"
	"LogAgent/error"
	"LogAgent/logger"
	"LogAgent/settings"
	"time"
)

var (
	ret = new(error.Error)
)

func main() {
	// 1.配置模块初始化
	if ret = settings.Init(); ret != error.Null() {
		record.Error(ret)
		return
	}
	record.Succeed("配置模块初始化成功")

	// 2.设置APP模式
	app.GlobalMode = settings.GetGlobalMode()
	record.Hint("当前运行模式%s", app.GlobalMode)

	// 3.初始化日志模块
	if ret = logger.Init(settings.Config.Log, app.GlobalMode); ret != error.Null() {
		record.Error(ret)
		return
	}
	defer logger.Sync()
	record.Succeed("日志模块初始化成功")

	logger.L.Infof("he")

	go app.Run()
	for {
		time.Sleep(time.Second)
	}
}
