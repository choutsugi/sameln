package watch

import (
	"LogAgent/universal/generic"
	"LogAgent/universal/logger"
)

func ConfigFileUpdate() {
	go func() {
		for {
			select {
			case msg := <-generic.ConfigFileUpdateChan:
				if msg.IsUnmarshal {
					if logger.IsInitialized() {
						logger.L().Infof("The Watch module monitors that the config(%s) has been updated and parsed successfully.", msg.FileName)
					}
					continue
				}
				if logger.IsInitialized() {
					logger.L().Errorf("The Watch module monitors that the config(%s) has been updated but parsed successfully! Error:%s", msg.FileName, msg.Raw.Error())
				}
			default:

			}
		}
	}()
}
