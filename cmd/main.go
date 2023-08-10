package main

import (
	"fmt"
	"goBTC"
	"goBTC/client"
	"goBTC/global"
	"goBTC/models"
	"goBTC/utils"
)

var srv *client.BTCClient

func main() {
	fmt.Println("vim-go")
	goBTC.MustLoad("./config.yml")
	srv = global.Client
	// GetBlockInfoByHash()
	// SignTx()
	GetWitness()
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
	witness := client.GetTxWitnessByTxHex(data.Hex)
	if witness == "" {
		return "", nil
	}
	// fmt.Printf("witness: %+v\n", witness)
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

func SignTx() {
	body := fmt.Sprintf(`{"p":"brc-20","op":"%s","tick":"%s","amt":"%s"}`, "deploy", "yyds", "21000")
	filter := models.OrdInscribeData{
		ContentType: "text/plain;charset=utf-8",
		Body:        body,
		Destination: "",
		TxFee:       10,
	}
	_, err := srv.SignTx(filter)
	if err != nil {
		fmt.Printf("wch---- err: %+v\n", err)
		return
	}
}
