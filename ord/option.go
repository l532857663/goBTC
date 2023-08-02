package ord

import (
	"goBTC/elastic"

	"github.com/btcsuite/btcd/btcjson"
)

func UpdateInscribeInfoOwner(esId string, txInfo *btcjson.TxRawResult) error {
	// 处理数据
	addr := txInfo.Vout[0].ScriptPubKey.Address
	vout := GetInscribeOutputStr(txInfo.Txid)

	updateInfo := elastic.UpdateInfo{}
	updateInfo.Doc = make(map[string]interface{})
	updateInfo.Doc["owner_address"] = addr
	updateInfo.Doc["owner_output"] = vout
	err := elastic.UpdateData(elastic.InscribeInfoType, esId, updateInfo)
	if err != nil {
		return err
	}
	return nil
}

func DeleteInscribeActivity(txId string) error {
	err := elastic.DeleteData(elastic.ActivityType, txId, nil)
	if err != nil {
		return err
	}
	return nil
}

func UpdateOrdTokenState(esId, state string) error {
	// 处理数据
	updateInfo := elastic.UpdateInfo{}
	updateInfo.Doc = make(map[string]interface{})
	updateInfo.Doc["state"] = state
	updateInfo.Doc["sync_state"] = elastic.StateSyncIsTrue
	err := elastic.UpdateData(elastic.OrdTokenType, esId, updateInfo)
	if err != nil {
		return err
	}
	return nil
}

func UpdateDeployTokenInfo(esId string, info map[string]interface{}) error {
	// 处理数据
	updateInfo := elastic.UpdateInfo{}
	updateInfo.Doc = info
	err := elastic.UpdateData(elastic.DeployType, esId, updateInfo)
	if err != nil {
		return err
	}
	return nil
}
