package server

import (
	"fmt"
	"goBTC/db/brc20_market"
	"goBTC/global"
	"goBTC/utils/logutils"
	"strconv"
)

func QueryPendingOrder4DB() {
	log := global.LOG
	// 查询pending的订单信息
	order := &brc20_market.Order{}
	list, err := order.GetPendingOrder()
	if err != nil {
		logutils.LogErrorf(log, "GetPendingOrder error: %+v", err)
		return
	}
	for _, data := range list {
		fmt.Printf("wch----- data: %+v\n", data)
		order1, err := QueryTransferInfo(data.TxHash, *data.InscribeID)
		if err != nil {
			logutils.LogErrorf(log, "QueryTransferInfo error: %+v", err)
			return
		}
		logutils.LogInfof(log, "QueryPendingOrder4DB get witness data: %+v", order1)
		ok := CheckOrderInfo(data, order1)
		if ok {
			// 处理状态
			switch data.State {
			case 2:
				data.State = 3
			case 5:
				data.State = 6
			default:
				logutils.LogErrorf(log, "QueryPendingOrder4DB order state error: %+v", data.State)
				return
			}
			// 更新数据库
			row, err := data.UpdatePendingOrderState()
			if err != nil {
				logutils.LogErrorf(log, "UpdatePendingOrderState error: %+v", err)
				return
			}
			logutils.LogInfof(log, "data.UpdatePendingOrderState row: %+v", row)
		}
	}
}

func QueryTransferInfo(hash, inscribeId string) (*brc20_market.Order, error) {
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
	// 查询到交易数据,整理铭文信息
	inscriberInfo, err := GetInscribeInfoByHash(data.Hex)
	if err != nil {
		logutils.LogErrorf(log, "GetInscribeInfoByHash error: %+v", err)
		return nil, err
	}
	logutils.LogInfof(log, "body len: %+v", inscriberInfo.ContentSize)
	logutils.LogInfof(log, "Brc20: %+v", inscriberInfo.Brc20.Tick)
	logutils.LogInfof(log, "Amount: %+v", inscriberInfo.Brc20.Amt)
	// 整理订单数据
	amount, _ := strconv.ParseInt(inscriberInfo.Brc20.Amt, 0, 64)

	order := &brc20_market.Order{
		InscribeID: &inscribeId,
		Tick:       inscriberInfo.Brc20.Tick,
		Amount:     amount,
	}
	GetInscriptionInfoByOrdinals(order)
	return order, nil
}

func CheckOrderInfo(order, checkOrder *brc20_market.Order) bool {
	log := global.LOG
	funcName := "CheckOrderInfo"
	// inscribeId
	if *order.InscribeID != *checkOrder.InscribeID {
		logutils.LogErrorf(log, "[%s]: %+v, %+v", funcName, order.InscribeID, checkOrder.InscribeID)
		return false
	}
	// tick
	if order.Tick != checkOrder.Tick {
		logutils.LogErrorf(log, "[%s]: %+v, %+v", funcName, order.Tick, checkOrder.Tick)
		return false
	}
	// number
	if order.Amount != checkOrder.Amount {
		logutils.LogErrorf(log, "[%s]: %+v, %+v", funcName, order.Amount, checkOrder.Amount)
		return false
	}
	return true
}
