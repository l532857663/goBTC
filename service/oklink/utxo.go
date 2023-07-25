package oklink

import (
	"goBTC/global"
	"goBTC/models"
	"goBTC/utils/http"
	"net/url"

	"go.uber.org/zap"
)

func (p *Platform) GetTransferUTXOByAddress(symbol, address string) ([]*models.GetTransferUTXOResp, error) {
	// 查询交易记录
	data, err := p.GetTransactionByAddressForUTXO(symbol, address)
	if err != nil {
		return nil, err
	}
	res := make([]*models.GetTransferUTXOResp, len(data))
	for i, info := range data {
		res[i] = &models.GetTransferUTXOResp{
			ChainFullName:    info.ChainFullName,
			ChainShortName:   info.ChainShortName,
			Txid:             info.Txid,
			Height:           info.Height,
			Amount:           info.Amount,
			Address:          info.Address,
			Unspent:          info.Unspent,
			Confirm:          info.Confirm,
			Index:            info.Index,
			TransactionIndex: info.TransactionIndex,
			Balance:          info.Balance,
			Symbol:           info.Symbol,
		}
	}
	return res, nil
}

func (p *Platform) GetTransactionByAddressForUTXO(symbol, address string) ([]*TransactionUTXO, error) {
	// 查询交易信息
	params := url.Values{}
	params.Add("chainShortName", symbol)
	params.Add("address", address)

	url := OkLinkApiHost + GetTransactionForUTXO + "?" + params.Encode()
	code, body, err := http.HttpGetWithHeader(url, p.HttpHeader)
	if err != nil || code != 200 {
		global.LOG.Error("Http Get error", zap.Error(err))
		return nil, err
	}

	var bodyInfo *BaseResp[*TransactionUTXO]
	data, err := DecodeBodyDataAll([]byte(body), bodyInfo)
	if err != nil {
		return nil, err
	}
	return data, nil
}
