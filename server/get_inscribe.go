package server

import (
	"fmt"
	"goBTC/client"
	"goBTC/global"
	"goBTC/ord"
	"goBTC/utils/logutils"
	"sync"
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/wire"
)

var (
	Wg  sync.WaitGroup
	Wg1 sync.WaitGroup
)

func CheckNewHeight(startHeight int64) {
	srv := global.Client
	log := global.LOG
	logutils.LogInfof(log, "[CheckNewHeight] Start")
	for {
		newHigh, err := srv.GetBlockCount()
		if err != nil {
			logutils.LogErrorf(log, "GetBlockCount error: %+v", err)
			return
		}
		if startHeight > newHigh {
			time.Sleep(5 * time.Minute)
			continue
		}
		GetBlockInfo(startHeight, newHigh)
		startHeight = newHigh + 1
		time.Sleep(5 * time.Minute)
		logutils.LogInfof(log, "[CheckNewHeight] Once time New high: %v", newHigh)
	}
}

func GetBlockInfo(startHeight, newHigh int64) {
	srv := global.Client
	log := global.LOG
	logutils.LogInfof(log, "[GetBlockInfo] Start startHeight: %v", startHeight)
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
		// logutils.LogInfof(log, "Get block header", blockInfo.Header)
		sum, sumBrc20 := new(int), new(int)
		for j := 0; j < txInfoLength; j++ {
			Wg.Add(1)
			if j%100 == 0 {
				time.Sleep(1 * time.Second)
			}
			go GetOneTxInfo(blockInfo, sum, sumBrc20, i, j, j)
		}
		Wg.Wait()
		eTime := time.Now().Unix()
		logutils.LogInfof(log, "Get block all time: [%v]", eTime-endTime)
		logutils.LogInfof(log, "the block get inscribe sum: [%v], sumBrc20: [%v]", *sum, *sumBrc20)
	}
	logutils.LogInfof(log, "[GetBlockInfo] End")
}

func GetOneTxInfo(blockInfo *wire.MsgBlock, sum, sumBrc20 *int, i int64, j, flag int) {
	srv := global.Client
	log := global.LOG
	tx := blockInfo.Transactions[j]
	txHash := tx.TxHash().String()
	defer func() {
		Wg.Done()
		if r := recover(); r != nil {
			logutils.LogErrorf(log, "panic error %v, j: %v, txHash: %v", r, j, txHash)
		}
	}()
	witnessStr := client.GetWitnessByTxInFor0(tx)
	// 判断该交易是否存在铭文流转
	txHaveInscribe, inscribeIndex := CheckTransferHaveInscribe(tx)
	if witnessStr == "" {
		if len(txHaveInscribe) == 0 {
			return
		}
	}
	txInfo, err := srv.GetRawTransactionByHash(txHash)
	if err != nil {
		logutils.LogErrorf(log, "GetRawTransactionByHash txHash: %v, error: %+v", txHash, err)
		if flag < j+3 {
			GetOneTxInfo(blockInfo, sum, sumBrc20, i, j, flag+1)
		}
		return
	}
	toAddr := txInfo.Vout[0].ScriptPubKey.Address
	if witnessStr == "" {
		logutils.LogInfof(log, "Get inscribe deal txHash: %v, j: %v", txHash, j)
		// 该交易是正常BTC转账，包含铭文数据
		for i, index := range inscribeIndex {
			Wg1.Add(1)
			go OldInscribeChange(tx, txInfo, index, txHaveInscribe[i])
		}
		Wg1.Wait()
		return
	}
	// 该交易存在铭文
	res := client.GetScriptString(witnessStr)
	if res == nil {
		return
	}
	logutils.LogInfof(log, "Get inscribe new txHaveInscribe: %v, txHash: %v, j: %v", len(txHaveInscribe), txHash, j)
	logStr := fmt.Sprintf("[%d] txHash: %s, [ord] : %v\n", j, txHash, res.ContentType)
	if len(txHaveInscribe) > 0 {
		res.TxHaveInscribe = txHaveInscribe[0]
	}
	// 保存铭文数据
	err = ord.SaveInscribeInfoByTxInfo(i, res, txInfo)
	if err != nil {
		logutils.LogErrorf(log, "SaveInscribeInfoByTxInfo j: %v error: %+v", logStr, err)
		return
	}
	if res.Brc20 != nil && res.Brc20.P != "" {
		// 保存BRC20铭文数据
		err := ord.SaveInscribeBrc20ByTxInfo(i, res, txInfo)
		if err != nil {
			logutils.LogErrorf(log, "SaveInscribeBrc20ByTxInfo j: %v error: %+v", logStr, err)
			return
		}
		*sumBrc20++
	}
	// 添加操作日志
	err = ord.SaveInscribeActivity(res.TxHaveInscribe, toAddr, res, txInfo)
	if err != nil {
		logutils.LogErrorf(log, "CreateActivityInfo j: %v error: %+v", logStr, err)

		return
	}
	*sum++
}

