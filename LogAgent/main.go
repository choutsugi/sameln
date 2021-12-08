package main

import (
	"LogAgent/settings"
)

func main() {
	// 1.加载配置
	err := settings.Init()
	if err != nil {
		return
	}
}
