package gracefulshutdown

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var cleanupOperations map[string]func() = nil
var cleanupOperationsMutuex sync.Mutex

func initCleanupOperations() {
	cleanupOperationsMutuex.Lock()
	if cleanupOperations == nil {
		cleanupOperations = make(map[string]func())
	}
	cleanupOperationsMutuex.Unlock()
}

func AddCleanupOperation(name string, op func()) {
	initCleanupOperations()
	cleanupOperationsMutuex.Lock()
	cleanupOperations[name] = op
	cleanupOperationsMutuex.Unlock()
}

func Wait(timeout time.Duration) {
	var wg sync.WaitGroup
	initCleanupOperations()
	signalReciever := make(chan os.Signal, 1)
	signal.Notify(signalReciever, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-signalReciever
	for _, op := range cleanupOperations {
		wg.Add(1)
		innerOp := op
		go func() {
			defer wg.Done()
			innerOp()
		}()
	}
	time.AfterFunc(timeout, func() {
		os.Exit(0)
	})
	wg.Wait()
}
