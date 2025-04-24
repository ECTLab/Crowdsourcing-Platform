package App

import (
	"navigation/pkg/Server"
	gracefulshutdown "navigation/tools/graceful_shutdown"
	"time"
)

func Run() {
	Init()
	go Server.ServeHttp()
	gracefulshutdown.Wait(2 * time.Second)
}
