package app

import "LogAgent/common/watch"

var GlobalMode string

func Run() {
	go watch.ConfigFileUpdate()
}
