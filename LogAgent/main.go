package main

import (
	"LogAgent/app"
	"LogAgent/blunder"
	"LogAgent/logger"
	"LogAgent/settings"
	"fmt"
	"go.uber.org/zap"
	"time"
)

var (
	res = new(blunder.Error)
)

func main() {
	// 1.配置模块初始化
	if res = settings.Init(); res.Err != nil {
		fmt.Println(time.Now().Local(), res.Msg, res.Err)
		return
	}
	fmt.Println(time.Now().Local(), res.Msg)

	// 2.设置APP模式
	switch settings.Config.App.Mode {
	case app.ModeRelease:
		app.GlobalMode = app.ModeRelease
	case app.ModeDevelop:
		app.GlobalMode = app.ModeDevelop
	default:
		fmt.Println(blunder.GetMsg(blunder.CodeSysUnknownAppMode))
		app.GlobalMode = app.ModeDevelop
	}

	// 3.初始化日志模块
	if res = logger.Init(settings.Config.Log, app.GlobalMode); res.Err != nil {
		fmt.Println(time.Now().Local(), res.Msg, res.Err)
		return
	}
	defer logger.Sync()
	zap.L().Debug(res.Msg)

}
