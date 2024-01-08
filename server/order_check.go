package server

import (
	"fmt"
	"goBTC/db/brc20_market"
	"goBTC/global"
	"goBTC/utils/logutils"
)

func QueryPendingOrder4DB() {
	log := global.LOG
	// // 查询pending的订单信息
	// order := &brc20_market.Order{}
	// list, err := order.GetPendingOrder()
	// if err != nil {
	// 	logutils.LogErrorf(log, "GetPendingOrder error: %+v", err)
	// 	return
	// }
	// for _, data := range list {
	// 	fmt.Printf("wch----- data: %+v\n", data)
	// }
	hash := "00ce62a9bf686e13a692d96748c575cbc310c29156821482797245ece0274322"
	data, err := QueryTransferInfo(hash)
	if err != nil {
		logutils.LogErrorf(log, "QueryTransferInfo error: %+v", err)
		return
	}
	logutils.LogInfof(log, "QueryPendingOrder4DB get witness data: %+v", data)
}

func QueryTransferInfo(hash string) (*brc20_market.Order, error) {
	srv := global.Client
	log := global.LOG
	data, err := srv.GetRawTransactionByHash(hash)
	if err != nil {
		logutils.LogErrorf(log, "GetRawTransactionByHash error: %+v", err)
		return nil, err
	}
	if data.BlockHash == "" {
		return nil, fmt.Errorf("the tx not ok")
	}
	// 获取块高
	blockInfo, err := srv.GetBlockStatus(data.BlockHash)
	if err != nil {
		return nil, err
	}
	fmt.Printf("wch---- data: %+v\n", blockInfo)
	// 查询到交易数据,整理铭文信息
	inscriberInfo, err := GetInscribeInfoByHash(data.Hex)
	if err != nil {
		logutils.LogErrorf(log, "GetInscribeInfoByHash error: %+v", err)
		return nil, err
	}
	logutils.LogInfof(log, "body len: %+v", inscriberInfo.ContentSize)
	logutils.LogInfof(log, "Brc20: %+v", inscriberInfo.Brc20.Tick)
	// 校验地址
	// 校验数量
	// 处理状态
	inscribeId := ""
	state := 2
	// 整理订单数据
	order := &brc20_market.Order{
		InscribeID:      &inscribeId,        // 这个
		InscribeContent: inscriberInfo.Body, // 这个
		ContentType:     &inscriberInfo.ContentType,
		BlockHeight:     &blockInfo.Height, // 这个
		Tick:            inscriberInfo.Brc20.Tick,
		State:           state,
	}
	fmt.Printf("wch---- order: %+v\n", order)
	return order, nil
}
