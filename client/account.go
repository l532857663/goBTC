package client

import (
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
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
