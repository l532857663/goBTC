package ord

import (
	"goBTC/elastic"
	"goBTC/models"
	"goBTC/utils"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/shopspring/decimal"
)

func SaveInscribeInfoByTxInfo(blockHeight int64, res *models.OrdInscribeData, txInfo *btcjson.TxRawResult) error {
	// tx
	txId := txInfo.Txid
	addr := txInfo.Vout[0].ScriptPubKey.Address
	vout := GetInscribeOutputStr(txId)
	inscribeId := GetInscribeIdStr(txId)
	ord := ""
	tokenType := elastic.InscribeTypeNFT
	state := elastic.InscriptionStatePending
	// 判断铭文类型
	if res.Brc20 != nil {
		ord = res.Body
		tokenType = elastic.InscribeTypeBRC20
	} else {
		if len(res.Body) < 500 {
			ord = res.Body
		}
	}
	// 判断铭文数据是否可用
	if res.TxHaveInscribe != "" {
		state = elastic.InscriptionStateInvalid
	}

	info := &elastic.InscribeInfo{
		InscribeId:      inscribeId,
		InscribeContent: ord,
		InscribeType:    tokenType,
		TxHash:          txId,
		ContentType:     res.ContentType,
		OwnerAddress:    addr,
		OwnerOutput:     vout,
		GenesisAddress:  addr,
		GenesisOutput:   vout,
		BlockHeight:     blockHeight,
		State:           state,
		SyncState:       elastic.StateSyncIsFalse,
		Brc20Info:       res.Brc20,
	}
	err := elastic.CreateData(elastic.InscribeInfoType, inscribeId, info)
	if err != nil {
		return err
	}
	return nil
}

func SaveInscribeBrc20ByTxInfo(blockHeight int64, res *models.OrdInscribeData, txInfo *btcjson.TxRawResult) error {
	// tx
	txId := txInfo.Txid
	inscribeId := GetInscribeIdStr(txId)
	addr := txInfo.Vout[0].ScriptPubKey.Address
	vout := GetInscribeOutputStr(txId)
	utxo := decimal.NewFromFloat(txInfo.Vout[0].Value).String()
	action := res.Brc20.OP
	// 判断铭文数据是否可用
	state := elastic.InscriptionStatePending
	if res.TxHaveInscribe != "" {
		state = elastic.InscriptionStateInvalid
	}

	ord := &elastic.OrdToken{
		InscribeId:   inscribeId,
		TxHash:       txId,
		OwnerAddress: addr,
		OwnerOutput:  vout,
		Value:        utxo,
		Action:       action,
		Brc20Info:    res.Brc20,
		BlockHeight:  blockHeight,
		BlockTime:    txInfo.Blocktime,
		State:        state,
		SyncState:    elastic.StateSyncIsFalse,
	}
	err := elastic.CreateData(elastic.OrdTokenType, inscribeId, ord)
	if err != nil {
		return err
	}
	return nil
}

func SaveInscribeActivity(esId, toAddr string, res *models.OrdInscribeData, txInfo *btcjson.TxRawResult) error {
	inscribeId := GetInscribeIdStr(txInfo.Txid)
	txHash := txInfo.Txid
	fromAddr := ""
	owner := ""
	activityType := ""
	activityAction := ""
	tokenType := ""
	var receiveInfo *elastic.ActivityInfo
	if res == nil {
		// 铭文ID不变
		inscribeId = esId
		// 铭文转移(普通交易)
		inscribeInfo, err := elastic.GetDataById(elastic.InscribeInfoType, inscribeId)
		if err != nil {
			return err
		}
		info := &elastic.InscribeInfo{}
		err = utils.Map2Struct(inscribeInfo.Source, info)
		if err != nil {
			return err
		}
		// 获取相应的数据
		fromAddr = info.OwnerAddress
		owner = info.OwnerAddress
		activityType = elastic.ActivityTypeTransfer
		activityAction = elastic.ActivityActionSend
		tokenType = info.InscribeType
		// 生成对应的接收数据
		receiveInfo = &elastic.ActivityInfo{
			Owner:          toAddr,
			ActivityAction: elastic.ActivityActionReceive,
		}
	} else {
		// 铭文操作(铭文交易)
		activityType = elastic.ActivityTypeInscribed
		activityAction = elastic.ActivityActionMint
		tokenType = elastic.InscribeTypeNFT
		owner = txInfo.Vout[0].ScriptPubKey.Address
		if res.Brc20 != nil {
			tokenType = elastic.InscribeTypeBRC20
			// BRC20的操作还有 deploy、transfer
			if res.Brc20.OP == elastic.ActivityActionDeploy {
				activityAction = elastic.ActivityActionDeploy
			} else if res.Brc20.OP == elastic.ActivityActionTransfer {
				activityAction = elastic.ActivityActionTransfer
			}
		}
	}
	active := &elastic.ActivityInfo{
		InscribeId:     inscribeId,
		InscribeType:   tokenType,
		TxHash:         txHash,
		ActivityType:   activityType,
		ActivityAction: activityAction,
		Owner:          owner,
		From:           fromAddr,
		To:             toAddr,
		BlockTime:      txInfo.Blocktime,
	}
	esId = txHash + "_" + active.ActivityAction
	err := elastic.CreateData(elastic.ActivityType, esId, active)
	if err != nil {
		return err
	}
	// 创建对应的接收数据
	if res == nil && receiveInfo != nil {
		active.ActivityAction = receiveInfo.ActivityAction
		active.Owner = receiveInfo.Owner
		esId = txHash + "_" + active.ActivityAction
		err := elastic.CreateData(elastic.ActivityType, esId, active)
		if err != nil {
			return err
		}
	}
	return nil
}
