package main

import (
	"fmt"
	"goBTC"
	"goBTC/client"
	"goBTC/global"
	"goBTC/utils"
)

var srv *client.BTCClient

func main() {
	fmt.Println("vim-go")
	goBTC.MustLoad("./config.yml")
	srv = global.Client
	GetBlockInfoByHash()
	// SignTx()
	// GetWitness()
	// GetWitnessScript()
	if global.MysqlFlag {
		utils.SignalHandler("main", goBTC.Shutdown)
	}
}

func GetBlockInfoByHash() {
	// hash := "000000000000000000029730547464f056f8b6e2e0a02eaf69c24389983a04f5"
	hash := "00000000ad1471fd5e8b872b9124b4d334fa594cb4ea87acac8895051bfdda1f"
	blockInfo, err := srv.GetBlockInfoByHash(hash)
	if err != nil {
		fmt.Printf("GetBlockInfoByHash error: %+v\n", err)
		return
	}
	fmt.Printf("%+v\n", blockInfo.Header)
	for i, tx := range blockInfo.Transactions {
		witnessStr := client.GetTxWitness(tx)
		if witnessStr == "" {
			continue
		}
		res := client.GetScriptString(witnessStr)
		if res != nil {
			fmt.Printf("[%d] txHash: %s, [ord] : %v\n", i, tx.TxHash().String(), res.ContentType)
		}
	}
}

func GetWitnessScript() {
	script := `70736274ff0100f4020000000300000000000000000000000000000000000000000000000000000000000000000000000000ffffffff00000000000000000000000000000000000000000000000000000000000000000100000000ffffffff40337d5bb0d29219ed84a8144f1e6039bf54bde8857351bc7c97aae5b0ffddd00000000000ffffffff0300000000000000001976a914000000000000000000000000000000000000000088ac00000000000000001976a914000000000000000000000000000000000000000088ac40420f000000000022512058a350ef35006b103344b06237ac02ac3133a462417348edf1ab8b0596780f37000000000001011f0000000000000000160014ae47938f7acd1623e6e10e1ebcc33c2a7cb6e30d0001011f0000000000000000160014ae47938f7acd1623e6e10e1ebcc33c2a7cb6e30d0001012b220200000000000022512058a350ef35006b103344b06237ac02ac3133a462417348edf1ab8b0596780f3701030483000000011720a9f6467`
	client.GetWitnessScript(script)
}
