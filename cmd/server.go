package main

import (
	"fmt"
	"goBTC"
	"goBTC/client"
	"goBTC/elastic"
	"goBTC/global"
	"goBTC/models"
	"goBTC/ord"
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
	global.MysqlFlag = true
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
			group.POST("get_inscription", getInscriptionByFilter)
			group.POST("get_activity", getActivityByFilter)
			group.POST("/webp", webpHandler)
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

func getInscriptionByFilter(c *gin.Context) {
	var param models.GetBalanceParam
	err := c.BindJSON(&param)
	if err != nil {
		global.LOG.Error("required param error！", zap.Any("err", err))
		resp.FailParamError(c)
		return
	}
	// 分页
	pageFrom, pageSize := utils.PageFilter(param.Page, param.Limit)
	// 查询主币
	filter := make(map[string]interface{})
	if param.Address != "" {
		filter["owner_address"] = param.Address
	}
	if param.InscriptionId != "" {
		filter["inscription_id"] = param.InscriptionId
	}
	result, err := ord.GetBalanceInfo(elastic.InscribeInfoType, pageFrom, pageSize, filter)
	if err != nil {
		global.LOG.Error("ord.GetBalanceInfo", zap.Error(err))
		resp.FailWithMessage(err.Error(), c)
	}

	resp.OkWithData(result, c)
}

func getActivityByFilter(c *gin.Context) {
	var param models.GetBalanceParam
	err := c.BindJSON(&param)
	if err != nil {
		global.LOG.Error("required param error！", zap.Any("err", err))
		resp.FailParamError(c)
		return
	}
	// 分页
	pageFrom, pageSize := utils.PageFilter(param.Page, param.Limit)
	// 查询主币
	filter := make(map[string]interface{})
	if param.Address != "" {
		filter["owner_address"] = param.Address
	}
	if param.InscriptionId != "" {
		filter["inscription_id"] = param.InscriptionId
	}
	result, err := ord.GetBalanceInfo(elastic.ActivityType, pageFrom, pageSize, filter)
	if err != nil {
		global.LOG.Error("ord.GetBalanceInfo", zap.Error(err))
		resp.FailWithMessage(err.Error(), c)
	}

	resp.OkWithData(result, c)
}

func webpHandler(c *gin.Context) {
	//// // var buf bytes.Buffer
	//// var width, height int
	//// // var data []byte
	//// var err error
	//// hash := "ff4d5e838adfe81c8486ed8630be945badf9a5e75d07262f9d56964eba6ca032" // IMAGE_1
	//// res, err := GetWitnessResByHash(hash)
	//// if err != nil {
	//// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	//// 	return
	//// }
	//// data, err := hex.DecodeString(res)
	//// if err != nil {
	//// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	//// 	return
	//// }

	//// // GetInfo
	//// if width, height, _, err = webp.GetInfo(data); err != nil {
	//// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	//// 	return
	//// }
	//// fmt.Printf("width = %d, height = %d\n", width, height)

	//// // GetMetadata
	//// if metadata, err := webp.GetMetadata(data, "ICCP"); err != nil {
	//// 	fmt.Printf("Metadata: err = %v\n", err)
	//// } else {
	//// 	fmt.Printf("Metadata: %s\n", string(metadata))
	//// }

	//// // Decode webp
	//// img, err := webp.Decode(bytes.NewReader(data))
	//// if err != nil {
	//// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	//// 	return
	//// }
	//// if err := png.Encode(w, img); err != nil {
	//// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	//// 	return
	//// }

	// // Encode lossless webp
	// if err = webp.Encode(&buf, m, &webp.Options{Lossless: true}); err != nil {
	// 	log.Println(err)
	// }
	// if err = ioutil.WriteFile("output.webp", buf.Bytes(), 0666); err != nil {
	// 	log.Println(err)
	// }

	fmt.Println("Save output.webp ok")
}
