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
	// GetBlockInfoByHash()
	// SignTx()
	// GetWitness()
	GetWitnessScript()
	if global.MysqlFlag {
		utils.SignalHandler("main", goBTC.Shutdown)
	}
}

func GetBlockInfoByHash() {
	// hash := "000000000000000000029730547464f056f8b6e2e0a02eaf69c24389983a04f5"
	hash := "00000000000000000000499cd89c4f19483a2081c2dcfbbbf7b2c2150c37d7b5"
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

func GetWitnessResByHash(hash string) (string, error) {
	// 查询Witness的铭文数据
	data, err := srv.GetRawTransactionByHash(hash)
	if err != nil {
		fmt.Printf("GetRawTransactionByHash error: %+v\n", err)
		return "", err
	}
	fmt.Printf("wch---- data: %+v\n", data)
	witness := client.GetTxWitnessByTxHex(data.Hex)
	if witness == "" {
		return "", nil
	}
	fmt.Printf("witness: %+v\n", witness)
	resList := client.GetScriptString(witness)
	if resList == nil {
		return "", nil
	}
	fmt.Printf("body len: %+v\n", resList.ContentSize)
	fmt.Printf("Brc20: %+v\n", resList.Brc20.Tick)
	return resList.Body, nil
}

func GetWitness() {
	// hash := "7fb631b7ed420c07b546ee4db8527a9523bbc44961f9983430166988cd6beeeb" // TEXT_1
	// hash := "bdbf2d7e385f650cbcba9a0ae83dc3f466dadc1e48732835e977cfefe2710b42" // TEXT_2
	// hash := "885441055c7bb5d1c54863e33f5c3a06e5a14cc4749cb61a9b3ff1dbe52a5bbb" // TEXT_3
	// hash := "ff4d5e838adfe81c8486ed8630be945badf9a5e75d07262f9d56964eba6ca032" // IMAGE_1
	// hash := "67df85eb1a66211b4e761d0b76464e5d07e758426214dab5d6fe42b664d979a4" // AUDIO_1
	// hash := "38d89d0506a5c936867b8a8c13b57d815cb2b2d86aee076ffec86b31c2cf51b5" // AUDIO_2
	// 铭文流转
	// hash := "5ee59cb5f2b88d1aa1dd7ef0f6263a2682412866e8cdb73275fa013429169623"
	hash := "231746e07440a6fa81d45f0d26e0510329175de1cac07b64c0a53faafb3b551d"

	body, _ := GetWitnessResByHash(hash)
	l := len(body)
	if l > 500 {
		fmt.Printf("body len: %s\n", l)
	} else {
		fmt.Printf("body: %s\n", body)
	}
}

func GetWitnessScript() {
	script := `70736274ff0100f4020000000300000000000000000000000000000000000000000000000000000000000000000000000000ffffffff00000000000000000000000000000000000000000000000000000000000000000100000000ffffffff40337d5bb0d29219ed84a8144f1e6039bf54bde8857351bc7c97aae5b0ffddd00000000000ffffffff0300000000000000001976a914000000000000000000000000000000000000000088ac00000000000000001976a914000000000000000000000000000000000000000088ac40420f000000000022512058a350ef35006b103344b06237ac02ac3133a462417348edf1ab8b0596780f37000000000001011f0000000000000000160014ae47938f7acd1623e6e10e1ebcc33c2a7cb6e30d0001011f0000000000000000160014ae47938f7acd1623e6e10e1ebcc33c2a7cb6e30d0001012b220200000000000022512058a350ef35006b103344b06237ac02ac3133a462417348edf1ab8b0596780f3701030483000000011720a9f6467`
	client.GetWitnessScript(script)
}
