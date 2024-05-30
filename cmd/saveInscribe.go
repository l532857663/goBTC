package main

import (
	"fmt"
	"goBTC"
	"goBTC/client"
	"goBTC/global"
	"goBTC/ord"
	"goBTC/utils"
	"time"

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
	GetInscribptionByTxhashAndHeight()
	// FixIncribeMoreTransferError()
	if global.MysqlFlag {
		utils.SignalHandler("brc20Assets", goBTC.Shutdown)
	}
}

func GetInscribptionByTxhashAndHeight() {
	//txHash := "5ee59cb5f2b88d1aa1dd7ef0f6263a2682412866e8cdb73275fa013429169623" // 创建BRC20
	//txHash := "231746e07440a6fa81d45f0d26e0510329175de1cac07b64c0a53faafb3b551d" // 转移BRC20
	tmp := []string{
		"c4270c5f80216758bc6f5ad1c169cfc433392625517118fd806212d1e9c7a219",
	}
	b := []int64{
		2568712,
	}
	for i, v := range tmp {
		GetHashInfo(v, b[i])
	}
}

func FixIncribeMoreTransferError() {
	tmp := []string{
		"d137417a11dfa300c57ac0e2e382f68dba1b14e6ab697eccbe5463b6c848fb8d",
		"18403eef89b824fbb4de5a1a3b5dad2f78b6e5ff67086308c94d59f7da5e1675",
		"156a787e50ef1b912830ca1f495fbb2032a159d477844aa84dd63f777bb4db33",
		"de1c1bacf4311067ab547a658afe47dbaf361423e62695bff4d5eac2a4124c9e",
	}
	for _, v := range tmp {
		// 删除旧的Activity
		err := ord.DeleteInscribeActivity(v)
		if err != nil {
			fmt.Printf("ord.DeleteInscribeActivity error: %+v, txHash: %+v\n", err, v)
			continue
		}
		// 重新添加
		GetHashInfo(v, 0)
	}
	fmt.Println("FixIncribeMoreTransferError END")
}

func GetHashInfo(txHash string, blockHeight int64) {
	log.Info("[GetHashInfo] Start", zap.Any("txHash", txHash))
	txInfo, err := srv.GetRawTransactionByHash(txHash)
	if err != nil {
		log.Info("GetRawTransactionByHash", zap.Error(err))
		return
	}
	// fmt.Printf("wch---- txInfo: %+v\n", txInfo)
	witnessStr := client.GetTxWitnessByTxInfo(txInfo)
	// 判断该交易是否存在铭文流转
	oldTxid := txInfo.Vin[0].Txid
	txHaveInscribe, err := ord.GetInscribeIsExist(oldTxid)
	if err != nil {
		log.Info("GetInscribeIsExist", zap.Error(err))
	}
	fmt.Printf("wch-----witnessStr: %+v, tx: %+v\n", len(witnessStr), txHaveInscribe)
	if witnessStr == "" {
		if txHaveInscribe == "" {
			log.Info("This tx not have inscription")
			return
		}
		var err error
		// 添加操作日志
		err = ord.SaveInscribeActivity(txHaveInscribe, nil, txInfo)
		if err != nil {
			log.Error("CreateActivityInfo", zap.Any("txHaveInscribe", txHaveInscribe), zap.Error(err))
			return
		}
		// 修改铭文拥有人
		err = ord.UpdateInscribeInfoOwner(txHaveInscribe, txInfo)
		if err != nil {
			log.Error("UpdateInscribeInfoOwner", zap.Any("txHaveInscribe", txHaveInscribe), zap.Error(err))
			return
		}
		return
	}
	res := client.GetScriptString(witnessStr)
	if res == nil {
		log.Info("GetScriptString not have inscription")
		return
	}
	res.TxHaveInscribe = txHaveInscribe
	err = ord.SaveInscribeInfoByTxInfo(blockHeight, res, txInfo)
	if err != nil {
		log.Error("SaveInscribeInfoByTxInfo", zap.Error(err))
		return
	}
	if res.Brc20 != nil && res.Brc20.P != "" {
		err := ord.SaveInscribeBrc20ByTxInfo(blockHeight, res, txInfo)
		if err != nil {
			log.Error("SaveInscribeBrc20ByTxInfo", zap.Error(err))
			return
		}
	}
	// 添加操作日志
	err = ord.SaveInscribeActivity(txHaveInscribe, res, txInfo)
	if err != nil {
		log.Error("CreateActivityInfo", zap.Any("txHaveInscribe", txHaveInscribe), zap.Error(err))
		return
	}
	fmt.Println("[GetHashInfo] End")
	time.Sleep(1 * time.Second)
}
