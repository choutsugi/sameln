package etcd

import (
	"LogAgent/logic/collector"
	"LogAgent/logic/types"
	"LogAgent/universal/codes"
	"LogAgent/universal/error"
	"LogAgent/universal/generic"
	"LogAgent/universal/logger"
	"LogAgent/universal/settings"
	"LogAgent/universal/system"
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"go.uber.org/atomic"
	"time"
)

var (
	cli         *clientv3.Client
	initialized atomic.Bool
)

func Init(etcdConfig *settings.EtcdConfigType) *error.Error {
	if initialized.Load() {
		logger.L().Error("The Etcd module unable to re-initialize!")
		return error.NewError(nil, codes.InitEtcdFailed)
	}
	var raw error.RawErr
	cli, raw = clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdConfig.Addr},
		DialTimeout: time.Second * 5,
	})
	if raw != nil {
		logger.L().Errorf("The Etcd module connects to Etcd service unsuccessfully! Error:%s", raw.Error())
		return error.NewError(raw, codes.EtcdConnectFailed)
	}
	initialized.Store(true)
	return error.Null()
}

func PutConfig(key string) *error.Error {
	ip, err := system.GetLocalIP()
	if err != error.Null() {
		return err
	}
	key = fmt.Sprintf(key, ip)
	str := `[{"path":"D:/ProgramData/LogAgent/logs/s4.log","topic":"s4_log"},{"path":"D:/ProgramData/LogAgent/logs/web.log","topic":"web_log"},{"path":"D:/ProgramData/LogAgent/logs/s5.log","topic":"s5_log"}]`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_, raw := cli.Put(ctx, key, str)
	cancel()
	if raw != nil {
		logger.L().Errorf("The Etcd module sets config(%s) unsuccessfully! Error:%s", key, raw.Error())
		return error.NewError(raw, codes.EtcdSetConfigFailed)
	}
	logger.L().Infof("The Etcd module sets config(%s) successfully!", key)

	return error.Null()
}

func GetConfig(key string) ([]types.CollectEntry, *error.Error) {
	var entries []types.CollectEntry
	var err *error.Error

	ip, err := system.GetLocalIP()
	if err != error.Null() {
		logger.L().Errorf("The Etcd module gets local ip unsuccessfully! Error:%s", err.RawErr().Error())
		return nil, error.NewError(err.RawErr(), codes.EtcdGetIpFailed)
	}
	key = fmt.Sprintf(key, ip)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, raw := cli.Get(ctx, key)
	cancel()
	if raw != nil {
		logger.L().Errorf("The Etcd module gets config unsuccessfully! Error:%s", raw.Error())
		return nil, error.NewError(raw, codes.EtcdGetConfigFailed)
	}

	if len(resp.Kvs) == 0 {
		logger.L().Errorf("The Etcd module cannot find the specified-config(%s) from the Etcd service!", key)
		return nil, error.NewError(nil, codes.EtcdConfigNotFound)
	}

	logger.L().Infof("The Etcd module gets config-file(%s) successfully!", key)

	ret := resp.Kvs[0]
	raw = json.Unmarshal(ret.Value, &entries)
	if raw != nil {
		logger.L().Errorf("The Etcd module parses config(%s) unsuccessfully! Error:%s", key, raw.Error())
		return nil, error.NewError(nil, codes.EtcdConfigParseFailed)
	}
	logger.L().Infof("The Etcd module parses config(%s) successfully.", key)
	return entries, error.Null()
}

func WatchConf(key string) {
	watchChan := cli.Watch(context.Background(), key)
	var entries []types.CollectEntry
	for resp := range watchChan {
		for _, event := range resp.Events {
			entries = []types.CollectEntry{}
			logger.L().Infof("The Etcd module monitors that the config(%s) has been updated! Type:%s Key:%s Value:%s", key, event.Type, event.Kv.Key, event.Kv.Value)
			raw := json.Unmarshal(event.Kv.Value, &entries)
			if raw != nil {
				logger.L().Errorf("The Etcd module parses config(%s) unsuccessfully! Error:%s", key, raw.Error())
				continue
			}
			logger.L().Infof("The Etcd module parses config(%s) successfully!", key)
			collector.UpdateConfig(entries)
		}
	}
}

func Close() {
	for tick := 0; tick < generic.TryCloseWithMaxTime; tick++ {
		if raw := cli.Close(); raw == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
}
