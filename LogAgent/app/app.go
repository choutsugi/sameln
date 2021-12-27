package app

import (
	"LogAgent/logic/collector"
	"LogAgent/logic/etcd"
	"LogAgent/logic/kafka"
	"LogAgent/logic/types"
	"LogAgent/universal/error"
	"LogAgent/universal/logger"
	"LogAgent/universal/record"
	"LogAgent/universal/settings"
	"LogAgent/universal/system"
	"LogAgent/universal/watch"
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
		record.Fatal(ret)
		return ret
	}
	record.Info("配置模块初始化成功")

	// 2.设置APP模式
	GlobalMode = settings.GetGlobalMode()
	record.Info("当前运行模式%s", GlobalMode)

	// 3.初始化日志模块
	if ret = logger.Init(settings.Config.Log, GlobalMode); ret != error.Null() {
		record.Fatal(ret)
		return ret
	}
	defer logger.Sync()
	logger.L().Info("日志模块初始化成功")

	// 4.初始化连接Kafka
	if ret = kafka.Init(settings.Config.Kafaka); ret != error.Null() {
		return ret
	}
	logger.L().Info("Kafka模块初始化成功")

	// 5.初始化连接Etcd
	if ret = etcd.Init(settings.Config.Etcd); ret != error.Null() {
		return ret
	}
	logger.L().Info("Etcd模块初始化成功")
	return ret

}

func server() {
	var err *error.Error
	var entries []types.CollectEntry

	// 监视系统配置文件更新
	go watch.ConfigFileUpdate()

	if err = etcd.PutConfig(settings.Config.Etcd.CollectKey); err != error.Null() {
		return
	}
	if entries, err = etcd.GetConfig(settings.Config.Etcd.CollectKey); err != error.Null() {
		return
	}

	defer logger.Sync()
	defer kafka.Close()
	defer etcd.Close()

	// 监视日志收集Key更新
	go etcd.WatchConf(settings.Config.Etcd.CollectKey)
	// 初始化收集器并启动
	collector.Start(entries)

	for {
		time.Sleep(time.Second)
		fmt.Println("current time:", system.LocalTime())
	}
}
