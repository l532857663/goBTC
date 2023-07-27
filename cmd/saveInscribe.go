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
	GetInscribptionByTxhashAndHeight()
	if global.MysqlFlag {
		utils.SignalHandler("brc20Assets", goBTC.Shutdown)
	}
}

func GetInscribptionByTxhashAndHeight() {
	//txHash := "5ee59cb5f2b88d1aa1dd7ef0f6263a2682412866e8cdb73275fa013429169623" // 创建BRC20
	//txHash := "231746e07440a6fa81d45f0d26e0510329175de1cac07b64c0a53faafb3b551d" // 转移BRC20
	tmp := []string{
		"b61b0172d95e266c18aea0c624db987e971a5d6d4ebc2aaed85da4642d635735",
		"24f2585e667e345c7b72a4969b4c70eb0e2106727d876217497c6cf86a8a354c",
		"885441055c7bb5d1c54863e33f5c3a06e5a14cc4749cb61a9b3ff1dbe52a5bbb",
		"628f019c4e3c30ccc0fd9aae872cb3720294a255127292bf61c38fbee39462fe",
		"c51c79684743c7cbc53189c17479439320c2f3996ef069d13674ef12ff9260c2",
		"778bf74299ba8b29df3fcf22ce66cdc45a87c21cc229eb0f6d86bc57539971d2",
		"fb15db73d1871165a3bb93ed22dc26da97f71095deb504fd033c0e53a3d7e712",
		"d7fd0b111cf6a7ffc4efd2d5102a48ba513f635ad8b325afce9984f2beef3e5b",
		"095357d354ea7a85d61c677bf93ebf30af077157bd11ee48f074d063082325b9",
		"ccc1702bc70ae324cba9f9a834b982a6437cb70e3dbd1c4275b0ea4407420ab0",
		"e2f3efce6f9f126a97e220f878154797bab204f218748cfd6b2299929cdae5b0",
	}
	b := []int64{
		779832,
		779833,
		779834,
		779835,
		779918,
		779960,
		779960,
		779960,
		779960,
		779960,
		779961,
	}
	for i, v := range tmp {
		GetHashInfo(v, b[i])
	}
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
	// fmt.Printf("wch-----witnessStr: %+v, tx: %+v\n", len(witnessStr), txHaveInscribe)
	if witnessStr == "" {
		if !txHaveInscribe {
			log.Info("This tx not have inscription")
			return
		}
		var err error
		// 添加操作日志
		err = ord.SaveInscribeActivity(oldTxid, nil, txInfo)
		if err != nil {
			log.Error("CreateActivityInfo", zap.Any("oldTxid", oldTxid), zap.Error(err))
			return
		}
		// 修改铭文拥有人
		err = ord.UpdateInscribeInfoOwner(oldTxid, txInfo)
		if err != nil {
			log.Error("UpdateInscribeInfoOwner", zap.Any("oldTxid", oldTxid), zap.Error(err))
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
	err = ord.SaveInscribeActivity(oldTxid, res, txInfo)
	if err != nil {
		log.Error("CreateActivityInfo", zap.Any("oldTxid", oldTxid), zap.Error(err))
		return
	}
	fmt.Println("[GetHashInfo] End")
}
