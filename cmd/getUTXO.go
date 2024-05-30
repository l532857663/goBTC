package main

import (
	"fmt"
	"goBTC"
	"goBTC/client"
	"goBTC/global"
	"goBTC/server"
	"goBTC/utils"
	"goBTC/utils/logutils"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
)

var (
	// 全局参数
	srv *client.BTCClient
	log *zap.Logger
	h   int64
)

func main() {
	fmt.Println("vim-go")
	global.MysqlFlag = false
	goBTC.MustLoad("./config.yml")
	srv = global.Client
	log = global.LOG

	// TestGetTx()
	go server.CheckNewHeight(0, server.GetTransferByBlockHeight)
	go InitHttpService()
	utils.SignalHandler("getUTXO", goBTC.Shutdown)
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

func TestGetBlock(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("wch------ Start\n")
	hh := r.FormValue("h")
	if hh != "" {
		h, _ := strconv.ParseInt(hh, 0, 64)
		fmt.Printf("wch----- h: %v\n", hh)
		server.GetTransferByBlockHeight(h, h)
	}
	fmt.Printf("wch------ END\n")
	w.Write([]byte("OK"))
}

func InitHttpService() {
	http.HandleFunc("/getBlock", TestGetBlock)
	logutils.LogInfof(global.LOG, "InitHttpService: %+v", global.CONFIG.Service.ServiceAddr)
	http.ListenAndServe(global.CONFIG.Service.ServiceAddr, nil)
}
