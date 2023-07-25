package oklink

import (
	"fmt"
	"goBTC/global"
	"goBTC/models"
	"goBTC/utils/http"
	"net/url"

	"go.uber.org/zap"
)

func (p *Platform) GetInscriptionList(filter models.InscribeFilter) (*models.GetInscribeResp, error) {
	// 查询主币信息
	params := url.Values{}
	// params.Add("token", filter.Token)
	params.Add("inscriptionId", filter.InscriptionId)
	// params.Add("inscriptionNumber", filter.InscriptionNumber)
	// params.Add("state", filter.State)
	// params.Add("page", filter.Page)
	// params.Add("limit", filter.Limit)

	url := OkLinkApiHost + GetInscriptionList + "?" + params.Encode()
	code, body, err := http.HttpGetWithHeader(url, p.HttpHeader)
	if err != nil {
		global.LOG.Error("Http Get error", zap.Error(err))
		return nil, err
	}
	if code != 200 {
		err = fmt.Errorf("Http get nothing, code: %v", code)
		global.LOG.Error("Http Get", zap.Any("code", code), zap.Any("body", body))
		return nil, err
	}

	fmt.Printf("wch------- body: %+v\n", body)
	var bodyInfo *BaseResp[*InscriptionsInfo]
	data, err := DecodeBodyDataOne([]byte(body), bodyInfo)
	if err != nil {
		return nil, err
	}
	inscribptionList := []models.InscriptionsList{}
	for _, info := range data.InscriptionsList {
		tmp := models.InscriptionsList{
			InscriptionId:     info.InscriptionId,
			InscriptionNumber: info.InscriptionNumber,
			Location:          info.Location,
			Token:             info.Token,
			State:             info.State,
			Msg:               info.Msg,
			TokenType:         info.TokenType,
			ActionType:        info.ActionType,
			LogoURL:           info.LogoURL,
			OwnerAddress:      info.OwnerAddress,
			TxID:              info.TxID,
			BlockHeight:       info.BlockHeight,
			ContentSize:       info.ContentSize,
			Time:              info.Time,
		}
		inscribptionList = append(inscribptionList, tmp)
	}
	res := &models.GetInscribeResp{
		Page:             data.Page,
		Limit:            data.Limit,
		TotalPage:        data.TotalPage,
		TotalInscription: data.TotalInscription,
		InscriptionsList: inscribptionList,
	}
	return res, nil
}
