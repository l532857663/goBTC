package server

import (
	"fmt"
	"goBTC/global"
	"goBTC/utils/logutils"
	"time"

	"github.com/btcsuite/btcd/wire"
)

var (
	// 处理数据
	SpendUTXOMap   map[string]interface{}
	UnspendUTXOMap map[string]interface{}
)

func GetTransferByBlockHeight(startHeight, newHigh int64) {
	srv := global.Client
	log := global.LOG
	logutils.LogInfof(log, "[GetBlockInfo] Start startHeight: %v, newHigh: %v", startHeight, newHigh)
	for i := startHeight; i <= newHigh; i++ {
		startTime := time.Now().Unix()
		blockInfo, err := srv.GetBlockInfoByHeight(i)
		if err != nil {
			logutils.LogErrorf(log, "GetBlockInfoByHash error: %+v", err)
			i--
			continue
		}
		endTime := time.Now().Unix()
		txInfoLength := len(blockInfo.Transactions)
		logutils.LogInfof(log, "Get block info, block height: [%v], have tx: [%v], time: [%v]", i, txInfoLength, endTime-startTime)
		SpendUTXOMap = make(map[string]interface{}, txInfoLength*10)
		UnspendUTXOMap = make(map[string]interface{}, txInfoLength*10)
		for i, txInfo := range blockInfo.Transactions {
			fmt.Printf("index: %+v\n", i)
			// 添加计数器
			Wg.Add(1)
			go GetUTXOInfoByTransferInfo(txInfo)
		}
		Wg.Wait()
		fmt.Printf("wch---- suMap len: %v, nuMap len: %v\n", len(SpendUTXOMap), len(UnspendUTXOMap))
	}
	fmt.Printf("wch------ END\n")
}

func GetUTXOInfoByTransferInfo(txInfo *wire.MsgTx) {
	fmt.Printf("txInfo: %+v\n", txInfo)
	defer Wg.Done()
	for _, txIn := range txInfo.TxIn {
		SpendUTXOMap[txIn.PreviousOutPoint.String()] = txIn
	}
	txHash := txInfo.TxHash().String()
	for i, txOut := range txInfo.TxOut {
		// fmt.Printf("txOut: %+v\n", txOut)
		// // PKScript -> address, addrType
		// addr, err := srv.GetAddressByPKScript(txOut.PkScript)
		// // value sats
		// amount := txOut.Value
		key := fmt.Sprintf("%s:%d", txHash, i)
		UnspendUTXOMap[key] = txOut
	}
}
