package etcd

import (
	"LogAgent/common/error"
	"LogAgent/common/logger"
	"LogAgent/common/system"
	"LogAgent/logic/models"
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

var (
	client *clientv3.Client
)

func Init(address []string) *error.Error {
	var raw error.RawErr
	client, raw = clientv3.New(clientv3.Config{
		Endpoints:   address,
		DialTimeout: time.Second * 5,
	})
	if raw != nil {
		logger.L().Errorw("Etcd模块初始化失败", "err", raw.Error())
		return error.NewError(raw, error.CodeEtcdConnectFailed)
	}
	return error.Null()
}

// GetConf 获取配置项
func GetConf(key string) (collectEntryList []models.CollectEntry, err *error.Error) {
	// 获取IP生成Key
	ip, err := system.LocalIP()
	if err != nil {
		return
	}
	key = fmt.Sprintf(key, ip)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	resp, raw := client.Get(ctx, key)
	if err != nil {
		logger.L().Warnw(fmt.Sprintf("读取Key为%s的配置失败", key), "err", raw.Error())
		return
	}

	if len(resp.Kvs) == 0 {
		logger.L().Errorf("etcd: conf of key:%s is not exist", key)
		return
	}

	logger.L().Infow(fmt.Sprintf("读取Key为%s的配置成功", key), "err", raw.Error())

	ret := resp.Kvs[0]
	// 对从etcd获取的Json格式的配置数据进行解析
	raw = json.Unmarshal(ret.Value, &collectEntryList)
	if err != nil {
		logger.L().Warnw(fmt.Sprintf("解析Key为%s的配置失败", key), "err", raw.Error())
		return
	}
	logger.L().Infow(fmt.Sprintf("解析Key为%s的配置成功", key), "err", raw.Error())
	return
}

// WatchConf 监视etcd配置变化
func WatchConf(key string) {
	watchChan := client.Watch(context.Background(), key)
	var newConf []models.CollectEntry
	for resp := range watchChan {
		for _, event := range resp.Events {
			newConf = []models.CollectEntry{}
			logger.L().Info("etcd: configuration has been updated.")
			fmt.Printf("type:%s, key:%s, value:%s", event.Type, event.Kv.Key, event.Kv.Value)
			err := json.Unmarshal(event.Kv.Value, &newConf)
			if err != nil {
				logger.L().Errorf("etcd: conf of key:%s unmarshal failed, err:%v", event.Kv.Key, err)
				continue
			}
			// 如果配置更新则通知tailfile刷新任务
			//tailfile.UpdateConf(newConf)
		}
	}
}
