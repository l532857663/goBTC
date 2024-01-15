package main

import (
	"fmt"
	"goBTC"
	"goBTC/client"
	"goBTC/global"
	"goBTC/models"
	"goBTC/utils"
	"goBTC/utils/resp"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	srv *client.BTCClient
)

func main() {
	fmt.Println("vim-go")
	// global.MysqlFlag = true
	goBTC.MustLoad("./config.yml")
	srv = global.Client
	RunServer()
	if global.MysqlFlag {
		utils.SignalHandler("main", goBTC.Shutdown)
	}
}

func RunServer() {
	Router := gin.Default()
	Groups := Router.Group("ord_info")
	{
		group := Groups.Group("v1")
		{
			group.POST("createTransfer", createInscribe)
			group.POST("sendTransfer", sendTransfer)
		}
	}
	// address := fmt.Sprintf(":%d", config.CONFIG.System.Addr)
	address := fmt.Sprintf(":%d", 4396)
	httpType := "HTTP"
	server := &http.Server{
		Addr:    address,
		Handler: Router,
	}
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			global.LOG.Info("ListenAndServe err", zap.String("err", err.Error()))
			panic(err)
		}
	}()

	// 保证文本顺序输出
	time.Sleep(10 * time.Microsecond)
	global.LOG.Info("server run success on ", zap.String("address", address), zap.String("type", httpType))
}

func createInscribe(c *gin.Context) {
	var param models.CreateTransferParam
	err := c.BindJSON(&param)
	if err != nil {
		global.LOG.Error("required param error！", zap.Any("err", err))
		resp.FailParamError(c)
		return
	}
	addr := param.Address
	body := fmt.Sprintf(`{"p":"brc-20","op":"%s","tick":"%s","amt":"%s"}`, "transfer", param.Tick, param.Amount)
	txFeeRate := param.FeeRate

	// 处理数据
	filter := models.CreateOrdFilter{
		ContentType:   "text/plain;charset=utf-8",
		Body:          body,
		Destination:   addr,
		TxFee:         txFeeRate,
		ChangeAddress: addr,
	}
	respData, err := srv.CreateInscribe(filter)
	if err != nil {
		global.LOG.Error("srv.CreateInscribe！", zap.Any("err", err))
		resp.FailWithCodeMessage(200, err.Error(), c)
		return
	}

	resp.OkWithData(respData, c)
	return
}

func sendTransfer(c *gin.Context) {
	var param models.SendTransferParam
	err := c.BindJSON(&param)
	if err != nil {
		global.LOG.Error("required param error！", zap.Any("err", err))
		resp.FailParamError(c)
		return
	}
	txs, err := srv.SendTransfer(param.Key, param.PSBTData)
	if err != nil {
		global.LOG.Error("srv.SendTransfer error！", zap.Any("err", err))
		resp.FailWithCodeMessage(200, err.Error(), c)
		return
	}
	global.LOG.Info("commit txs:", zap.Any("info:", txs))
	commitTxHash, err := srv.SendRawTransaction(txs.CommitTx)
	if err != nil {
		global.LOG.Error("srv.SendRawTransaction commit error！", zap.Any("err", err))
		resp.FailWithCodeMessage(200, err.Error(), c)
		return
	}
	revealTxHash := []string{}
	for i, revealTx := range txs.RevealTxs {
		txHash, err := srv.SendRawTransaction(revealTx)
		if err != nil {
			global.LOG.Error("srv.SendRawTransaction reveal error！", zap.Any("i", i), zap.Any("err", err))
			resp.FailWithCodeMessage(200, err.Error(), c)
			return
		}
		revealTxHash = append(revealTxHash, txHash)
	}
	// 返回结果
	respData := &models.SendTransferResp{
		CommitTxHash: commitTxHash,
		RevealTxHash: revealTxHash,
	}

	resp.OkWithData(respData, c)
	return
}
