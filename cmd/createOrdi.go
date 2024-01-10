package main

import (
	"fmt"
	"goBTC"
	"goBTC/client"
	"goBTC/global"
	"goBTC/models"
	"goBTC/ord"
	"goBTC/utils"
)

var srv *client.BTCClient

func main() {
	fmt.Println("vim-go")
	goBTC.MustLoad("./config.yml")
	srv = global.Client
	CreateOrdi()
	if global.MysqlFlag {
		utils.SignalHandler("main", goBTC.Shutdown)
	}
}

func CreateOrdi() {
	hexPrivateKey := "cQSreoKBANpfNxLHD6v1crHE3rz44Q7hZPsV2XaJVQv6dA5eXGQV"
	// utxoTaprootAddress, err := srv.GetTaprootAddressByPriKey(hexPrivateKey)
	// if err != nil {
	// 	fmt.Printf("wch---- err: %+v\n", err)
	// 	return
	// }
	commitTxPrevOutputList := make([]*ord.PrevOutput, 0)
	commitTxPrevOutputList = append(commitTxPrevOutputList, &ord.PrevOutput{
		TxId:       "12fd266d657045ef596a4a611e857dbb7331eeffcb4ade835a27237312898a72",
		VOut:       1,
		Amount:     1759377,
		Address:    "tb1q7vs6lm6tercj0a725ksy23vd8s27czlagpsd2m",
		PrivateKey: hexPrivateKey,
	})
	body := fmt.Sprintf(`{"p":"brc-20","op":"%s","tick":"%s","amt":"%s"}`, "mint", "dddd", "21000000")
	filter := models.CreateOrdFilter{
		ContentType:   "text/plain;charset=utf-8",
		Body:          body,
		Destination:   "tb1q7vs6lm6tercj0a725ksy23vd8s27czlagpsd2m",
		TxFee:         3,
		ChangeAddress: "tb1pg0uc7ujx6rplw4wj73etg505jh49k63s7wc3kyngf73ze7ffue4skru6ld",
	}
	_, err1 := srv.CreateOrdi(filter, commitTxPrevOutputList)
	if err1 != nil {
		fmt.Printf("wch---- err: %+v\n", err1)
		return
	}
}
