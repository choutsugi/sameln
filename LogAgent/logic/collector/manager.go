package collector

import (
	"LogAgent/logic/types"
	"LogAgent/universal/error"
	"LogAgent/universal/logger"
	"go.uber.org/atomic"
)

type taskManager struct {
	tasks            map[string]*task
	CollectEntryList []types.CollectEntry
	queue            chan []types.CollectEntry
}

var (
	manager *taskManager
	started atomic.Bool
)

func Start(entries []types.CollectEntry) {
	if started.Load() {
		return
	}
	manager = &taskManager{
		tasks:            make(map[string]*task, 20),
		CollectEntryList: entries,
		queue:            make(chan []types.CollectEntry),
	}

	for _, entry := range entries {
		task := createTask(entry)
		if err := task.init(); err != error.Null() {
			logger.L().Warnf("TailFile: create task %s failed.", entry.Topic)
			continue
		}
		manager.tasks[task.path] = task
		logger.L().Infof("TailFile: task %s is ready to start.", task.topic)
		go task.run()
	}

	go manager.watch()
	started.Store(true)
	return
}

func (mgr *taskManager) isExist(conf types.CollectEntry) (ok bool) {
	_, ok = mgr.tasks[conf.Path]
	return
}

func (mgr *taskManager) watch() {
	for {
		entries := <-mgr.queue
		logger.L().Info("TailFile: configuration has been updated from etcd.")
		for _, conf := range entries {
			if mgr.isExist(conf) {
				continue
			}
			task := createTask(conf)
			if err := task.init(); err != error.Null() {
				logger.L().Warnf("TailFile: create task %s failed.", conf.Topic)
				continue
			}
			mgr.tasks[task.path] = task
			logger.L().Infof("TailFile: task %s is ready to start.", task.topic)
			go task.run()
		}

		for key, task := range mgr.tasks {
			var isExist bool
			for _, entry := range entries {
				if key == entry.Path {
					isExist = true
					break
				}
			}
			if !isExist {
				logger.L().Infof("TailFile: task %s is ready to stop.", task.topic)
				task.cancel()
				delete(mgr.tasks, task.path)
			}
		}
	}
}

func UpdateConfig(entries []types.CollectEntry) {
	manager.queue <- entries
}
