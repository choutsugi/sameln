// Package collector 基于tail库的日志收集模块
package collector

import (
	"LogAgent/logic/kafka"
	"LogAgent/logic/types"
	"LogAgent/universal/codes"
	"LogAgent/universal/error"
	"LogAgent/universal/logger"
	"context"
	"github.com/hpcloud/tail"
	"strings"
)

type task struct {
	path   string     // log file's path
	topic  string     // kafka topic
	ins    *tail.Tail // tail instance
	ctx    context.Context
	cancel context.CancelFunc
}

func (t *task) init() *error.Error {
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

	t.ins, raw = tail.TailFile(t.path, config)
	if raw != nil {
		logger.L().Warnf("The Collector module initializes tail-task(%s) unsuccessfully!", t.topic)
		return error.NewError(raw, codes.CollectorInitTaskFailed)
	}

	return error.Null()
}

func (t *task) run() {
	logger.L().Infof("The Collector module starts to run tail-task(%s).", t.topic)
	for {
		select {
		case <-t.ctx.Done():
			logger.L().Infof("The Collector module stops to run tail-task(%s).", t.topic)
			return
		case line, ready := <-t.ins.Lines:
			if !ready {
				logger.L().Warnf("The Collector module's tail-task(%s) reads log-file(%s) unsuccessfully!", t.topic, t.path)
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
			logger.L().Debugf("The Collector module's tail-task(%s) sends message to Kafka module successfully!", t.topic)
		}
	}
}

func createTask(config types.CollectEntry) *task {
	ctx, cancel := context.WithCancel(context.Background())
	task := &task{
		path:   config.Path,
		topic:  config.Topic,
		ins:    nil,
		ctx:    ctx,
		cancel: cancel,
	}
	return task
}
