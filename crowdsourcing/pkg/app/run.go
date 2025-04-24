package app

import (
	"crowdsourcing/pkg/server"
	gracefulshutdown "crowdsourcing/tools/graceful_shutdown"
	"time"
)

func Run() {
	Init()
	go server.ServeHttp()
	gracefulshutdown.Wait(2 * time.Second)
}
