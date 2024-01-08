package client

import (
	"log"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

// 查询节点块高
func (c *BTCClient) GetBlockCount() (int64, error) {
	return c.Client.GetBlockCount()
}

// 根据块高查HASH
func (c *BTCClient) GetBlockHashByHeight(height int64) (string, error) {
	hash, err := c.Client.GetBlockHash(height)
	if err != nil {
		return "", err
	}
	return hash.String(), err
}

// 根据块高查询数据
func (c *BTCClient) GetBlockInfoByHeight(height int64) (*wire.MsgBlock, error) {
	hash, err := c.Client.GetBlockHash(height)
	if err != nil {
		return nil, err
	}
	return c.Client.GetBlock(hash)
}

// 根据块HASH查询数据
func (c *BTCClient) GetBlockInfoByHash(hash string) (*wire.MsgBlock, error) {
	// hash处理
	h, err := chainhash.NewHashFromStr(hash)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return c.Client.GetBlock(h)
}

// 根据块HASH查询数据
func (c *BTCClient) GetBlockStatus(hashOrHeight string) (*btcjson.GetBlockStatsResult, error) {
	return c.Client.GetBlockStats(hashOrHeight, nil)
}
