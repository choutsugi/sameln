package main

import (
	"LogAgent/blunder"
	"LogAgent/logger"
	"LogAgent/settings"
	"fmt"
	"go.uber.org/zap"
)

func main() {
	// 1.加载配置
	if res := settings.Init(); res != blunder.NewWithSuccess() {
		fmt.Println(res.Msg, res.Err)
		return
	}

	// 2.初始化日志模块
	if err := logger.Init(settings.Config); err != nil {
		fmt.Println()
		return
	}
	defer logger.Sync()

	zap.L().Debug("success!")

}
