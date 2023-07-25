package oklink

import (
	"fmt"
	"goBTC/models"
	"goBTC/utils/http"
	"net/url"
)

func (p *Platform) GetBalanceByAddress(symbol, address, protocolType string, filter models.Filter) (*models.GetBalanceResp, error) {
	// 查询主币
	res := &models.GetBalanceResp{}
	if filter.Page == "1" {
		data, err := p.GetBalanceByAddressForSymbol(symbol, address)
		if err != nil {
			return nil, err
		}
		res.Balance = data.Balance
		res.BalanceSymbol = data.BalanceSymbol
		res.ChainFullName = data.ChainFullName
		res.ChainShortName = data.ChainShortName
	}
	if protocolType != "" {
		data, err := p.GetBalanceByAddressForToken(symbol, address, protocolType, filter.Page, filter.Limit)
		if err != nil {
			return nil, err
		}
		var tokenList []models.TokenInfo
		for _, token := range data.TokenList {
			tmp := models.TokenInfo{
				Token:                token.Token,
				TokenContractAddress: token.TokenContractAddress,
				Amount:               token.HoldingAmount,
			}
			tokenList = append(tokenList, tmp)
		}
		res.TotalPage = data.TotalPage
		res.ChainFullName = data.ChainFullName
		res.ChainShortName = data.ChainShortName
		res.TokenList = tokenList
	}
	if res.Balance == "" && res.TotalPage == "" {
		return nil, nil
	}
	return res, nil
}

func (p *Platform) GetBalanceByAddressForSymbol(symbol, address string) (*AddressInfo, error) {
	// 查询主币信息
	params := url.Values{}
	params.Add("chainShortName", symbol)
	params.Add("address", address)

	url := OkLinkApiHost + GetBalanceForSymbol + "?" + params.Encode()
	code, body, err := http.HttpGetWithHeader(url, p.HttpHeader)
	if err != nil || code != 200 {
		return nil, err
	}

	var bodyInfo *BaseResp[*AddressInfo]
	data, err := DecodeBodyDataOne([]byte(body), bodyInfo)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (p *Platform) GetBalanceByAddressForToken(symbol, address, protocolType, page, limit string) (*AddressTokenInfo, error) {
	// 查询主币信息
	params := url.Values{}
	params.Add("chainShortName", symbol)
	params.Add("address", address)
	params.Add("protocolType", protocolType)
	params.Add("page", page)
	params.Add("limit", limit)

	url := OkLinkApiHost + GetBalanceForToken + "?" + params.Encode()
	code, body, err := http.HttpGetWithHeader(url, p.HttpHeader)
	if err != nil || code != 200 {
		err := fmt.Errorf("GetBalanceForToken error, code = %v, msg = %v", code, string(body))
		return nil, err
	}

	var bodyInfo *BaseResp[*AddressTokenInfo]
	data, err := DecodeBodyDataOne([]byte(body), bodyInfo)
	if err != nil {
		return nil, err
	}
	return data, nil
}
