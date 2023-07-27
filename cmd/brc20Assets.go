package main

import (
	"fmt"
	"goBTC"
	"goBTC/client"
	"goBTC/global"
	"goBTC/ord"
	"goBTC/utils"

	"go.uber.org/zap"
)

var (
	srv *client.BTCClient
	log *zap.Logger
)

func main() {
	fmt.Println("vim-go")
	// global.MysqlFlag = true
	goBTC.MustLoad("./config.yml")
	srv = global.Client
	log = global.LOG
	CheckBrc20Assets()
	if global.MysqlFlag {
		utils.SignalHandler("brc20Assets", goBTC.Shutdown)
	}
}

func CheckBrc20Assets() {
	res, err := ord.GetUnSyncOrdToken()
	if err != nil {
		global.LOG.Error("ord.GetUnSyncOrdToken", zap.Error(err))
		return
	}
	for _, hit := range res {
		fmt.Printf("wch----- hit: %+v\n", hit)
	}
}
