package ord

import (
	"goBTC/elastic"

	"github.com/btcsuite/btcd/btcjson"
)

func UpdateInscribeInfoOwner(txId string, txInfo *btcjson.TxRawResult) error {
	// 处理数据
	addr := txInfo.Vout[0].ScriptPubKey.Address
	vout := GetInscribeOutputStr(txInfo.Txid)
	esId := GetInscribeIdStr(txId)

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
	updateInfo := elastic.Query{}
	searchInfo.Match = make(map[string]interface{})
	searchInfo.Match["tx_hash"] = txId
	err := elastic.DeleteData(elastic.ActivityType, txId, updateInfo)
	if err != nil {
		return err
	}
	return nil
}
