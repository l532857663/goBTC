package main

import (
	"fmt"
	"goBTC"
	"goBTC/global"
	"goBTC/server"
	"goBTC/utils"
	"goBTC/utils/logutils"
)

func main() {
	fmt.Println("vim-go")
	// global.MysqlFlag = false
	goBTC.MustLoad("./config.yml")
	go server.CheckNewHeight(2543620)
	// TestGetTx()
	utils.SignalHandler("scan", goBTC.Shutdown)
}

func TestGetTx() {
	fmt.Printf("wch------ test\n")
	blockList := []int64{788344, 789649, 789792, 789793}
	txList := []int{2826, 2398, 1653, 309}
	for index, i := range blockList {
		blockInfo, err := global.Client.GetBlockInfoByHeight(i)
		if err != nil {
			logutils.LogErrorf(global.Log, "GetBlockInfoByHash error: %+v", err)
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
