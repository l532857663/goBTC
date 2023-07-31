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
