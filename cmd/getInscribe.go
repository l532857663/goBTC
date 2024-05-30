package main

import (
	"fmt"
	"goBTC"
	"goBTC/client"
	"goBTC/elastic"
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
	GetInscribeInfo()
	if global.MysqlFlag {
		utils.SignalHandler("getInscribe", goBTC.Shutdown)
	}
}

func GetInscribeInfo() {
	res, err := ord.GetUnSyncOrdInscribe()
	if err != nil {
		global.LOG.Error("ord.GetInscribeInfo", zap.Error(err))
		return
	}
	global.LOG.Info("GetInscribeInfo", zap.Any("total", res.Total.Value))
	for _, hit := range res.Hits {
		inscribeInfo := &elastic.InscribeInfo{}
		err := utils.Map2Struct(hit.Source, inscribeInfo)
		if err != nil {
			global.LOG.Error("The inscribe utils.Map2Struct error", zap.Any("index", hit.Index), zap.Any("id", hit.Id))
			continue
		}
		fmt.Printf("wch----- data: %+v\n", inscribeInfo)
	}
}
