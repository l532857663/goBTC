package oklink

import (
	"fmt"
	"goBTC/global"
	"goBTC/models"
	"goBTC/utils/http"
	"net/url"

	"go.uber.org/zap"
)

func (p *Platform) GetTransaction(askPath string, paramMap map[string]string) (*TransactionInfo, error) {
	// 查询交易信息
	params := url.Values{}
	for k, v := range paramMap {
		params.Add(k, v)
	}

	askURL := OkLinkApiHost + askPath + "?" + params.Encode()
	code, body, err := http.HttpGetWithHeader(askURL, p.HttpHeader)
	if err != nil {
		global.LOG.Error("Http Get error", zap.Error(err))
		return nil, err
	}
	if code != 200 {
		err = fmt.Errorf("Http get nothing, code: %v", code)
		global.LOG.Error("Http Get", zap.Any("code", code), zap.Any("body", body))
		return nil, err
	}

	var bodyInfo *BaseResp[*TransactionInfo]
	data, err := DecodeBodyDataOne([]byte(body), bodyInfo)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// 查询地址交易记录
func (p *Platform) GetTransferByAddress(symbol, address, protocolType string, filter models.Filter) (*models.GetTransferResp, error) {
	param := map[string]string{
		"chainShortName":       symbol,
		"address":              address,
		"protocolType":         protocolType,
		"tokenContractAddress": filter.ContractAddress,
		"page":                 filter.Page,
		"limit":                filter.Limit,
	}
	data, err := p.GetTransaction(GetAddressTransaction, param)
	if err != nil {
		return nil, err
	}
	transferList := []models.TransferInfo{}
	for _, info := range data.TransactionList {
		tmp := models.TransferInfo{
			Txid:                 info.Txid,
			BlockHash:            info.BlockHash,
			Height:               info.Height,
			TransactionTime:      info.TransactionTime,
			From:                 info.From,
			To:                   info.To,
			IsFromContract:       info.IsFromContract,
			IsToContract:         info.IsToContract,
			FromTag:              info.FromTag,
			ToTag:                info.ToTag,
			Amount:               info.Amount,
			TransactionSymbol:    info.TransactionSymbol,
			TokenContractAddress: info.TokenContractAddress,
			TxFee:                info.TxFee,
			State:                info.State,
		}
		transferList = append(transferList, tmp)
	}
	res := &models.GetTransferResp{
		ChainFullName:  data.ChainFullName,
		ChainShortName: data.ChainShortName,
		TotalPage:      data.TotalPage,
		TransferList:   transferList,
	}
	return res, nil
}

// 查询区块交易记录
func (p *Platform) GetTransferByBlockNum(symbol, height, protocolType string, filter models.Filter) (*models.GetTransferResp, error) {
	param := map[string]string{
		"chainShortName": symbol,
		"height":         height,
		"protocolType":   protocolType,
		"page":           filter.Page,
		"limit":          filter.Limit,
	}
	data, err := p.GetTransaction(GetBlockTransaction, param)
	if err != nil {
		return nil, err
	}
	transferList := []models.TransferInfo{}
	for _, info := range data.BlockList {
		tmp := models.TransferInfo{
			Txid:                 info.Txid,
			BlockHash:            info.BlockHash,
			Height:               info.Height,
			TransactionTime:      info.TransactionTime,
			From:                 info.From,
			To:                   info.To,
			IsFromContract:       info.IsFromContract,
			IsToContract:         info.IsToContract,
			FromTag:              info.FromTag,
			ToTag:                info.ToTag,
			Amount:               info.Amount,
			TransactionSymbol:    info.TransactionSymbol,
			TokenContractAddress: info.TokenContractAddress,
			TxFee:                info.TxFee,
			State:                info.State,
			TokenID:              info.TokenID,
			MethodId:             info.MethodId,
		}
		transferList = append(transferList, tmp)
	}
	res := &models.GetTransferResp{
		ChainFullName:  data.ChainFullName,
		ChainShortName: data.ChainShortName,
		TotalPage:      data.TotalPage,
		TransferList:   transferList,
	}
	return res, nil
}
