package app

import (
	"LogAgent/common/error"
	"LogAgent/common/logger"
	"LogAgent/common/record"
	"LogAgent/common/settings"
	"LogAgent/common/system"
	"LogAgent/common/watch"
	"LogAgent/logic/etcd"
	"LogAgent/logic/kafka"
	"fmt"
	"time"
)

var (
	ret        = error.Null()
	GlobalMode string
)

func Run() {
	if initialize() != error.Null() {
		return
	}
	server()
}

func initialize() *error.Error {
	// 1.配置模块初始化
	if ret = settings.Init(); ret != error.Null() {
		record.Error(ret)
		return ret
	}
	record.Succeed("配置模块初始化成功")

	// 2.设置APP模式
	GlobalMode = settings.GetGlobalMode()
	record.Hint("当前运行模式%s", GlobalMode)

	// 3.初始化日志模块
	if ret = logger.Init(settings.Config.Log, GlobalMode); ret != error.Null() {
		record.Error(ret)
		return ret
	}
	defer logger.Sync()
	record.Succeed("日志模块初始化成功")

	// 4.初始化连接Kafka
	if ret = kafka.Init([]string{settings.Config.Kafaka.Addr}, settings.Config.Kafaka.Port, settings.Config.Kafaka.ChanSize); ret != error.Null() {
		return ret
	}
	logger.L().Info("Kafka模块初始化成功")

	// 5.初始化连接Etcd
	if ret = etcd.Init([]string{settings.Config.Etcd.Addr}); ret != error.Null() {
		return ret
	}
	logger.L().Info("Etcd模块初始化成功")
	return ret

}

func server() {
	defer logger.Sync()
	// 监视配置文件更新
	go watch.ConfigFileUpdate()

	for {
		time.Sleep(time.Second)
		fmt.Println(system.LocalTime())
	}
}
