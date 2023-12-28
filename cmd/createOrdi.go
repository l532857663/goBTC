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
	SignTx()
	if global.MysqlFlag {
		utils.SignalHandler("main", goBTC.Shutdown)
	}
}

func SignTx() {
	body := fmt.Sprintf(`{"p":"brc-20","op":"%s","tick":"%s","amt":"%s"}`, "deploy", "yyds", "21000")
	filter := models.OrdInscribeData{
		ContentType: "text/plain;charset=utf-8",
		Body:        body,
		Destination: "",
		TxFee:       3,
	}
	_, err := srv.SignTx(filter)
	if err != nil {
		fmt.Printf("wch---- err: %+v\n", err)
		return
	}
}
