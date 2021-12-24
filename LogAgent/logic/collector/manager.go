package collector

import (
	"LogAgent/common/error"
	"LogAgent/common/logger"
	"LogAgent/logic/models"
)

type taskManager struct {
	taskMap          map[string]*collectTask
	CollectEntryList []models.CollectEntry
	queue            chan []models.CollectEntry
}

var (
	mgr *taskManager
)

func Init(allConfig []models.CollectEntry) {
	manager := &taskManager{
		taskMap:          make(map[string]*collectTask, 20),
		CollectEntryList: allConfig,
		queue:            make(chan []models.CollectEntry),
	}

	for _, config := range allConfig {
		task := newTask(config)
		if err := task.init(); err != error.Null() {
			logger.L().Warnf("TailFile: create task %s failed.", config.Path)
			continue
		}
		manager.taskMap[task.path] = task
		logger.L().Infow("TailFile: task %s is ready to start.", task.path)
		go task.run()
	}

	go manager.watch()
	return
}

func (mgr *taskManager) isExist(conf models.CollectEntry) (ok bool) {
	_, ok = mgr.taskMap[conf.Path]
	return
}

func (mgr *taskManager) watch() {
	for {
		allConf := <-mgr.queue
		logger.L().Infof("TailFile: configuration has been updated from etcd.")
		for _, conf := range allConf {
			if mgr.isExist(conf) {
				continue
			}
			task := newTask(conf)
			if err := task.init(); err != error.Null() {
				logger.L().Warnf("TailFile: create task %s failed.", conf.Path)
				continue
			}
			mgr.taskMap[task.path] = task
			logger.L().Infow("TailFile: task %s is ready to start.", task.path)
			go task.run()
		}

		for key, task := range mgr.taskMap {
			var isExist bool
			for _, conf := range allConf {
				if key == conf.Path {
					isExist = true
					break
				}
			}
			if !isExist {
				logger.L().Infof("TailFile: task:%s is ready to stop.", key)
				task.cancel()
				delete(mgr.taskMap, task.path)
			}
		}
	}
}

func UpdateConf(allConf []models.CollectEntry) {
	mgr.queue <- allConf
}
