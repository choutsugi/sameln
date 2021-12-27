// Package collector 基于tail库的日志收集模块
package collector

import (
	"LogAgent/logic/kafka"
	"LogAgent/logic/types"
	"LogAgent/universal/error"
	"LogAgent/universal/logger"
	"context"
	"github.com/hpcloud/tail"
	"strings"
)

// 日志收集任务结构
type task struct {
	path   string             // 日志文件路径
	topic  string             // kafka主题
	ins    *tail.Tail         // tail实例
	ctx    context.Context    // 用于控制收集任务结束
	cancel context.CancelFunc // 用于控制收集任务结束
}

// 初始化日志收集任务：构造collectTask。
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
		logger.L().Warnf("collector: init task %s failed.", t.path)
		return error.NewError(raw, error.CodeTailInitTaskFailed)
	}

	return error.Null()
}

// 运行日志收集任务
func (t *task) run() {
	logger.L().Infof("TailFile: task %s started.", t.path)
	for {
		select {
		case <-t.ctx.Done():
			logger.L().Warnf("TailFile: task %s stopped.", t.path)
			return
		case line, ready := <-t.ins.Lines:
			if !ready {
				logger.L().Warnf("TailFile: task %s failed to read log.", t.path)
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
			logger.L().Debugf("TailFile: task %s sent message successfully.", t.path)
		}
	}
}

// 创建日志收集任务
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
