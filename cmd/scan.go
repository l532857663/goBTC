package main

import (
	"fmt"
	"goBTC"
	"goBTC/client"
	"goBTC/global"
	"goBTC/server"
	"goBTC/utils"
	"goBTC/utils/logutils"

	"go.uber.org/zap"
)

var (
	srv *client.BTCClient
	log *zap.Logger
)

func main() {
	fmt.Println("vim-go")
	global.MysqlFlag = false
	goBTC.MustLoad("./config.yml")
	srv = global.Client
	log = global.LOG
	// go server.CheckNewHeight(2543620)
	// TestGetTx()
	TestGetByTxhash()
	utils.SignalHandler("scan", goBTC.Shutdown)
}

func TestGetTx() {
	fmt.Printf("wch------ test\n")
	blockList := []int64{788344, 789649, 789792, 789793}
	txList := []int{2826, 2398, 1653, 309}
	for index, i := range blockList {
		blockInfo, err := global.Client.GetBlockInfoByHeight(i)
		if err != nil {
			logutils.LogErrorf(global.LOG, "GetBlockInfoByHash error: %+v", err)
			return
		}
		fmt.Printf("wch------ test1\n")
		sum, sumBrc20 := new(int), new(int)
		j := txList[index]
		server.Wg.Add(1)
		server.GetOneTxInfo(blockInfo, sum, sumBrc20, i, j, j)
		fmt.Printf("wch------ END\n")
	}
}

func TestGetByTxhash() {
	txHash := "136696530eca2daa4639c64025cd4d09611fbc30700af5576b9a0462cf99aa13"
	log.Info("[GetHashInfo] Start", zap.Any("txHash", txHash))
	txInfo, err := srv.GetRawTransactionByHash(txHash)
	if err != nil {
		log.Info("GetRawTransactionByHash", zap.Error(err))
		return
	}
	// fmt.Printf("wch---- txInfo: %+v\n", txInfo)
	witnessStr := client.GetTxWitnessByTxInfo(txInfo)
	res := client.GetScriptStringPlus(witnessStr)
	if res == nil {
		log.Info("GetScriptString not have inscription")
		return
	}
	fmt.Printf("wch----- witness info: %+v\n", witnessStr)
	return
}
