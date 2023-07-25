package client

import (
	"encoding/json"
	"fmt"
	"goBTC/models"
	"goBTC/ord"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/wire"

	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
)

func (c *BTCClient) GetTransactionByHash(hash string) (*btcutil.Tx, error) {
	h, err := chainhash.NewHashFromStr(hash)
	if err != nil {
		return nil, err
	}
	return c.Client.GetRawTransaction(h)
}

func (c *BTCClient) GetRawTransactionByHash(hash string) (*btcjson.TxRawResult, error) {
	h, err := chainhash.NewHashFromStr(hash)
	if err != nil {
		return nil, err
	}
	return c.Client.GetRawTransactionVerbose(h)
}

func (c *BTCClient) SignTx(filter models.OrdInscribeData) (string, error) {
	// hexPrivateKey := "Kz5sBtKKjK8QDWs1ph6tFXnARdgePx21VMj1v77nzJG6NQyVnihC"
	hexPrivateKey := "cQSreoKBANpfNxLHD6v1crHE3rz44Q7hZPsV2XaJVQv6dA5eXGQV"
	wif, err := btcutil.DecodeWIF(hexPrivateKey)
	if err != nil {
		return "", fmt.Errorf("SignTx DecodeWIF fatal, " + err.Error())
	}
	if !wif.IsForNet(c.Params) {
		return "", fmt.Errorf("SignTx IsForNet fatal")
	}
	pubK := wif.PrivKey.PubKey()
	utxoTaprootAddress, err := btcutil.NewAddressTaproot(schnorr.SerializePubKey(txscript.ComputeTaprootKeyNoScript(pubK)), c.Params)
	if err != nil {
		return "", err
	}
	fmt.Printf("wch----- taproot: %+v, pubKey: %x\n", utxoTaprootAddress.EncodeAddress(), utxoTaprootAddress.ScriptAddress())
	// 查询未花费的UTXO列表
	unspentList, err := c.MempoolClient.ListUnspent(utxoTaprootAddress)
	if err != nil {
		err = fmt.Errorf("GetListUnspent error: %+v", err)
		return "", err
	}
	if len(unspentList) == 0 {
		err = fmt.Errorf("no utxo for %s", utxoTaprootAddress)
		return "", err
	}
	fmt.Printf("unspentList: %+v\n", unspentList)
	vinAmount := 0
	commitTxOutPointList := make([]*wire.OutPoint, 0)
	commitTxPrivateKeyList := make([]*btcec.PrivateKey, 0)
	for i := range unspentList {
		value := unspentList[i].Output.Value
		fmt.Printf("wch---- unspent value: %+v\n", value)
		if value < 10000 {
			continue
		}
		commitTxOutPointList = append(commitTxOutPointList, unspentList[i].Outpoint)
		commitTxPrivateKeyList = append(commitTxPrivateKeyList, wif.PrivKey)
		vinAmount += int(value)
	}

	dataList := make([]ord.InscriptionData, 0)

	ordData := ord.InscriptionData{
		ContentType: filter.ContentType,
		Body:        []byte(filter.Body),
		Destination: utxoTaprootAddress.EncodeAddress(),
	}
	dataList = append(dataList, ordData)

	request := ord.InscriptionRequest{
		CommitTxOutPointList:   commitTxOutPointList,
		CommitTxPrivateKeyList: commitTxPrivateKeyList,
		CommitFeeRate:          filter.TxFee,
		FeeRate:                filter.TxFee,
		DataList:               dataList,
		SingleRevealTxOnly:     false,
	}

	// tool, err := ord.NewInscriptionTool(c.Params, c.Client, &request)
	tool, err := ord.NewInscriptionToolWithBtcApiClient(c.Params, c.MempoolClient, &request)
	if err != nil {
		fmt.Printf("wch------ ord.NewInscriptionToolWithBtcApiClient error: %+v\n", err)
		return "", err
	}
	tool.ShowInfo()
	// tool.DealWithTxByHex()
	return "", nil

	commitTxHash, revealTxHashList, _, _, err := tool.Inscribe()
	if err != nil {
		err = fmt.Errorf("send tx errr, %v", err)
		return "", err
	}

	txid := commitTxHash.String()
	fmt.Printf("commit txid: %+v\n", txid)
	for i := range revealTxHashList {
		// txids = append(txids, revealTxHashList[i].String())
		revealTxid := revealTxHashList[i].String()
		fmt.Printf("reveal txid: %+v\n", revealTxid)
	}

	return "", nil
}

func (c *BTCClient) SendRawTransaction(hexTx string) (string, error) {
	method := "sendrawtransaction"
	var txHash string
	marshalledParam, err := json.Marshal(hexTx)
	if err != nil {
		return txHash, err
	}
	rawMessage := json.RawMessage(marshalledParam)
	result, err := c.Client.RawRequest(method, []json.RawMessage{rawMessage})
	if err != nil {
		return txHash, err
	}
	err = json.Unmarshal(result, &txHash)
	if err != nil {
		return txHash, err
	}
	return txHash, nil
}
