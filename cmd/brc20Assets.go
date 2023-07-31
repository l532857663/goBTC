package main

import (
	"fmt"
	"goBTC"
	"goBTC/client"
	"goBTC/elastic"
	"goBTC/global"
	"goBTC/models"
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
	global.LOG.Info("GetUnSyncOrdToken", zap.Any("total", res.Total.Value))
	for _, hit := range res.Hits {
		ordToken := &elastic.OrdToken{}
		err := utils.Map2Struct(hit.Source, ordToken)
		if err != nil {
			global.LOG.Error("The ord_token utils.Map2Struct error", zap.Any("index", hit.Index), zap.Any("id", hit.Id))
			continue
		}
		brc20 := &models.OrdBRC20{}
		err = utils.Map2Struct(ordToken.Brc20Info, brc20)
		if err != nil {
			global.LOG.Error("The brc20 utils.Map2Struct error", zap.Any("index", hit.Index), zap.Any("id", hit.Id), zap.Any("info", ordToken.Brc20Info), zap.Error(err))
			continue
		}
		if brc20.P != "brc-20" {
			continue
		}
		err = ord.DealWithBrc20Info(ordToken, brc20)
		if err != nil {
			global.LOG.Error("ord.DealWithBrc20Info error", zap.Any("index", hit.Index), zap.Any("id", hit.Id), zap.Any("info", ordToken.Brc20Info), zap.Error(err))
			continue
		}
	}
}
