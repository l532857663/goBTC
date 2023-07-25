package client

import (
	"fmt"
	"testing"
)

var srv *BTCClient

func init() {
	// 构建节点客户端
	nodeInfo := BTC_GETBLOCK_MAIN
	var err error
	srv, err = NewBTCClient(nodeInfo)
	if err != nil {
		fmt.Printf("NewBTCClient error: %+v, nodeInfo: %+v\n", err, nodeInfo)
		return
	}
}

func Test_GetBlockCount(t *testing.T) {
	blockHigh, err := srv.GetBlockCount()
	if err != nil {
		fmt.Printf("GetBlockInfoByHash error: %+v\n", err)
		return
	}
	fmt.Printf("%+v\n", blockHigh)
}

func Test_GetBlockHashByHeight(t *testing.T) {
	blockHigh := int64(789498)
	blockHash, err := srv.GetBlockHashByHeight(blockHigh)
	if err != nil {
		fmt.Printf("GetBlockHashByHeight error: %+v\n", err)
		return
	}
	fmt.Printf("%+v\n", blockHash)
}

func Test_GetBlockInfoByHash(t *testing.T) {
	hash := "000000000000000000018a6a6d1cf6d0b16f1d2cbd303cf4c6b5eadf2dc40a0f"
	blockInfo, err := srv.GetBlockInfoByHash(hash)
	if err != nil {
		fmt.Printf("GetBlockInfoByHash error: %+v, hash: %+v\n", err, hash)
		return
	}
	fmt.Printf("%+v\n", blockInfo.Header)
}

func Test_GetBlockInfo(t *testing.T) {
	blockHigh, err := srv.GetBlockCount()
	if err != nil {
		fmt.Printf("GetBlockInfoByHash error: %+v\n", err)
		return
	}
	fmt.Printf("%+v\n", blockHigh)
	blockHash, err := srv.GetBlockHashByHeight(blockHigh)
	if err != nil {
		fmt.Printf("GetBlockHashByHeight error: %+v\n", err)
		return
	}
	fmt.Printf("%+v\n", blockHash)
	blockInfo, err := srv.GetBlockInfoByHash(blockHash)
	if err != nil {
		fmt.Printf("GetBlockInfoByHash error: %+v, blockHash: %+v\n", err, blockHash)
		return
	}
	fmt.Printf("%+v\n", blockInfo.Header)
}
