package main

import (
	"fmt"
	"goBTC"
	"goBTC/client"
	"goBTC/global"
	"goBTC/ord"
	"goBTC/utils"
	"sync"
	"time"

	"github.com/btcsuite/btcd/wire"
	"go.uber.org/zap"
)

var (
	srv *client.BTCClient
	log *zap.Logger
	wg  sync.WaitGroup
)

func main() {
	fmt.Println("vim-go")
	// global.MysqlFlag = false
	goBTC.MustLoad("./config.yml")
	srv = global.Client
	log = global.LOG
	go CheckNewHeight(789792)
	// TestGetTx()
	utils.SignalHandler("scan", goBTC.Shutdown)
}

func CheckNewHeight(startHeight int64) {
	fmt.Println("[CheckNewHeight] Start")
	for {
		newHigh, err := srv.GetBlockCount()
		if err != nil {
			log.Error("GetBlockCount error", zap.Error(err))
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
		sum, sumBrc20 := new(int), new(int)
		for j := 0; j < txInfoLength; j++ {
			wg.Add(1)
			if j%100 == 0 {
				time.Sleep(2 * time.Second)
			}
			go GetOneTxInfo(blockInfo, sum, sumBrc20, i, j, j)
		}
		wg.Wait()
		eTime := time.Now().Unix()
		log.Info("Get block", zap.Any("all time", eTime-endTime))
		log.Info("the block get inscribe:", zap.Any("sum", sum), zap.Any("sumBrc20", sumBrc20))
	}
	log.Info("[GetBlockInfo] End")
}

func GetOneTxInfo(blockInfo *wire.MsgBlock, sum, sumBrc20 *int, i int64, j, flag int) {
	tx := blockInfo.Transactions[j]
	txHash := tx.TxHash().String()
	defer func() {
		wg.Done()
		if r := recover(); r != nil {
			log.Error("panic error", zap.Any("err", r), zap.Any("j", j), zap.Any("txHash", txHash))
		}
	}()
	witnessStr := client.GetWitnessByTxInFor0(tx)
	// 判断该交易是否存在铭文流转
	txHaveInscribe, oldTxid := "", ""
	var inscribeIndex int
	for i, vin := range tx.TxIn {
		output := ""
		oldTxid, output = client.GetVinHashAndOutput(vin.PreviousOutPoint)
		var err error
		txHaveInscribe, err = ord.GetInscribeIsExist(oldTxid, output)
		if err != nil {
			log.Info("GetInscribeIsExist", zap.Error(err), zap.Any("oldTxid", oldTxid))
			continue
		}
		if txHaveInscribe != "" {
			// 获取到铭文信息
			inscribeIndex = i
			break
		}
	}
	if witnessStr == "" {
		if txHaveInscribe == "" {
			return
		}
	}
	txInfo, err := srv.GetRawTransactionByHash(txHash)
	if err != nil {
		log.Error("GetRawTransactionByHash", zap.Any("txHash", txHash), zap.Error(err))
		if flag < j+3 {
			GetOneTxInfo(blockInfo, sum, sumBrc20, i, j, flag+1)
		}
		return
	}
	toAddr := txInfo.Vout[0].ScriptPubKey.Address
	if witnessStr == "" {
		// 旧铭文ID是创建铭文的ID，非转移后的
		var err error
		// 查询转移铭文的输出地址
		toAddrIndex := GetInscribeToAddrByRawTx(tx, inscribeIndex)
		if toAddrIndex == nil {
			return
		}
		toAddr = txInfo.Vout[*toAddrIndex].ScriptPubKey.Address
		log.Info("Get inscribe deal", zap.Any("txHaveInscribe", txHaveInscribe), zap.Any("oldTxid", oldTxid), zap.Any("txHash", txHash), zap.Any("toAddr", toAddr), zap.Any("j", j))
		// 添加操作日志
		err = ord.SaveInscribeActivity(txHaveInscribe, toAddr, nil, txInfo)
		if err != nil {
			log.Error("CreateActivityInfo", zap.Any("txHaveInscribe", txHaveInscribe), zap.Any("txHash", txHash), zap.Error(err))
			return
		}
		// 修改铭文拥有人
		err = ord.UpdateInscribeInfoOwner(txHaveInscribe, txInfo)
		if err != nil {
			log.Error("UpdateInscribeInfoOwner", zap.Any("txHaveInscribe", txHaveInscribe), zap.Any("txHash", txHash), zap.Error(err))
			return
		}
		return
	}
	// 该交易存在铭文
	res := client.GetScriptString(witnessStr)
	if res == nil {
		return
	}
	log.Info("Get inscribe new", zap.Any("txHaveInscribe", txHaveInscribe), zap.Any("oldTxid", oldTxid), zap.Any("txHash", txHash), zap.Any("j", j))
	logStr := fmt.Sprintf("[%d] txHash: %s, [ord] : %v\n", j, txHash, res.ContentType)
	res.TxHaveInscribe = txHaveInscribe
	// 保存铭文数据
	err = ord.SaveInscribeInfoByTxInfo(i, res, txInfo)
	if err != nil {
		log.Error("SaveInscribeInfoByTxInfo", zap.Any("tx", logStr), zap.Error(err))
		return
	}
	if res.Brc20 != nil && res.Brc20.P != "" {
		// 保存BRC20铭文数据
		err := ord.SaveInscribeBrc20ByTxInfo(i, res, txInfo)
		if err != nil {
			log.Error("SaveInscribeBrc20ByTxInfo", zap.Any("tx", logStr), zap.Error(err))
			return
		}
		*sumBrc20++
	}
	// 添加操作日志
	err = ord.SaveInscribeActivity(txHaveInscribe, toAddr, res, txInfo)
	if err != nil {
		log.Error("CreateActivityInfo", zap.Any("tx", logStr), zap.Error(err))
		return
	}
	*sum++
}

func GetInscribeToAddrByRawTx(tx *wire.MsgTx, inscribeIn int) *int {
	// 根据铭文输入的位置，定位铭文输出的位置(通过ordi规则，sat位置匹配原则)
	var inSum int64 = 0
	var fromAddrIndex int
	for i, vin := range tx.TxIn {
		fromAddrIndex = i
		if i == inscribeIn {
			break
		}
		txInfo, err := srv.Client.GetRawTransaction(&vin.PreviousOutPoint.Hash)
		if err != nil {
			log.Error("GetRawTransaction error", zap.Error(err), zap.Any("txHash", tx.TxHash().String()))
			return nil
		}
		inSum += txInfo.MsgTx().TxOut[vin.PreviousOutPoint.Index].Value
	}
	var outSum int64 = 0
	var toAddrIndex int
	for i, vout := range tx.TxOut {
		outSum += vout.Value
		if outSum > inSum {
			toAddrIndex = i
			break
		}
	}
	if fromAddrIndex != 0 && toAddrIndex == 0 {
		log.Error("GetInscribeToAddrByRawTx error", zap.Any("txHash", tx.TxHash().String()))
		return nil
	}
	return &toAddrIndex
}

func TestGetTx() {
	fmt.Printf("wch------ test\n")
	blockList := []int64{788344, 789649, 789792, 789793}
	txList := []int{2826, 2398, 1653, 309}
	for index, i := range blockList {
		blockInfo, err := srv.GetBlockInfoByHeight(i)
		if err != nil {
			log.Error("GetBlockInfoByHash", zap.Error(err))
			return
		}
		fmt.Printf("wch------ test1\n")
		sum, sumBrc20 := new(int), new(int)
		j := txList[index]
		wg.Add(1)
		GetOneTxInfo(blockInfo, sum, sumBrc20, i, j, j)
		fmt.Printf("wch------ END\n")
	}
}
