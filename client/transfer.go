package client

import (
	"encoding/json"
	"fmt"
	"goBTC/models"
	"goBTC/ord"

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

func (c *BTCClient) GetTaprootAddressByPriKey(wifKey string) (string, error) {
	wif, err := btcutil.DecodeWIF(wifKey)
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
	return utxoTaprootAddress.EncodeAddress(), nil
}

func (c *BTCClient) CreateOrdi(filter models.CreateOrdFilter, commitTxPrevOutputList []*ord.PrevOutput) (string, error) {
	dataList := make([]ord.InscriptionData, 0)
	ordData := ord.InscriptionData{
		ContentType: filter.ContentType,
		Body:        []byte(filter.Body),
		RevealAddr:  filter.Destination,
	}
	dataList = append(dataList, ordData)

	request := ord.InscriptionRequest{
		CommitTxPrevOutputList: commitTxPrevOutputList,
		CommitFeeRate:          filter.TxFee,
		RevealFeeRate:          filter.TxFee,
		RevealOutValue:         546,
		InscriptionDataList:    dataList,
		ChangeAddress:          filter.ChangeAddress,
	}

	txs, err := ord.Inscribe(c.Params, &request)
	if err != nil {
		err = fmt.Errorf("send tx errr, %v", err)
		return "", err
	}

	fmt.Printf("commit txs: %+v\n", txs)

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
