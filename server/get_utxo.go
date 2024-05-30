package server

import (
	"fmt"
	"goBTC/elastic"
	"goBTC/global"
	"goBTC/prometheus"
	"goBTC/utils/logutils"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/btcsuite/btcd/wire"
	"github.com/shopspring/decimal"
)

var (
	// 处理数据
	SpendUTXOMap   sync.Map
	UnspendUTXOMap sync.Map
	AddrUTXO       map[string][]interface{}
	AddrSent       map[string]decimal.Decimal
)

func GetTransferByBlockHeight(startHeight, newHigh int64) {
	srv := global.Client
	log := global.LOG
	logutils.LogInfof(log, "[GetBlockInfo] Start startHeight: %v, newHigh: %v", startHeight, newHigh)
	for i := startHeight; i <= newHigh; i++ {
		startTime := time.Now().Unix()
		// 查询块的数据
		blockInfo, err := srv.GetBlockInfoByHeight(i)
		if err != nil {
			logutils.LogErrorf(log, "GetBlockInfoByHash error: %+v", err)
			i--
			continue
		}
		endTime := time.Now().Unix()
		txInfoLength := len(blockInfo.Transactions)
		logutils.LogInfof(log, "Get block info, block height: [%v], have tx: [%v], time: [%v]", i, txInfoLength, endTime-startTime)
		// 遍历交易获取输入输出
		prometheus.GoroutineUsage.Set(0)
		for _, txInfo := range blockInfo.Transactions {
			// 添加计数器
			Wg.Add(1)
			go GetUTXOInfoByTransferInfo(txInfo)
		}
		Wg.Wait()
		// 处理输入
		AddrSent = make(map[string]decimal.Decimal)
		DealTxInByBlock()
		// for key, val := range AddrSent {
		// 	fmt.Printf("wch---- AddrSent: %v, %v\n", key, val)
		// }
		// 处理输出
		AddrUTXO = make(map[string][]interface{})
		DealTxOutByBlock(startHeight)
		// for key, val := range AddrUTXO {
		// 	fmt.Printf("wch---- AddrUTXO: %v, %+v\n", key, val)
		// }
	}
}

func GetUTXOInfoByTransferInfo(txInfo *wire.MsgTx) {
	// fmt.Printf("txInfo: %+v\n", txInfo)
	defer Wg.Done()
	prometheus.GoroutineUsage.Add(1)
	for _, txIn := range txInfo.TxIn {
		key := txIn.PreviousOutPoint.String()
		SpendUTXOMap.Store(key, txIn)
	}
	txHash := txInfo.TxHash().String()
	for i, txOut := range txInfo.TxOut {
		key := fmt.Sprintf("%s:%d", txHash, i)
		UnspendUTXOMap.Store(key, txOut)
	}
}

func DealTxInByBlock() {
	count := 0
	SpendUTXOMap.Range(func(key, value interface{}) bool {
		count++
		txIn := key.(string)
		// 判断是否有输入在当前区块生成
		if value, ok := UnspendUTXOMap.Load(key); ok {
			// 存在的话这个输出不需要保存
			_, _, _, _, amount := GetUTXOInfoByTxOut(key, value)
			AddrSent[txIn] = amount
			return true
		}
		AddrSent[txIn] = decimal.Zero
		return true
	})
}

func DealTxOutByBlock(height int64) {
	count := 0
	addrUTXO := make(map[string][]interface{})
	UnspendUTXOMap.Range(func(key, value interface{}) bool {
		addr, pkScript, txId, vout, amount := GetUTXOInfoByTxOut(key, value)
		if addr == "" {
			return true
		}
		// UTXOInfo
		UTXOInfo := elastic.UnSpentsUTXO{
			TxId:         txId,
			Vout:         vout,
			ScriptPubKey: pkScript,
			Amount:       amount,
			Height:       height,
		}
		addrUTXO[addr] = append(addrUTXO[addr], UTXOInfo)
		count++
		return true
	})
	AddrUTXO = addrUTXO
}

func GetUTXOInfoByTxOut(key, value interface{}) (string, string, string, int64, decimal.Decimal) {
	zero := decimal.Zero
	txOut := value.(*wire.TxOut)
	// PKScript -> address, addrType
	addr, err := global.Client.GetAddressByPKScript(txOut.PkScript)
	if err != nil {
		return "", "", "", 0, zero
	}
	pkScript := fmt.Sprintf("%x", txOut.PkScript)
	// txInfo
	txInfo := strings.Split(key.(string), ":")
	if len(txInfo) != 2 {
		return "", "", "", 0, zero
	}
	// vout
	vout, err := strconv.ParseInt(txInfo[1], 0, 64)
	if err != nil {
		return "", "", "", 0, zero
	}
	// value sats
	amount := decimal.NewFromInt(txOut.Value)
	return addr, pkScript, txInfo[0], vout, amount
}
