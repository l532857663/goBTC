package ord

import (
	"fmt"
	"goBTC/elastic"
	"goBTC/models"
	"goBTC/utils"
)

func DealWithBrc20Info(tokenInfo *elastic.OrdToken, brc20 *models.OrdBRC20) error {
	switch brc20.OP {
	case "deploy":
		fmt.Printf("wch----- deploy: %+v\n", brc20)
		err := CheckBrc20ForDeploy(tokenInfo, brc20)
		if err != nil {
			return err
		}
	case "mint":
		fmt.Printf("wch----- mint: %+v\n", brc20)
		// 判断币种mint完没有
		err := CheckBrc20ForMint(tokenInfo, brc20)
		if err != nil {
			return err
		}
	case "transfer":
		fmt.Printf("wch----- transfer: %+v\n", brc20)
	default:
		fmt.Printf("wch----- default: %+v\n", brc20)
	}
	return nil
}

func CheckBrc20ForDeploy(tokenInfo *elastic.OrdToken, brc20 *models.OrdBRC20) error {
	// ID信息
	dId := GetDeployIdStr(brc20.Tick)
	res, err := elastic.GetDataById(elastic.DeployType, dId)
	if err != nil {
		return err
	}
	esId := tokenInfo.InscribeId
	// 判断之前块高的数据有没有
	if res.Id != "" && res.Error == nil {
		// 有数据已存在，新数据为失败
		return UpdateOrdTokenState(esId, elastic.InscriptionStateInvalid)
	}

	// 之前没数据，添加新数据
	ord := &elastic.DeployToken{
		InscribeId: esId,
		TxHash:     tokenInfo.TxHash,
		Tick:       brc20.Tick,
		Lim:        brc20.Lim,
		Max:        brc20.Max,
		DeployAddr: tokenInfo.OwnerAddress,
		DeployTime: tokenInfo.BlockTime,
		State:      elastic.TokenMintedInit,
	}
	err = elastic.CreateData(elastic.DeployType, dId, ord)
	if err != nil {
		return err
	}
	// 更新OrdToken状态
	return UpdateOrdTokenState(esId, elastic.InscriptionStateSuccess)
}

func CheckBrc20ForMint(tokenInfo *elastic.OrdToken, brc20 *models.OrdBRC20) error {
	// ID信息
	dId := GetDeployIdStr(brc20.Tick)
	res, err := elastic.GetDataById(elastic.DeployType, dId)
	if err != nil {
		return err
	}
	esId := tokenInfo.InscribeId
	if res.Id == "" && res.Error != nil {
		// 没有deploy数据，操作失败
		return UpdateOrdTokenState(esId, elastic.InscriptionStateInvalid)
	}
	token := &elastic.DeployToken{}
	err = utils.Map2Struct(res.Source, token)
	if err != nil {
		return err
	}
	fmt.Printf("token: %+v\n", token)
	return nil
}
