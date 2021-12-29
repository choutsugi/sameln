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
	"LogAgent/universal/watch"
	"time"
)

var (
	GlobalMode string
)

func Run() {
	if initialize() != error.Null() {
		record.Warn("The App fails to initialize and will exit!")
		return
	}
	server()
}

func initialize() (err *error.Error) {
	if err = settings.Init(); err != error.Null() {
		record.Error(err)
		return
	}
	record.Info("Initialize the Settings module successfully.")

	GlobalMode = settings.GetGlobalMode()
	record.Info("Current running mode: %s", GlobalMode)

	if err = logger.Init(settings.Config.Log, GlobalMode); err != error.Null() {
		record.Error(err)
		return
	}
	defer logger.Sync()
	logger.L().Info("Initialize the Logger module successfully.")

	if err = kafka.Init(settings.Config.Kafaka); err != error.Null() {
		record.Error(err)
		return
	}
	logger.L().Info("Initialize the Kafka module successfully.")

	if err = etcd.Init(settings.Config.Etcd); err != error.Null() {
		record.Error(err)
		return
	}
	logger.L().Info("Initialize the Etcd module successfully.")
	return
}

func server() {
	var err *error.Error
	var entries []types.CollectEntry

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

	go etcd.WatchConf(settings.Config.Etcd.CollectKey)
	collector.Start(entries)

	for {
		time.Sleep(time.Second)
	}
}
