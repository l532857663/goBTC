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

	"go.uber.org/zap"
)

var (
	// 全局参数
	srv *client.BTCClient
	log *zap.Logger
)

func main() {
	fmt.Println("vim-go")
	global.MysqlFlag = false
	goBTC.MustLoad("./config.yml")
	srv = global.Client
	log = global.LOG

	fmt.Printf("wch------ Start\n")
	// TestGetTx()
	TestGetBlock()
	fmt.Printf("wch------ END\n")
	// go server.CheckNewHeight(845492, server.GetTransferByBlockHeight)
	utils.SignalHandler("scanUTXO", goBTC.Shutdown)
}

func TestGetTx() {
	startTime := time.Now().Unix()
	blockInfo, err := srv.GetBlockInfoByHeight(845492)
	if err != nil {
		logutils.LogErrorf(global.LOG, "GetBlockInfoByHash error: %+v", err)
		return
	}
	endTime := time.Now().Unix()
	fmt.Printf("wch---- get block time: %+v\n", endTime-startTime)
	// fmt.Printf("wch---- blockInfo: %+v\n", blockInfo)
	fmt.Printf("wch---- tx len: %+v\n", len(blockInfo.Transactions))
	for i, txInfo := range blockInfo.Transactions {
		fmt.Printf("index: %+v\n", i)
		if i == 4 {
			break
		}
		// 添加计数器
		server.Wg.Add(1)
		go server.GetUTXOInfoByTransferInfo(txInfo)
	}
	server.Wg.Wait()
}

func TestGetBlock() {
	server.GetTransferByBlockHeight(2818745, 2818745)
}
