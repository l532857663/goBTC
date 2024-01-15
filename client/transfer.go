package client

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"goBTC/models"
	"goBTC/ord"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

var (
	ToolsMap = make(map[string]InscriptionBuilder)
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

func (c *BTCClient) ContinueTransfer(builderStr, psbtData string, inputKey []string, commitTxPrevOutputList []*ord.PrevOutput) (*ord.InscribeTxs, error) {
	tool := &ord.InscriptionBuilder{
		Network:                   c.Params,
		CommitTxPrevOutputFetcher: txscript.NewMultiPrevOutFetcher(nil),
		RevealTxPrevOutputFetcher: txscript.NewMultiPrevOutFetcher(nil),
	}
	err := json.Unmarshal([]byte(builderStr), tool)
	if err != nil {
		return nil, err
	}
	for i, prikey := range inputKey {
		privateKeyBytes, _ := hex.DecodeString(prikey)
		privateKey, _ := btcec.PrivKeyFromBytes(privateKeyBytes)
		tool.InscriptionTxCtxDataList[i].PrivateKey = privateKey
	}
	for _, prevOutput := range commitTxPrevOutputList {
		txHash, err := chainhash.NewHashFromStr(prevOutput.TxId)
		if err != nil {
			return nil, err
		}
		outPoint := wire.NewOutPoint(txHash, prevOutput.VOut)
		pkScript, err := ord.AddrToPkScript(prevOutput.Address, tool.Network)
		if err != nil {
			return nil, err
		}
		txOut := wire.NewTxOut(prevOutput.Amount, pkScript)
		tool.CommitTxPrevOutputFetcher.AddPrevOut(*outPoint, txOut)
	}

	tool.UpdateCommitInfo(psbtData)
	err = tool.CompleteRevealTx()
	if err != nil {
		return nil, err
	}
	commitTx, err := tool.GetCommitTxHex()
	if err != nil {
		return nil, err
	}
	revealTxs, err := tool.GetRevealTxHexList()
	if err != nil {
		return nil, err
	}

	commitTxFee, revealTxFees := tool.CalculateFee()

	return &ord.InscribeTxs{
		CommitTx:     commitTx,
		RevealTxs:    revealTxs,
		CommitTxFee:  commitTxFee,
		RevealTxFees: revealTxFees,
		CommitAddrs:  tool.CommitAddrs,
	}, nil
}

func (c *BTCClient) CreateInscribe(filter models.CreateOrdFilter) (*models.CreateTransferResp, error) {
	dataList := make([]InscriptionData, 0)
	ordData := InscriptionData{
		ContentType: filter.ContentType,
		Body:        []byte(filter.Body),
		RevealAddr:  filter.Destination,
	}
	dataList = append(dataList, ordData)

	// 查询未花费的UTXO列表
	address := filter.Destination
	addr, err := btcutil.DecodeAddress(address, c.Params)
	if err != nil {
		fmt.Printf("invalid recipet address: %v", err)
		return nil, err
	}
	unspendList, err := c.MempoolClient.ListUnspent(addr)
	if err != nil {
		fmt.Printf("GetListUnspent error: %+v", err)
		return nil, err
	}
	if len(unspendList) == 0 {
		err := fmt.Errorf("not enougth utxo for %v", addr)
		return nil, err
	}
	commitTxPrevOutputList := make([]*PrevOutput, 0)
	for _, unspend := range unspendList {
		amount := unspend.Output.Value
		if amount < 1000 {
			// 小于1000的有可能是铭文UTXO
			continue
		}
		commitTxPrevOutputList = append(commitTxPrevOutputList, &PrevOutput{
			TxId:       unspend.Outpoint.Hash.String(),
			VOut:       unspend.Outpoint.Index,
			Amount:     amount,
			Address:    address,
			PrivateKey: DefaultKey,
		})
	}

	request := &InscriptionRequest{
		CommitTxPrevOutputList: commitTxPrevOutputList,
		CommitFeeRate:          filter.TxFee,
		RevealFeeRate:          filter.TxFee,
		RevealOutValue:         546,
		InscriptionDataList:    dataList,
		ChangeAddress:          filter.ChangeAddress,
	}

	tool, err := NewInscriptionTool(c.Params, request)
	if err != nil && err.Error() == "insufficient balance" {
		return nil, err
	}
	// 使用UNISAT签名处理
	pbstData, err := tool.GetSignHash()
	if err != nil {
		err = fmt.Errorf("tool.GetSignHash error, %v", err)
		return nil, err
	}

	hasher := md5.New()
	if _, err := hasher.Write(pbstData); err != nil {
		return nil, err
	}
	key := hex.EncodeToString(hasher.Sum(nil))
	ToolsMap[key] = *tool

	// 算手续费
	err = tool.signCommitTx()
	if err != nil {
		return nil, errors.New("sign commit tx error")
	}
	err = tool.CompleteRevealTx()
	if err != nil {
		return nil, err
	}
	commitTxFee, revealTxFees := tool.CalculateFee()
	revealTxFee := int64(0)
	for _, f := range revealTxFees {
		revealTxFee += f
	}

	respData := &models.CreateTransferResp{
		PSBTData:    hex.EncodeToString(pbstData),
		Key:         key,
		CommitFee:   commitTxFee,
		RevealFee:   revealTxFee,
		RevealValue: request.RevealOutValue,
	}

	return respData, nil
}

func (c *BTCClient) SendTransfer(key, psbtData string) (*InscribeTxs, error) {
	tool := ToolsMap[key]

	tool.UpdateCommitInfo(psbtData)
	err := tool.CompleteRevealTx()
	if err != nil {
		return nil, err
	}
	commitTx, err := tool.GetCommitTxHex()
	if err != nil {
		return nil, err
	}
	revealTxs, err := tool.GetRevealTxHexList()
	if err != nil {
		return nil, err
	}

	commitTxFee, revealTxFees := tool.CalculateFee()

	return &InscribeTxs{
		CommitTx:     commitTx,
		RevealTxs:    revealTxs,
		CommitTxFee:  commitTxFee,
		RevealTxFees: revealTxFees,
		CommitAddrs:  tool.CommitAddrs,
	}, nil
}
