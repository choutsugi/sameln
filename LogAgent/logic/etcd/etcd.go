package etcd

import (
	"LogAgent/logic/collector"
	"LogAgent/logic/types"
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
	client      *clientv3.Client
	initialized atomic.Bool
)

func Init(etcdConfig *settings.EtcdConfigType) *error.Error {
	if initialized.Load() {
		return error.Null()
	}
	var raw error.RawErr
	client, raw = clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdConfig.Addr},
		DialTimeout: time.Second * 5,
	})
	if raw != nil {
		logger.L().Errorw("Etcd模块初始化失败", "err", raw.Error())
		return error.NewError(raw, error.CodeEtcdConnectFailed)
	}
	initialized.Store(true)
	return error.Null()
}

// PutConfig 设置配置项
func PutConfig(key string) *error.Error {
	// 获取IP生成Key
	ip, err := system.LocalIP()
	if err != error.Null() {
		return err
	}
	key = fmt.Sprintf(key, ip)
	str := `[{"path":"D:/ProgramData/LogAgent/logs/s4.log","topic":"s4_log"},{"path":"D:/ProgramData/LogAgent/logs/web.log","topic":"web_log"},{"path":"D:/ProgramData/LogAgent/logs/s5.log","topic":"s5_log"}]`

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	_, raw := client.Put(ctx, key, str)
	cancel()

	if raw != nil {
		logger.L().Warnw(fmt.Sprintf("设置Key为%s的配置失败", key), "err", raw.Error())
		return error.NewError(raw, error.CodeEtcdPutConfFailed)
	}

	return error.Null()
}

// GetConfig 获取配置项
func GetConfig(key string) (collectEntryList []types.CollectEntry, err *error.Error) {
	// 获取IP生成Key
	ip, err := system.LocalIP()
	if err != error.Null() {
		return
	}
	key = fmt.Sprintf(key, ip)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, raw := client.Get(ctx, key)
	cancel()

	if raw != nil {
		logger.L().Warnw(fmt.Sprintf("读取Key为%s的配置失败", key), "err", raw.Error())
		return
	}

	if len(resp.Kvs) == 0 {
		logger.L().Errorf("etcd: conf of key:%s is not exist", key)
		return
	}

	logger.L().Infow(fmt.Sprintf("读取Key为%s的配置成功", key))

	ret := resp.Kvs[0]
	// 对从etcd获取的Json格式的配置数据进行解析
	raw = json.Unmarshal(ret.Value, &collectEntryList)
	if raw != nil {
		logger.L().Warnw(fmt.Sprintf("解析Key为%s的配置失败", key), "err", raw.Error())
		return
	}
	logger.L().Infow(fmt.Sprintf("解析Key为%s的配置成功", key))
	return
}

// WatchConf 监视etcd配置变化
func WatchConf(key string) {
	watchChan := client.Watch(context.Background(), key)
	var newConf []types.CollectEntry
	for resp := range watchChan {
		for _, event := range resp.Events {
			newConf = []types.CollectEntry{}
			logger.L().Info("etcd: configuration has been updated.")
			fmt.Printf("type:%s, key:%s, value:%s", event.Type, event.Kv.Key, event.Kv.Value)
			err := json.Unmarshal(event.Kv.Value, &newConf)
			if err != nil {
				logger.L().Errorf("etcd: conf of key:%s unmarshal failed, err:%v", event.Kv.Key, err)
				continue
			}
			// 如果配置更新则通知tailfile刷新任务
			collector.UpdateConfig(newConf)
		}
	}
}

func Close() {
	for tick := 0; tick < generic.TryCloseWithMaxTime; tick++ {
		if raw := client.Close(); raw == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
}
