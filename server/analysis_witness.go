package server

import (
	"fmt"
	"goBTC/client"
	"goBTC/models"
)

func GetInscribeInfoByHash(txHex string) (*models.OrdInscribeData, error) {
	// 查询Witness的铭文数据
	witness := client.GetTxWitnessByTxHex(txHex)
	if witness == "" {
		return nil, fmt.Errorf("GetTxWitnessByTxHex null")
	}
	resList := client.GetScriptString(witness)
	if resList == nil {
		return nil, fmt.Errorf("GetScriptString error")
	}
	return resList, nil
}
