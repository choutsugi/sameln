package watch

import (
	"LogAgent/common/models"
	"LogAgent/logger"
)

func ConfigFileUpdate() {
	go func() {
		for {
			select {
			case msg := <-models.ConfigFileUpdateChan:
				if msg.IsUnmarshal {
					if logger.IsInitialized {
						logger.L.Infof("配置文件%s已更新，解析成功", msg.FileName)
					}
					continue
				}
				if logger.IsInitialized {
					logger.L.Warnf("配置文件%s已更新，解析失败", msg.FileName)
				}
			}
		}
	}()
}
