package collector

import (
	"LogAgent/logic/types"
	"LogAgent/universal/error"
	"LogAgent/universal/logger"
	"go.uber.org/atomic"
)

type taskManager struct {
	tasks   map[string]*task
	Entries []types.CollectEntry
	queue   chan []types.CollectEntry
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
		tasks:   make(map[string]*task, 20),
		Entries: entries,
		queue:   make(chan []types.CollectEntry),
	}

	for _, entry := range entries {
		task := createTask(entry)
		if err := task.init(); err != error.Null() {
			logger.L().Errorf("The Collector module initializes tail-task(%s) unsuccessfully!", task.topic)
			continue
		}
		manager.tasks[task.path] = task
		logger.L().Infof("The Collector module initializes tail-task(%s) successfully and ready to start.", task.topic)
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
		logger.L().Info("The Collector module has been notified by the Etcd module that the config has been updated.")
		for _, conf := range entries {
			if mgr.isExist(conf) {
				continue
			}
			task := createTask(conf)
			if err := task.init(); err != error.Null() {
				logger.L().Errorf("The Collector module initializes tail-t(%s) unsuccessfully!", task.topic)
				continue
			}
			mgr.tasks[task.path] = task
			logger.L().Infof("The Collector module initializes tail-t(%s) successfully and ready to start.", task.topic)
			go task.run()
		}

		for key, t := range mgr.tasks {
			var isExist bool
			for _, entry := range entries {
				if key == entry.Path {
					isExist = true
					break
				}
			}
			if !isExist {
				logger.L().Infof("The Collector module's tail-t(%s) is ready to stop.", t.topic)
				t.cancel()
				t.ins.Cleanup()
				go func(t *task) {
					for {
						if t.ins.Stop() == nil {
							break
						}
					}
				}(t)
				delete(mgr.tasks, t.path)
			}
		}
	}
}

func UpdateConfig(entries []types.CollectEntry) {
	manager.queue <- entries
}
