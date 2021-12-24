package collector

import (
	"LogAgent/common/error"
	"LogAgent/common/logger"
	"LogAgent/logic/kafka"
	"LogAgent/logic/models"
	"context"
	"github.com/hpcloud/tail"
	"strings"
)

// collectTask 每个Etcd中的配置项对应一个Task
type collectTask struct {
	path     string
	topic    string
	instance *tail.Tail
	ctx      context.Context
	cancel   context.CancelFunc
}

func (t *collectTask) init() *error.Error {
	var raw error.RawErr
	config := tail.Config{
		Location: &tail.SeekInfo{
			Offset: 0,
			Whence: 2,
		},
		ReOpen:    true,
		MustExist: false,
		Poll:      true,
		Follow:    true,
	}

	t.instance, raw = tail.TailFile(t.path, config)
	if raw != nil {
		return error.NewError(raw, error.CodeTailInitTaskFailed)
	}

	return error.Null()
}

func (t *collectTask) run() {
	logger.L().Infof("TailFile: collectTask %s started.", t.path)
	for {
		select {
		case <-t.ctx.Done():
			logger.L().Warnf("TailFile: collectTask %s stopped.", t.path)
			return
		case line, ok := <-t.instance.Lines:
			if !ok {
				logger.L().Warnf("TailFile: collectTask %s failed to read log.", t.path)
				continue
			}
			if len(strings.Trim(line.Text, "\r")) == 0 {
				continue
			}
			msg := &kafka.ProducerMessage{
				Topic: t.topic,
				Value: kafka.StringEncoder(line.Text),
			}
			kafka.Write(msg)
			logger.L().Infof("TailFile: collectTask %s sent message successfully.", t.path)
		}
	}
}

func newTask(config models.CollectEntry) *collectTask {
	ctx, cancel := context.WithCancel(context.Background())
	task := &collectTask{
		path:     config.Path,
		topic:    config.Topic,
		instance: nil,
		ctx:      ctx,
		cancel:   cancel,
	}
	return task
}