func GetInscribeToAddrByRawTx(tx *wire.MsgTx, inscribeIn int) *int {
	srv := global.Client
	log := global.LOG
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
			logutils.LogErrorf(log, "GetRawTransaction error: %+v, txHash: %v", err, tx.TxHash().String())
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
		logutils.LogErrorf(log, "GetInscribeToAddrByRawTx txHash: %v", tx.TxHash().String())
		return nil
	}
	return &toAddrIndex
}

// 判断交易是否存在铭文
func CheckTransferHaveInscribe(tx *wire.MsgTx) ([]string, []int) {
	log := global.LOG
	var inscribeList []string
	var inscribeIndex []int
	for i, vin := range tx.TxIn {
		oldTxid, output := client.GetVinHashAndOutput(vin.PreviousOutPoint)
		txHaveInscribe, err := ord.GetInscribeIsExist(oldTxid, output)
		if err != nil {
			logutils.LogErrorf(log, "GetInscribeIsExist error: %+v, oldTxid: %v", err, oldTxid)
			continue
		}
		if txHaveInscribe != "" {
			// 获取到铭文信息
			inscribeIndex = append(inscribeIndex, i)
			inscribeList = append(inscribeList, txHaveInscribe)
			logutils.LogInfof(log, "GetInscribeIsExist oldTxid: %v", oldTxid)
			break
		}
	}
	return inscribeList, inscribeIndex
}

// 铭文UTXO转移
func OldInscribeChange(tx *wire.MsgTx, txInfo *btcjson.TxRawResult, inscribeIndex int, txHaveInscribe string) {
	log := global.LOG
	txHash := tx.TxHash().String()
	defer Wg1.Done()
	// 旧铭文ID是创建铭文的ID，非转移后的
	var err error
	// 查询转移铭文的输出地址
	toAddrIndex := GetInscribeToAddrByRawTx(tx, inscribeIndex)
	if toAddrIndex == nil {
		return
	}
	toAddr := txInfo.Vout[*toAddrIndex].ScriptPubKey.Address
	logutils.LogInfof(log, "Get old inscribe txHaveInscribe: %v, toAddr: %v", txHaveInscribe, toAddr)
	// 添加操作日志
	err = ord.SaveInscribeActivity(txHaveInscribe, toAddr, nil, txInfo)
	if err != nil {
		logutils.LogErrorf(log, "CreateActivityInfo txHaveInscribe: %v, txHash: %v, error: %+v", txHaveInscribe, txHash, err)
		return
	}
	// 修改铭文拥有人
	err = ord.UpdateInscribeInfoOwner(txHaveInscribe, txInfo)
	if err != nil {
		logutils.LogErrorf(log, "UpdateInscribeInfoOwner txHaveInscribe: %v, txHash: %v, error: %+v", txHaveInscribe, txHash, err)
		return
	}
	return
}
