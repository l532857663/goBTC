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
	// global.MysqlFlag = false
	goBTC.MustLoad("./config.yml")
	srv = global.Client
	log = global.LOG
	go CheckNewHeight(767430)
	utils.SignalHandler("scan", goBTC.Shutdown)
}

func CheckNewHeight(startHeight int64) {
	fmt.Println("[CheckNewHeight] Start")
	for {
		newHigh, err := srv.GetBlockCount()
		if err != nil {
			fmt.Printf("GetBlockCount error: %+v\n", err)
			return
		}
		if startHeight > newHigh {
			time.Sleep(5 * time.Minute)
			continue
		}
		GetBlockInfo(startHeight, newHigh)
		startHeight = newHigh + 1
		time.Sleep(5 * time.Minute)
		fmt.Println("[CheckNewHeight] Once time New high:", newHigh)
	}
}

func GetBlockInfo(startHeight, newHigh int64) {
	log.Info("[GetBlockInfo] Start", zap.Any("startHeight", startHeight))
	for i := startHeight; i <= newHigh; i++ {
		startTime := time.Now().Unix()
		blockInfo, err := srv.GetBlockInfoByHeight(i)
		if err != nil {
			log.Error("GetBlockInfoByHash", zap.Error(err))
			i--
			continue
		}
		endTime := time.Now().Unix()
		txInfoLength := len(blockInfo.Transactions)
		log.Info("Get block info", zap.Any("block height", i), zap.Any("have tx", txInfoLength), zap.Any("time", endTime-startTime))
		// log.Info("Get block", zap.Any("header", blockInfo.Header))
		sum, sumBrc20 := 0, 0
		for j := 0; j < txInfoLength; j++ {
			tx := blockInfo.Transactions[j]
			txHash := tx.TxHash().String()
			witnessStr, oldTxid := client.GetWitnessAndTxHashByTxIn(tx)
			// 判断该交易是否存在铭文流转
			txHaveInscribe, err := ord.GetInscribeIsExist(oldTxid)
			if err != nil {
				log.Info("GetInscribeIsExist", zap.Error(err))
			}
			// log.Info("Get tx", zap.Any("txHash", txHash), zap.Any("witness len", len(witnessStr)), zap.Any("txHaveInscribe", txHaveInscribe))
			if witnessStr == "" {
				if !txHaveInscribe {
					continue
				}
			}
			txInfo, err := srv.GetRawTransactionByHash(txHash)
			if err != nil {
				log.Error("GetRawTransactionByHash", zap.Any("txHash", txHash), zap.Error(err))
				j--
				continue
			}
			if witnessStr == "" {
				var err error
				// 添加操作日志
				err = ord.SaveInscribeActivity(oldTxid, nil, txInfo)
				if err != nil {
					log.Error("CreateActivityInfo", zap.Any("oldTxid", oldTxid), zap.Error(err))
					continue
				}
				// 修改铭文拥有人
				err = ord.UpdateInscribeInfoOwner(oldTxid, txInfo)
				if err != nil {
					log.Error("UpdateInscribeInfoOwner", zap.Any("oldTxid", oldTxid), zap.Error(err))
					continue
				}
				continue
			}
			// 该交易存在铭文
			res := client.GetScriptString(witnessStr)
			if res == nil {
				continue
			}
			logStr := fmt.Sprintf("[%d] txHash: %s, [ord] : %v\n", j, txHash, res.ContentType)
			res.TxHaveInscribe = txHaveInscribe
			// 保存铭文数据
			err = ord.SaveInscribeInfoByTxInfo(i, res, txInfo)
			if err != nil {
				log.Error("SaveInscribeInfoByTxInfo", zap.Any("tx", logStr), zap.Error(err))
				continue
			}
			if res.Brc20 != nil && res.Brc20.P != "" {
				// 保存BRC20铭文数据
				err := ord.SaveInscribeBrc20ByTxInfo(i, res, txInfo)
				if err != nil {
					log.Error("SaveInscribeBrc20ByTxInfo", zap.Any("tx", logStr), zap.Error(err))
					continue
				}
				sumBrc20++
			}
			// 添加操作日志
			err = ord.SaveInscribeActivity(oldTxid, res, txInfo)
			if err != nil {
				log.Error("CreateActivityInfo", zap.Any("tx", logStr), zap.Error(err))
				continue
			}
			sum++
		}
		eTime := time.Now().Unix()
		log.Info("Get block", zap.Any("all time", eTime-endTime))
		log.Info("the block get inscribe:", zap.Any("sum", sum), zap.Any("sumBrc20", sumBrc20))
	}
	fmt.Println("[GetBlockInfo] End")
}
