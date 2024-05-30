package client

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
)

func (c *BTCClient) GetAddressByPrivateKey(priKey *btcec.PrivateKey) (*btcutil.AddressPubKeyHash, error) {
	pubKey := priKey.PubKey()
	pkHash := btcutil.Hash160(pubKey.SerializeCompressed())
	addr, err := btcutil.NewAddressPubKeyHash(pkHash, c.Params)
	if err != nil {
		return nil, err
	}
	return addr, nil
}

func (c *BTCClient) GetAddressByPKScript(pkScript []byte) (string, error) {
	_, addr, required, err := txscript.ExtractPkScriptAddrs(pkScript, c.Params)
	if err != nil {
		return "", err
	}
	if len(addr) == 0 || required == 0 {
		return "", fmt.Errorf("Not have address")
	}
	return addr[0].EncodeAddress(), nil
}
