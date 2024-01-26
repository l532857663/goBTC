package main

import (
	"fmt"
	"goBTC"
	"goBTC/client"
	"goBTC/global"
	"goBTC/server"
	"goBTC/utils"
	"goBTC/utils/logutils"
	"time"
)

var (
	srv *client.BTCClient
)

func main() {
	fmt.Println("vim-go")
	global.MysqlFlag = true
	goBTC.MustLoad("./config.yml")
	srv = global.Client
	go CheckOrder()
	if global.MysqlFlag {
		utils.SignalHandler("orderCheck", goBTC.Shutdown)
	}
}

func CheckOrder() {
	for {
		logutils.LogInfof(global.LOG, "Start check order")
		server.QueryPendingOrder4DB()
		logutils.LogInfof(global.LOG, "End check order")
		time.Sleep(5 * time.Second)
	}
}
