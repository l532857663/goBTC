package server

import (
	"goBTC/db/brc20_market"
	"goBTC/utils/logutils"
)

func SaveOrder() {
	// 保存铭文订单信息
	order := &brc20_market.InscribeConfig{
		ProtocolName: "",
		InscribeType: "",
		InscribeID:   nil,
		Tick:         "",
	}
	list, err := order.Create()
	if err != nil {
		logutils.LogErrorf(log, "GetPendingOrder error: %+v", err)
		return
	}
}
