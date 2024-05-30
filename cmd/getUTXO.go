package main

import (
	"fmt"
	"goBTC"
	"goBTC/client"
	"goBTC/global"
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
	TestGetTx()
	// utils.SignalHandler("scanUTXO", goBTC.Shutdown)
}

func TestGetTx() {
	fmt.Printf("wch------ test\n")
	blockInfo, err := srv.GetBlockInfoByHeight(845492)
	if err != nil {
		logutils.LogErrorf(global.LOG, "GetBlockInfoByHash error: %+v", err)
		return
	}
	fmt.Printf("wch---- blockInfo: %+v\n", blockInfo)
	fmt.Printf("wch---- tx len: %+v\n", len(blockInfo.Transactions))
	for i, txInfo := range blockInfo.Transactions {
		fmt.Printf("index: %+v\n", i)
		fmt.Printf("txInfo: %+v\n", txInfo)
		if i > 4 {
			break
		}
		for _, txIn := range txInfo.TxIn {
			fmt.Printf("txIn: %+v\n", txIn)
		}
		for _, txOut := range txInfo.TxOut {
			fmt.Printf("txOut: %+v\n", txOut)
			// PKScript -> address, addrType
			addr, sc := srv.GetAddressByPKScript(txOut.PkScript)
			// value
		}
	}
	fmt.Printf("wch------ END\n")
}
