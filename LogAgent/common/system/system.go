package system

import "time"

// LocalTime 获取本地时间（毫秒级， 字符串）
func LocalTime() string {
	return time.Now().Local().Format("2006-01-02 15:04:05.000")
}
