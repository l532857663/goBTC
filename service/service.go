package service

import (
	"fmt"
	"goBTC/global"
	"goBTC/models"
	"goBTC/service/oklink"
	"strings"

	"go.uber.org/zap"
)

// const (
// 	OkLink = "OkLink"
// )

var (
	// platformName
	OkLink string = "OkLink"
	// 使用查询平台及开关
	PlatformInfo = map[string]bool{
		OkLink: true,
	}
	// 全局可使用的平台映射
	PlatformMap   = map[string]Platform{}
	ChainMainCoin = map[string]string{
		"BTC":  "BTC",
		"ETH":  "ETH",
		"BSC":  "BNB",
		"TRON": "TRX",
	}
)

func InitPlatformMap() {
	// 初始化全局可使用的平台信息
	for platform, ok := range PlatformInfo {
		if !ok {
			continue
		}
		switch platform {
		case OkLink:
			PlatformMap[platform] = oklink.NewPlatformInfo()
		default:
			global.LOG.Warn("Not support the platform!", zap.String("name", platform))
		}
	}
}

func GetBalanceInfo(req GetAddressInfoReq) *models.GetBalanceResp {
	// 使用多平台查询资产信息
	symbol := strings.ToUpper(req.Symbol)
	for platform, obj := range PlatformMap {
		global.LOG.Info("use the", zap.String("platform", platform), zap.Any("req", req))
		filter := models.Filter{
			Page:  req.PageIndex,
			Limit: req.PageSize,
		}
		client := obj.(GetBalance)
		res, err := client.GetBalanceByAddress(symbol, req.Address, req.ProtocolType, filter)
		if err != nil {
			global.LOG.Error("The platform get balance error", zap.Any("platform", platform))
			continue
		}
		fmt.Printf("res: %+v\n", res)
		if res != nil {
			return res
		}
	}
	coin, _ := ChainMainCoin[symbol]
	res := &models.GetBalanceResp{
		Balance:       "0",
		BalanceSymbol: coin,
		TotalPage:     "0",
		TokenList:     NilTokenList,
	}
	return res
}

func GetTransferInfoByAddress(req GetAddressInfoReq) *models.GetTransferResp {
	// 使用多平台查询资产信息
	symbol := strings.ToUpper(req.Symbol)
	for platform, obj := range PlatformMap {
		global.LOG.Info("use the", zap.String("platform", platform), zap.Any("req", req))
		filter := models.Filter{
			Page:            req.PageIndex,
			Limit:           req.PageSize,
			ContractAddress: req.TokenContractAddress,
		}
		client := obj.(GetTransfer)
		res, err := client.GetTransferByAddress(symbol, req.Address, req.ProtocolType, filter)
		if err != nil {
			global.LOG.Error("The platform get balance error", zap.Any("platform", platform))
			continue
		}
		fmt.Printf("res: %+v\n", res)
		if res != nil {
			return res
		}
	}
	res := NilTransferResp
	return res
}

func GetTransferInfoForUTXO(req GetUnspentReq) []*models.GetTransferUTXOResp {
	// 使用多平台查询资产信息
	symbol := strings.ToUpper(req.Symbol)
	for platform, obj := range PlatformMap {
		global.LOG.Info("use the", zap.String("platform", platform), zap.Any("req", req))
		client := obj.(GetTransferUTXO)
		res, err := client.GetTransferUTXOByAddress(symbol, req.Address)
		if err != nil {
			global.LOG.Error("The platform get balance error", zap.Any("platform", platform))
			continue
		}
		fmt.Printf("res: %+v\n", res)
		if res != nil {
			return res
		}
	}
	return nil
}

func GetTransferInfoForBlock(req GetTransferReq) *models.GetTransferResp {
	// 使用多平台查询资产信息
	symbol := strings.ToUpper(req.Symbol)
	for platform, obj := range PlatformMap {
		global.LOG.Info("use the", zap.String("platform", platform), zap.Any("req", req))
		filter := models.Filter{
			Page:  req.PageIndex,
			Limit: req.PageSize,
		}
		client := obj.(GetTransfer)
		res, err := client.GetTransferByBlockNum(symbol, req.Height, req.ProtocolType, filter)
		if err != nil {
			global.LOG.Error("The platform get balance error", zap.Any("platform", platform))
			continue
		}
		if res != nil {
			return res
		}
	}
	res := NilTransferResp
	return res
}

func GetInscriptionInfo(req GetInscriptionReq) *models.GetInscribeResp {
	// 使用多平台查询资产信息
	// symbol := strings.ToUpper(req.Symbol)
	for platform, obj := range PlatformMap {
		global.LOG.Info("use the", zap.String("platform", platform), zap.Any("req", req))
		filter := models.InscribeFilter{
			Page:              req.PageIndex,
			Limit:             req.PageSize,
			Token:             req.Token,
			InscriptionId:     req.InscriptionId,
			InscriptionNumber: req.InscriptionNumber,
			State:             req.State,
		}
		client := obj.(GetInscriptions)
		res, err := client.GetInscriptionList(filter)
		if err != nil {
			global.LOG.Error("The platform get balance error", zap.Any("platform", platform))
			continue
		}
		if res != nil {
			return res
		}
	}
	res := NilInscribeResp
	return res
}
