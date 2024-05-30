package server

import (
	"goBTC/global"
	"goBTC/utils/logutils"
	"sync"
	"time"
)

type dealFunc func(int64, int64)

var (
	Wg  sync.WaitGroup
	Wg1 sync.WaitGroup
)

func CheckNewHeight(startHeight int64, f dealFunc) {
	srv := global.Client
	log := global.LOG
	logutils.LogInfof(log, "[CheckNewHeight] Start")
	for {
		newHigh, err := srv.GetBlockCount()
		if err != nil {
			logutils.LogErrorf(log, "GetBlockCount error: %+v", err)
			return
		}
		if startHeight > newHigh {
			time.Sleep(5 * time.Minute)
			continue
		}
		f(startHeight, newHigh)
		startHeight = newHigh + 1
		time.Sleep(5 * time.Minute)
		logutils.LogInfof(log, "[CheckNewHeight] Once time New high: %v", newHigh)
	}
}
