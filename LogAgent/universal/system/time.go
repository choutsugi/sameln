package system

import "time"

func GetLocalTime() string {
	return time.Now().Local().Format("2006-01-02 15:04:05.000")
}

func GetUtcTime() string {
	return time.Now().Local().Format("2006-01-02T15:04:05.000+0800")
}
